package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kingqaquuu/SubPilot/apps/server/internal/config"
	"go.uber.org/zap"
)

func TestHealthEndpoint(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Env:  "test",
			Name: "SubPilot",
			Port: "8080",
		},
		Server: config.ServerConfig{
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		API: config.APIConfig{
			Prefix: "/api/v1",
		},
	}

	engine := New(cfg, zap.NewNop())
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Service string `json:"service"`
			Status  string `json:"status"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Code != 0 || body.Message != "success" {
		t.Fatalf("unexpected response envelope: %+v", body)
	}
	if body.Data.Service != "SubPilot" || body.Data.Status != "ok" {
		t.Fatalf("unexpected health payload: %+v", body.Data)
	}
}
