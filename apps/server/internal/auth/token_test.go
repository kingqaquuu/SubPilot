package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTokenManagerGenerateAndParse(t *testing.T) {
	manager, err := NewTokenManager("test-secret", time.Hour)
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}

	userID := uuid.New()
	token, err := manager.Generate(userID)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	got, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if got != userID {
		t.Fatalf("expected user id %s, got %s", userID, got)
	}
}

func TestTokenManagerRejectsInvalidToken(t *testing.T) {
	manager, err := NewTokenManager("test-secret", time.Hour)
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}

	if _, err := manager.Parse("not-a-token"); err == nil {
		t.Fatal("expected invalid token error")
	}
}

func TestTokenManagerRejectsExpiredToken(t *testing.T) {
	now := time.Date(2026, 5, 13, 10, 0, 0, 0, time.UTC)
	manager, err := NewTokenManager("test-secret", time.Hour)
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}
	manager.now = func() time.Time { return now }

	token, err := manager.Generate(uuid.New())
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	manager.now = func() time.Time { return now.Add(2 * time.Hour) }
	if _, err := manager.Parse(token); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected invalid token error for expired token, got %v", err)
	}
}

func TestNewTokenManagerRequiresSecret(t *testing.T) {
	if _, err := NewTokenManager("", time.Hour); err == nil {
		t.Fatal("expected missing secret error")
	}
}
