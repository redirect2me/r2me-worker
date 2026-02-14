package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func RealtimeSend(data any) {
	if Config.RealtimeEndpoint == "" {
		return
	}

	payload, err := json.Marshal(data)
	if err != nil {
		Logger.Error("unable to marshal realtime data", "error", err)
		return
	}
	Logger.Trace("sending realtime data", "payload", string(payload))

	req, err := http.NewRequest("POST", Config.RealtimeEndpoint, bytes.NewReader(payload))
	if err != nil {
		Logger.Error("unable to create realtime request", "error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "r2me-worker/1.0 (commit="+COMMIT+")")
	if Config.RealtimeCredentials != "" {
		req.Header.Set("Authorization", "Bearer "+Config.RealtimeCredentials)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Logger.Error("unable to send realtime data", "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		Logger.Error("realtime endpoint returned non-200 status", "status", resp.StatusCode, "request", req, "response", resp)
		return
	}

	Logger.Trace("realtime data sent successfully")
}
