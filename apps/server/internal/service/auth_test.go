package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/auth"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/dto"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestAuthServiceRegisterAndLogin(t *testing.T) {
	users := newFakeUserRepository()
	tokens := newTestTokenManager(t)
	service := NewAuthService(users, tokens)

	registered, err := service.Register(context.Background(), dto.RegisterRequest{
		Email:    " User@Example.com ",
		Password: "password123",
		Name:     "User",
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if registered.User.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %s", registered.User.Email)
	}
	if registered.AccessToken == "" || registered.TokenType != "Bearer" {
		t.Fatalf("unexpected auth response: %+v", registered)
	}

	stored := users.byEmail["user@example.com"]
	if stored.PasswordHash == "password123" {
		t.Fatal("password stored as plaintext")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(stored.PasswordHash), []byte("password123")); err != nil {
		t.Fatalf("stored hash does not match password: %v", err)
	}

	loggedIn, err := service.Login(context.Background(), dto.LoginRequest{
		Email:    "user@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if loggedIn.AccessToken == "" {
		t.Fatal("expected login access token")
	}
}

func TestAuthServiceRejectsDuplicateEmail(t *testing.T) {
	users := newFakeUserRepository()
	service := NewAuthService(users, newTestTokenManager(t))

	req := dto.RegisterRequest{Email: "user@example.com", Password: "password123", Name: "User"}
	if _, err := service.Register(context.Background(), req); err != nil {
		t.Fatalf("register first user: %v", err)
	}
	if _, err := service.Register(context.Background(), req); !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("expected duplicate email error, got %v", err)
	}
}

func TestAuthServiceRejectsInvalidCredentials(t *testing.T) {
	users := newFakeUserRepository()
	service := NewAuthService(users, newTestTokenManager(t))

	if _, err := service.Login(context.Background(), dto.LoginRequest{Email: "missing@example.com", Password: "password123"}); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials for missing user, got %v", err)
	}

	_, err := service.Register(context.Background(), dto.RegisterRequest{Email: "user@example.com", Password: "password123", Name: "User"})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, err := service.Login(context.Background(), dto.LoginRequest{Email: "user@example.com", Password: "wrongpassword"}); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials for bad password, got %v", err)
	}
}

func TestAuthServiceCurrentUser(t *testing.T) {
	users := newFakeUserRepository()
	service := NewAuthService(users, newTestTokenManager(t))

	registered, err := service.Register(context.Background(), dto.RegisterRequest{
		Email:    "user@example.com",
		Password: "password123",
		Name:     "User",
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	user, err := service.CurrentUser(context.Background(), registered.User.ID)
	if err != nil {
		t.Fatalf("current user: %v", err)
	}
	if user.Email != "user@example.com" {
		t.Fatalf("unexpected current user: %+v", user)
	}
}

func newTestTokenManager(t *testing.T) *auth.TokenManager {
	t.Helper()

	manager, err := auth.NewTokenManager("test-secret", time.Hour)
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}

	return manager
}

type fakeUserRepository struct {
	byID    map[uuid.UUID]*model.User
	byEmail map[string]*model.User
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		byID:    map[uuid.UUID]*model.User{},
		byEmail: map[string]*model.User{},
	}
}

func (r *fakeUserRepository) Create(ctx context.Context, user *model.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	copy := *user
	r.byID[user.ID] = &copy
	r.byEmail[user.Email] = &copy
	return nil
}

func (r *fakeUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copy := *user
	return &copy, nil
}

func (r *fakeUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user, ok := r.byEmail[email]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copy := *user
	return &copy, nil
}

func (r *fakeUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	_, ok := r.byEmail[email]
	return ok, nil
}
