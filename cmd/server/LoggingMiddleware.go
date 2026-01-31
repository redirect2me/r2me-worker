package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type contextKey string

const responseIdKey contextKey = "response_id"

type loggingWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader overrides the default WriteHeader to record the status code.
func (rec *loggingWriter) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uuid, uuidErr := uuid.NewV7()

		var responseID string
		if uuidErr != nil {
			responseID = "resp_no_uuid_generated"
			Logger.Warn("Failed to generate UUID", "error", uuidErr)
		} else {
			responseID = fmt.Sprintf("resp_%s", strings.ReplaceAll(uuid.String(), "-", ""))
		}

		loggingW := &loggingWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		newCtx := context.WithValue(r.Context(), responseIdKey, responseID)
		rWithCtx := r.WithContext(newCtx)

		// Call the next handler in the chain
		next.ServeHTTP(loggingW, rWithCtx)

		duration := time.Since(start)
		Logger.Info("Request complete",
			"status", loggingW.statusCode,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"duration_ms", duration.Milliseconds(),
			"response_id", responseID,
		)
	})
}
