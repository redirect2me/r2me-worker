package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
)

type mapResultKeyType string

const mapResultKey mapResultKeyType = "map_result"

type MapResult struct {
	Action      string `json:"action"`
	Destination string `json:"destination"`
	ResultCode  string `json:"result_code"`
	StatusCode  int    `json:"status_code"`
	Debug       bool   `json:"debug,omitempty"`
}

func (mr *MapResult) LogValue() slog.Value {
	if mr == nil {
		return slog.AnyValue(nil)
	}
	return slog.GroupValue(
		slog.Attr{Key: "action", Value: slog.StringValue(mr.Action)},
		slog.Attr{Key: "destination", Value: slog.StringValue(mr.Destination)},
		slog.Attr{Key: "result_code", Value: slog.StringValue(mr.ResultCode)},
		slog.Attr{Key: "status_code", Value: slog.IntValue(mr.StatusCode)},
		slog.Attr{Key: "debug", Value: slog.BoolValue(mr.Debug)},
	)
}

var supportUrl = "https://www.redirect2.me/support/${error}.html"

func makeError(errcode string) *MapResult {
	return &MapResult{
		Destination: strings.Replace(supportUrl, "${error}", errcode, -1),
		ResultCode:  errcode,
		StatusCode:  http.StatusTemporaryRedirect,
	}
}

func mapAddWww(r *http.Request) *MapResult {

	host := r.Host

	if strings.HasPrefix(host, "www.") {
		return makeError("error-already-has-www")
	}

	loc := url.URL(*r.URL)
	loc.Scheme = getScheme(r)
	loc.Host = "www." + host
	return &MapResult{
		Destination: loc.String(),
		ResultCode:  "success",
		StatusCode:  http.StatusTemporaryRedirect,
	}
}

func customDNSDialer(ctx context.Context, network, address string) (net.Conn, error) {
	d := net.Dialer{}
	return d.DialContext(ctx, "udp", "1.1.1.1:53")
}

func lookupTxt(domain string) string {
	prefix := "redirect2me="

	r := net.Resolver{
		PreferGo: true,
		Dial:     customDNSDialer,
	}

	results, err := r.LookupTXT(context.Background(), domain)

	if err != nil || len(results) == 0 {
		return ""
	}

	for rIndex := 0; rIndex < len(results); rIndex++ {
		result := results[rIndex]
		Logger.Info("DNS TXT Lookup", "domain", domain, "rIndex", rIndex, "result", result)
		if len(result) > len(prefix)+1 && result[0:len(prefix)] == prefix {
			return results[rIndex][len(prefix):]
		}
	}

	return ""
}

func mapLookup(r *http.Request) *MapResult {

	host := r.Host

	lookup := lookupTxt(host)

	if lookup == "" {
		return makeError("error-lookup-not-found")
	}

	u, err := url.Parse(lookup)
	if err != nil {
		return makeError("lookup-parse-error")
	}

	result := url.URL(*r.URL)
	if u.Scheme == "" {
		result.Scheme = getScheme(r)
	} else {
		result.Scheme = u.Scheme
	}
	result.Host = u.Host
	if u.Path != "" {
		result.Path = u.Path
	}
	if u.RawQuery != "" {
		result.RawQuery = u.RawQuery
	}
	return &MapResult{
		Destination: result.String(),
		ResultCode:  "success",
		StatusCode:  http.StatusTemporaryRedirect,
	}
}

func mapRemoveWww(r *http.Request) *MapResult {

	host := r.Host

	if !strings.HasPrefix(host, "www.") {
		return makeError("error-does-not-have-www")
	}

	loc := url.URL(*r.URL)
	loc.Scheme = getScheme(r)
	loc.Host = host[4:]
	return &MapResult{
		Destination: loc.String(),
		ResultCode:  "success",
		StatusCode:  http.StatusTemporaryRedirect,
	}
}

func RequestLogValue(r *http.Request) slog.Value {
	// Redact the password field and return a group of attributes
	return slog.GroupValue(
		slog.Attr{Key: "source_ip", Value: slog.StringValue(r.RemoteAddr)},
		slog.Attr{Key: "user_agent", Value: slog.StringValue(r.Header.Get("User-Agent"))},
		slog.Attr{Key: "host", Value: slog.StringValue(r.Host)},
		slog.Attr{Key: "scheme", Value: slog.StringValue(getScheme(r))},
	)
}

func getScheme(r *http.Request) string {

	proxyValue := r.Header.Get("X-Forwarded-Proto")
	if proxyValue != "" {
		return proxyValue
	}

	urlValue := r.URL.Scheme
	if urlValue != "" {
		return urlValue
	}

	return "http"
}

func GetMapper(action string) (http.HandlerFunc, error) {
	var mapFn func(r *http.Request) *MapResult

	switch action {
	case "addwww":
		mapFn = mapAddWww
	case "removewww":
		mapFn = mapRemoveWww
	case "lookup":
		mapFn = mapLookup
	default:
		Logger.Error("unknown action for mapper", "action", action)
		return nil, errors.New("ErrUnknownAction")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var result *MapResult

		_, ipErr := netip.ParseAddr(r.Host)
		if ipErr == nil {
			result = makeError("error-host-is-ip-address")
		} else {
			result = mapFn(r)
		}
		result.Action = Config.Action

		if r.Header.Get("X-Redirect2Me-Debug") == "1" {
			result.Debug = true
		}

		RealtimeSend(result)

		if lw, ok := w.(*loggingWriter); ok {
			lw.mapResult = result
		} else {
			// this should never happen
			Logger.Warn("ResponseWriter is not a loggingWriter, cannot attach mapResult", "request", RequestLogValue(r), "result", result)
		}

		if result.Debug {
			HandleJson(w, r, result)
		} else {
			w.Header().Set("X-Redirect-By", "redirect2.me")
			http.Redirect(w, r, result.Destination, result.StatusCode)
			Metrics.RedirectsTotal.WithLabelValues(result.ResultCode).Inc()
		}
	}, nil
}
