package server

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeRecommender struct {
	result []string
	err    error
}

func (f fakeRecommender) Recommend(ctx context.Context, userID string) ([]string, error) {
	return f.result, f.err
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestHealthAndReady(t *testing.T) {
	handler := newRouter(fakeRecommender{}, testLogger())

	for _, path := range []string{"/health", "/ready"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("%s: status expected 200, got %d", path, rec.Code)
		}
		if rec.Body.String() != "OK" {
			t.Errorf("%s: body expected \"OK\", got %q", path, rec.Body.String())
		}
	}
}

func TestRecommend(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		fake       fakeRecommender
		wantStatus int
		wantBody   string
	}{
		{
			name:       "valid request",
			body:       `{"user_id":"u123"}`,
			fake:       fakeRecommender{result: []string{"a", "b"}},
			wantStatus: http.StatusOK,
			wantBody:   "[\"a\",\"b\"]\n",
		},
		{
			name:       "invalid json",
			body:       `not json`,
			fake:       fakeRecommender{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "recommender error",
			body:       `{"user_id":"u123"}`,
			fake:       fakeRecommender{err: errors.New("boom")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := newRouter(tc.fake, testLogger())
			req := httptest.NewRequest(http.MethodPost, "/recommend", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status expected %d, got %d", tc.wantStatus, rec.Code)
			}
			if tc.wantBody != "" && rec.Body.String() != tc.wantBody {
				t.Errorf("body expected %q, got %q", tc.wantBody, rec.Body.String())
			}
		})
	}
}
