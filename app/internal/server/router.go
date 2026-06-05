package server

import (
	"context"
	"log/slog"
	"net/http"
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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /ready", s.handleReady)
	mux.HandleFunc("POST /recommend", s.handleRecommend)

	return loggingMiddleware(logger, mux)
}
