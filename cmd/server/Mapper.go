package main

import (
	"errors"
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
