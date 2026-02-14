package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var MetricsHandler http.Handler

type MetricsData struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration prometheus.Histogram
	RedirectsTotal  *prometheus.CounterVec
}

var Metrics *MetricsData

func initMetrics() {

	reg := prometheus.NewRegistry()

	requestsTotal := promauto.With(reg).NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Tracks the number of HTTP requests.",
		}, []string{"status_code", "method"},
	)

	requestDuration := promauto.With(reg).NewHistogram(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Tracks the latencies for HTTP requests.",
			Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 5),
		},
	)

	redirectsTotal := promauto.With(reg).NewCounterVec(
		prometheus.CounterOpts{
			Name: "redirects_total",
			Help: "Tracks the number of redirects.",
		}, []string{"result_code"},
	)

	Metrics = &MetricsData{
		RequestsTotal:   requestsTotal,
		RequestDuration: requestDuration,
		RedirectsTotal:  redirectsTotal,
	}

	MetricsHandler = promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
}
