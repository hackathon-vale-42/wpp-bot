package api

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

func promMiddleware(next handlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(r.Method, r.URL.Path))
		defer timer.ObserveDuration()

		status, _ := next(w, r)

		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(status)).Inc()
	})
}
