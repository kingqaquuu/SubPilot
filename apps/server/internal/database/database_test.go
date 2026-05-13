package database

import (
	"strings"
	"testing"

	"github.com/kingqaquuu/SubPilot/apps/server/internal/config"
)

func TestDSN(t *testing.T) {
	dsn := DSN(config.PostgresConfig{
		Host:     "postgres",
		Port:     "5432",
		Database: "subpilot",
		User:     "subpilot",
		Password: "secret",
		SSLMode:  "disable",
	})

	expectedParts := []string{
		"postgres://subpilot:secret@postgres:5432/subpilot",
		"TimeZone=UTC",
		"sslmode=disable",
	}

	for _, part := range expectedParts {
		if !strings.Contains(dsn, part) {
			t.Fatalf("expected DSN to contain %q, got %q", part, dsn)
		}
	}
}
