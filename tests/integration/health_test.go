package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ramdhanrizkij/next-store-api/internal/config"
	"github.com/ramdhanrizkij/next-store-api/internal/server"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/response"
)

func TestHealthEndpoints(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Env:  "test",
			Port: "8080",
		},
		JWT: config.JWTConfig{
			Secret:          "testsecret",
			ExpirationHours: 1,
		},
	}

	router := server.NewRouter(cfg, nil)

	t.Run("GET /ping should return pong", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/ping", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp response.Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if !resp.Success {
			t.Errorf("expected success to be true, got %v", resp.Success)
		}
	})

	t.Run("GET /health should return healthy status", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp response.Response
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if !resp.Success {
			t.Errorf("expected success to be true, got %v", resp.Success)
		}
	})
}
