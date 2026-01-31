package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redirect2me/r2me-worker/ui"
)

func main() {

	Logger.Trace("Server starting") //, "config", Config)

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

	var quit = make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var httpSrv *http.Server

	if Config.HttpPort > 0 {
		httpSrv = HttpServer(fmt.Sprintf("%s:%d", Config.HttpHost, Config.HttpPort), handler)
	}
	if Config.HttpsPort > 0 {
		HttpsServer(fmt.Sprintf("%s:%d", Config.HttpsHost, Config.HttpsPort), handler)
	}

	<-quit

	Logger.Info("Starting graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if httpSrv != nil {
		if err := httpSrv.Shutdown(ctx); err != nil {
			Logger.Error("Server forced to shutdown", "error", err)
		}
	}

	Logger.Info("Server exiting")
}
