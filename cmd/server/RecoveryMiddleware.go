package main

import (
	"net/http"
	"os"
	"runtime/debug"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				Logger.Error("panic recovered in RecoveryMiddleware", "error", err, "stack", string(debug.Stack()))
				os.Stderr.WriteString(string(debug.Stack()))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
