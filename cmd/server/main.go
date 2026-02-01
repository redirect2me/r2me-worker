package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/redirect2me/r2me-worker/ui"
)

func addAdminRoutes(mux *http.ServeMux, hostname string) {
	staticHandler := ui.GetStaticHandler(Logger.Logger)

	mux.Handle(hostname+"/favicon.ico", staticHandler)
	mux.Handle(hostname+"/favicon.svg", staticHandler)
	mux.Handle(hostname+"/robots.txt", staticHandler)
	mux.Handle(hostname+"/css/", staticHandler)
	mux.Handle(hostname+"/js/", staticHandler)
	mux.Handle(hostname+"/images/", staticHandler)
	mux.HandleFunc(hostname+"/status.json", StatusHandler)
	mux.HandleFunc(hostname+"/{$}", RootHandler)
}

func GetPublicIP() (string, error) {
	resp, err := http.Get("https://resolve.rs/ip/whatsmyip.txt")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(ip)), nil
}

func main() {

	Logger.Trace("Server starting", "config", Config)
	if Config.LoadError != nil {
		Logger.Error("Error loading config", "error", Config.LoadError)
	}

	mux := http.NewServeMux()

	if Config.AdminHost != "none" && Config.AdminHost != "" {
		addAdminRoutes(mux, Config.AdminHost)
	}

	if Config.AdminIP == "none" || Config.AdminIP == "" {
		// do nothing
	} else if Config.AdminIP == "auto" {
		autoIp, autoIpErr := GetPublicIP()
		if autoIpErr != nil {
			Logger.Error("unable to determine public IP for admin_ip", "error", autoIpErr)
		} else {
			Logger.Debug("adding admin routes for auto-detected public IP", "ip", autoIp)
			addAdminRoutes(mux, autoIp)
		}
	} else {
		adminIp, adminIpErr := netip.ParseAddr(Config.AdminIP)
		if adminIpErr != nil {
			Logger.Error("invalid admin_ip config value", "error", adminIpErr)
		} else {
			Logger.Debug("adding admin routes for specific IP", "ip", adminIp.String())
			addAdminRoutes(mux, adminIp.String())
		}
	}

	mapper, mapErr := GetMapper(Config.Action)
	if mapErr != nil {
		Logger.Error("unable to get mapper", "error", mapErr)
		return
	}
	mux.HandleFunc("/", mapper)

	handler := RecoveryMiddleware(LoggingMiddleware(HeaderMiddleware(mux)))

	var quit = make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var httpSrv, httpsSrv *http.Server

	if Config.HttpPort > 0 {
		httpSrv = HttpServer(fmt.Sprintf("%s:%d", Config.HttpAddr, Config.HttpPort), handler)
	}
	if Config.HttpsPort > 0 {
		httpsSrv = HttpsServer(fmt.Sprintf("%s:%d", Config.HttpsAddr, Config.HttpsPort), handler)
	}

	<-quit

	Logger.Info("Starting graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if httpSrv != nil {
		httpSrv.SetKeepAlivesEnabled(false)
		if err := httpSrv.Shutdown(ctx); err != nil {
			Logger.Error("HTTP server force shutdown", "error", err)
		}
	}
	if httpsSrv != nil {
		httpsSrv.SetKeepAlivesEnabled(false)
		if err := httpsSrv.Shutdown(ctx); err != nil {
			Logger.Error("HTTPS server force shutdown", "error", err)
		}
	}

	Logger.Info("Server exiting")
}
