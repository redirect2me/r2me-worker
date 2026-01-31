package main

import (
	"errors"
	"net/http"
)

func HttpServer(httpAddress string, mux http.Handler) *http.Server {
	srv := &http.Server{
		Addr:    httpAddress,
		Handler: mux,
	}

	go func() {
		httpListenErr := srv.ListenAndServe()
		if httpListenErr != nil && !errors.Is(httpListenErr, http.ErrServerClosed) {
			Logger.Error("unable to listen", "address", httpAddress, "error", httpListenErr)
		}
	}()

	return srv
}
