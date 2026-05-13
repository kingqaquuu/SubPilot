package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/auth"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/dto"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/model"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService struct {
	users  repository.UserRepository
	tokens *auth.TokenManager
}

func NewAuthService(users repository.UserRepository, tokens *auth.TokenManager) *AuthService {
	return &AuthService{
		users:  users,
		tokens: tokens,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	email := normalizeEmail(req.Email)
	exists, err := s.users.EmailExists(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("check email: %w", err)
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		Email:        email,
		PasswordHash: string(hash),
		Name:         strings.TrimSpace(req.Name),
	}
	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return s.authResponse(user)
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.users.FindByEmail(ctx, normalizeEmail(req.Email))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.authResponse(user)
}

func (s *AuthService) CurrentUser(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("find current user: %w", err)
	}

	resp := userResponse(user)
	return &resp, nil
}

func (s *AuthService) authResponse(user *model.User) (*dto.AuthResponse, error) {
	token, err := s.tokens.Generate(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.tokens.ExpiresInSeconds(),
		User:        userResponse(user),
	}, nil
}

func userResponse(user *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
