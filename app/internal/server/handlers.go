package server

import (
	"encoding/json"
	"net/http"
)

type recommendRequest struct {
	UserID string `json:"user_id"`
}

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("OK"))
}

func (s *server) handleReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("OK"))
}

func (s *server) handleRecommend(w http.ResponseWriter, r *http.Request) {
	var request recommendRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.logger.Error("failed to decode request", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recommendations, err := s.rec.Recommend(r.Context(), request.UserID)
	if err != nil {
		s.logger.Error("failed to recommend", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(recommendations)
}
