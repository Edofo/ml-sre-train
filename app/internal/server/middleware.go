package server

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

type metrics struct {
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

func newMetrics(reg prometheus.Registerer) *metrics {
	factory := promauto.With(reg)
	return &metrics{
		requests: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests by method, path and status.",
			},
			[]string{"method", "path", "status"},
		),
		duration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latency in seconds by method and path.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}
}

func instrument(logger *slog.Logger, m *metrics, next http.Handler) http.Handler {
	skip := map[string]bool{"/health": true, "/ready": true, "/metrics": true}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if skip[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		sr := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sr, r)
		elapsed := time.Since(start)

		m.requests.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(sr.status)).Inc()
		m.duration.WithLabelValues(r.Method, r.URL.Path).Observe(elapsed.Seconds())

		logger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", sr.status,
			"duration", elapsed,
		)
	})
}
