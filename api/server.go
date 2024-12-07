package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "route", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)
)

type handlerFunc func(http.ResponseWriter, *http.Request) (int, any)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(listenAddr string) error {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)

	http.Handle("/metrics", promhttp.Handler())

	http.Handle("/health", promMiddleware(healthHandler))

	slog.Info("Server started", "ListenAddr", listenAddr)

	return http.ListenAndServe(listenAddr, nil)
}

func writeJson(w http.ResponseWriter, status int, v any) (int, any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("writeJson: failed to encode")
		panic(err)
	}

	return status, v
}
