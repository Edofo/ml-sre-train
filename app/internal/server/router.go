package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type recommender interface {
	Recommend(context.Context, string) ([]string, error)
}

type server struct {
	rec    recommender
	logger *slog.Logger
}

func newRouter(rec recommender, logger *slog.Logger) http.Handler {
	s := &server{rec: rec, logger: logger}

	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	m := newMetrics(reg)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /ready", s.handleReady)
	mux.HandleFunc("POST /recommend", s.handleRecommend)
	mux.Handle("GET /metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	return instrument(logger, m, mux)
}
