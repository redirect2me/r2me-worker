package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type MapResult struct {
	Destination string `json:"destination"`
	StatusCode  int    `json:"status_code"`
}

var supportUrl = "https://www.redirect2.me/support/${error}.html"

func makeError(errcode string) MapResult {
	return MapResult{
		Destination: strings.Replace(supportUrl, "${error}", errcode, -1),
		StatusCode:  http.StatusTemporaryRedirect,
	}
}

func mapAddWww(r *http.Request) MapResult {

	host := r.Host

	if strings.HasPrefix(host, "www.") {
		return makeError("error-already-has-www")
	}

	loc := url.URL(*r.URL)
	loc.Scheme = getScheme(r)
	loc.Host = "www." + host
	return MapResult{
		Destination: loc.String(),
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

func mapLookup(r *http.Request) MapResult {

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
	return MapResult{
		Destination: result.String(),
		StatusCode:  http.StatusTemporaryRedirect,
	}
}

func mapRemoveWww(r *http.Request) MapResult {

	host := r.Host

	if !strings.HasPrefix(host, "www.") {
		return makeError("error-does-not-have-www")
	}

	loc := url.URL(*r.URL)
	loc.Scheme = getScheme(r)
	loc.Host = host[4:]
	return MapResult{
		Destination: loc.String(),
		StatusCode:  http.StatusTemporaryRedirect,
	}
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
	var mapFn func(r *http.Request) MapResult

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
		result := mapFn(r)

		if r.Header.Get("X-Redirect2Me-Debug") == "1" {
			HandleJson(w, r, result)
			return
		}

		http.Redirect(w, r, result.Destination, result.StatusCode)
	}, nil
}
