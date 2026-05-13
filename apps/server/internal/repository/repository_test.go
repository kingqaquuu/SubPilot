package repository

import (
	"testing"

	"gorm.io/gorm"
)

func TestNewRepositories(t *testing.T) {
	repos := New(&gorm.DB{})

	if repos.Users == nil {
		t.Fatal("expected user repository")
	}
	if repos.Categories == nil {
		t.Fatal("expected category repository")
	}
	if repos.Subscriptions == nil {
		t.Fatal("expected subscription repository")
	}
	if repos.Reminders == nil {
		t.Fatal("expected reminder repository")
	}
	if repos.PaymentRecords == nil {
		t.Fatal("expected payment record repository")
	}
}
