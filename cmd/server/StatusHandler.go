package main

import (
	"net/http"
	"runtime"
	"time"
)

var COMMIT string
var LASTMOD string
var BUILTBY string
var VERSION string

type Status struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Commit    string `json:"commit"`
	LastMod   string `json:"lastmod"`
	Timestamp string `json:"timestamp"`
	Tech      string `json:"tech"`
	BuiltBy   string `json:"built_by"`
	Version   string `json:"version"`
	Action    string `json:"action"`
	NodeID    string `json:"node_id"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	status := Status{}

	status.Success = true
	status.Message = "OK"
	status.Timestamp = time.Now().UTC().Format(time.RFC3339)
	status.Commit = COMMIT
	status.LastMod = LASTMOD
	status.Tech = runtime.Version()
	status.BuiltBy = BUILTBY
	status.Version = VERSION
	status.Action = Config.Action
	status.NodeID = Config.NodeID

	HandleJson(w, r, status)
}
