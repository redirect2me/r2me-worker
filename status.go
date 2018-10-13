package main

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"time"
)


type Status struct {
	Success  bool		`json:"success"`
	Version  string		`json:"version"`
	//Environ  []string
	Getwd    string
	Hostname string
	Seconds  int64
	TempDir  string
}

func Status_handler(w http.ResponseWriter, r *http.Request) {
	var err error
	status := Status{}

	status.Getwd, err = os.Getwd()
	if err != nil {
		status.Getwd = "ERROR!"
	}

	status.Hostname, err = os.Hostname()
	if err != nil {
		status.Hostname = "ERROR"
	}

	status.TempDir = os.TempDir()
	//status.Environ = os.Environ()
	status.Version = runtime.Version()
	status.Seconds = time.Now().Unix()
	status.Success = true
	callback := r.FormValue("callback");

	w.Header().Set("Content-Type", "text/plain; charset=utf8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	w.Header().Set("Access-Control-Max-Age", "604800") // 1 week

	var b []byte
	b, err = json.Marshal(status)
	if err != nil {
		b = []byte("{\"success\":false,\"err\":\"json.Marshal failed\"}")
	}

	if callback > "" {
		w.Write([]byte(callback))
		w.Write([]byte("("))
		w.Write(b)
		w.Write([]byte(");"))
	} else {
		w.Write(b)
	}
}


