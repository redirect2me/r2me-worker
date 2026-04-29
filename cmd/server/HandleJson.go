package main

import (
	"encoding/json"
	"net/http"
)

func HandleJson(w http.ResponseWriter, r *http.Request, data any) {

	var jsonData []byte
	var jsonErr error
	if r.FormValue("pretty") != "" {
		jsonData, jsonErr = json.MarshalIndent(data, "", "  ")
	} else {
		jsonData, jsonErr = json.Marshal(data)
	}

	if jsonErr != nil {
		Logger.Error("json.Marshal failed", "error", jsonErr, "data", data)
		jsonData = []byte("{\"success\":false,\"err\":\"json.Marshal failed\"}")
	}

	var callback = r.FormValue("callback")
	if callback != "" {
		w.Header().Set("Content-Type", "application/javascript; charset=utf8")
		w.Write([]byte(callback + "("))
		w.Write(jsonData)
		w.Write([]byte(");"))
	} else {
		//w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.Header().Set("Content-Type", "text/plain; charset=utf8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		w.Header().Set("Access-Control-Max-Age", "604800") // 1 week
		w.Write(jsonData)
	}
}
