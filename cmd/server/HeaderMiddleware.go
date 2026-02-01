package main

import (
	"net/http"
)

func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "redirect2.me")
		next.ServeHTTP(w, r)
	})
}
