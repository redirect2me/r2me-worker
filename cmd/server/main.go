package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/redirect2me/r2me-worker/ui"
)

func main() {

	Logger.Trace("Starting", "timestamp", time.Now().UTC().Format(time.RFC3339)) //, "config", Config)

	staticHandler := ui.GetStaticHandler(Logger.Logger)

	mux := http.NewServeMux()

	mux.Handle(Config.Hostname+"/favicon.ico", staticHandler)
	mux.Handle(Config.Hostname+"/favicon.svg", staticHandler)
	mux.Handle(Config.Hostname+"/robots.txt", staticHandler)
	mux.Handle(Config.Hostname+"/css/", staticHandler)
	mux.Handle(Config.Hostname+"/js/", staticHandler)
	mux.Handle(Config.Hostname+"/images/", staticHandler)
	mux.HandleFunc(Config.Hostname+"/status.json", StatusHandler)
	mux.HandleFunc(Config.Hostname+"/{$}", RootHandler)

	mapper, mapErr := GetMapper(Config.Action)
	if mapErr != nil {
		Logger.Error("unable to get mapper", "error", mapErr)
		return
	}
	mux.HandleFunc("/", mapper)

	handler := RecoveryMiddleware(LoggingMiddleware(mux))

	var done = make(chan bool)

	if Config.HttpPort > 0 {
		go HttpServer(fmt.Sprintf("%s:%d", Config.HttpHost, Config.HttpPort), handler)
	}
	if Config.HttpsPort > 0 {
		go HttpsServer(fmt.Sprintf("%s:%d", Config.HttpsHost, Config.HttpsPort), handler)
	}

	<-done
}
