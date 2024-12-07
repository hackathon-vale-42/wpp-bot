package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(listenAddr string) error {
	http.Handle("/metrics", promhttp.Handler())

	slog.Info("Server started", "ListenAddr", listenAddr)

	return http.ListenAndServe(listenAddr, nil)
}

func writeJson(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("writeJson: failed to encode")
		panic(err)
	}
}
