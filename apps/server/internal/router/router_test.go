package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
			Port: "18080",
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

func TestSwaggerSpecIsServed(t *testing.T) {
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.MkdirAll(docsDir, 0o755); err != nil {
		t.Fatalf("create docs dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(docsDir, "swagger.yaml"), []byte("openapi: 3.0.3\n"), 0o644); err != nil {
		t.Fatalf("write swagger spec: %v", err)
	}

	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previousDir); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})

	cfg := &config.Config{
		App: config.AppConfig{
			Env:  "test",
			Name: "SubPilot",
			Port: "18080",
		},
		API: config.APIConfig{
			Prefix: "/api/v1",
		},
	}

	engine := New(cfg, zap.NewNop())
	req := httptest.NewRequest(http.MethodGet, "/docs/swagger.yaml", nil)
	rec := httptest.NewRecorder()

	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	if rec.Body.String() != "openapi: 3.0.3\n" {
		t.Fatalf("unexpected swagger body: %q", rec.Body.String())
	}
}
