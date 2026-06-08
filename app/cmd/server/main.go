package main

import (
	"log/slog"
	"os"

	"github.com/edofo/ml-sre-train/internal/recommender"
	"github.com/edofo/ml-sre-train/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		logger.Error("PORT environment variable is not set")
		os.Exit(1)
	}

	rec := recommender.Stub{}

	err := server.RunServer(port, &rec, logger)
	if err != nil {
		logger.Error("failed to run server", "error", err)
		os.Exit(1)
	}
}
