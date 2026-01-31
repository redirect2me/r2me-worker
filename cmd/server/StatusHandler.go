package main

import (
	"net/http"
	"runtime"
	"time"
)

var COMMIT string
var LASTMOD string

type Status struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Commit    string `json:"commit"`
	LastMod   string `json:"lastmod"`
	Timestamp string `json:"timestamp"`
	Tech      string `json:"tech"`
	Version   string `json:"version"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	status := Status{}

	status.Success = true
	status.Message = "OK"
	status.Timestamp = time.Now().UTC().Format(time.RFC3339)
	status.Commit = COMMIT
	status.LastMod = LASTMOD
	status.Tech = runtime.Version()

	status.Version = runtime.Version()

	HandleJson(w, r, status)
}
