package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadReadsEnvironmentVariables(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_NAME", "TestPilot")
	t.Setenv("APP_PORT", "18080")
	t.Setenv("API_PREFIX", "/custom")
	t.Setenv("SERVER_READ_TIMEOUT", "3s")
	t.Setenv("SERVER_WRITE_TIMEOUT", "4s")
	t.Setenv("POSTGRES_HOST", "db")
	t.Setenv("POSTGRES_PORT", "15432")
	t.Setenv("POSTGRES_DB", "test_db")
	t.Setenv("POSTGRES_USER", "test_user")
	t.Setenv("POSTGRES_PASSWORD", "test_password")
	t.Setenv("POSTGRES_SSL_MODE", "require")
	t.Setenv("REDIS_HOST", "cache")
	t.Setenv("REDIS_PORT", "16379")
	t.Setenv("REDIS_PASSWORD", "redis_password")
	t.Setenv("REDIS_DB", "2")
	t.Setenv("JWT_SECRET", "test-secret")
	t.Setenv("JWT_EXPIRES_IN", "12h")

	cfg := loadFromTempDir(t)

	if cfg.App.Env != "test" || cfg.App.Name != "TestPilot" || cfg.App.Port != "18080" {
		t.Fatalf("unexpected app config: %+v", cfg.App)
	}
	if cfg.API.Prefix != "/custom" {
		t.Fatalf("unexpected api prefix: %s", cfg.API.Prefix)
	}
	if cfg.Server.ReadTimeout != 3*time.Second || cfg.Server.WriteTimeout != 4*time.Second {
		t.Fatalf("unexpected server timeouts: %+v", cfg.Server)
	}
	if cfg.Postgres.Host != "db" || cfg.Postgres.Port != "15432" || cfg.Postgres.Database != "test_db" || cfg.Postgres.User != "test_user" || cfg.Postgres.Password != "test_password" || cfg.Postgres.SSLMode != "require" {
		t.Fatalf("unexpected postgres config: %+v", cfg.Postgres)
	}
	if cfg.Redis.Host != "cache" || cfg.Redis.Port != "16379" || cfg.Redis.Password != "redis_password" || cfg.Redis.DB != 2 {
		t.Fatalf("unexpected redis config: %+v", cfg.Redis)
	}
	if cfg.JWT.Secret != "test-secret" || cfg.JWT.ExpiresIn != 12*time.Hour {
		t.Fatalf("unexpected jwt config: %+v", cfg.JWT)
	}
}

func TestLoadRejectsDefaultProductionSecrets(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("POSTGRES_PASSWORD", "change-me")
	t.Setenv("JWT_SECRET", "change-me-to-a-long-random-secret")

	_, err := loadFromTempDirWithError(t)
	if err == nil {
		t.Fatal("expected production defaults to be rejected")
	}
}

func TestLoadRejectsDefaultJWTSecretOutsideProduction(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("JWT_SECRET", "change-me-to-a-long-random-secret")

	_, err := loadFromTempDirWithError(t)
	if err == nil {
		t.Fatal("expected default jwt secret to be rejected")
	}
}

func TestLoadAllowsProductionWithNonDefaultSecrets(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("POSTGRES_PASSWORD", "prod-postgres-password")
	t.Setenv("JWT_SECRET", "prod-jwt-secret-with-enough-entropy")

	cfg := loadFromTempDir(t)

	if cfg.App.Env != "production" {
		t.Fatalf("unexpected app env: %s", cfg.App.Env)
	}
}

func loadFromTempDir(t *testing.T) *Config {
	t.Helper()

	cfg, err := loadFromTempDirWithError(t)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	return cfg
}

func loadFromTempDirWithError(t *testing.T) (*Config, error) {
	t.Helper()

	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}

	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previousDir); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})

	return Load()
}
