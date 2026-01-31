package main

import (
	"net/http"
)

func HttpServer(httpAddress string, mux http.Handler) {
	httpListenErr := http.ListenAndServe(httpAddress, mux)
	if httpListenErr != nil {
		Logger.Error("unable to listen", "address", httpAddress, "error", httpListenErr)
	}

}
