package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/auth"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/dto"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/middleware"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/model"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/service"
	"gorm.io/gorm"
)

func TestAuthHandlerRegisterLoginAndMe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	users := newHandlerFakeUserRepository()
	tokenManager := newHandlerTestTokenManager(t)
	authService := service.NewAuthService(users, tokenManager)
	authHandler := NewAuthHandler(authService)

	router := gin.New()
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.GET("/me", middleware.Auth(tokenManager), authHandler.Me)

	registerBody := `{"email":"user@example.com","password":"password123","name":"User"}`
	registerRec := performJSON(router, http.MethodPost, "/register", registerBody, "")
	if registerRec.Code != http.StatusOK {
		t.Fatalf("expected register status %d, got %d: %s", http.StatusOK, registerRec.Code, registerRec.Body.String())
	}
	var registerResp struct {
		Data dto.AuthResponse `json:"data"`
	}
	if err := json.Unmarshal(registerRec.Body.Bytes(), &registerResp); err != nil {
		t.Fatalf("decode register response: %v", err)
	}
	if registerResp.Data.AccessToken == "" {
		t.Fatal("expected register access token")
	}

	loginBody := `{"email":"user@example.com","password":"password123"}`
	loginRec := performJSON(router, http.MethodPost, "/login", loginBody, "")
	if loginRec.Code != http.StatusOK {
		t.Fatalf("expected login status %d, got %d: %s", http.StatusOK, loginRec.Code, loginRec.Body.String())
	}
	var loginResp struct {
		Data dto.AuthResponse `json:"data"`
	}
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}

	meRec := performJSON(router, http.MethodGet, "/me", "", loginResp.Data.AccessToken)
	if meRec.Code != http.StatusOK {
		t.Fatalf("expected me status %d, got %d: %s", http.StatusOK, meRec.Code, meRec.Body.String())
	}
}

func TestAuthHandlerRejectsMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	users := newHandlerFakeUserRepository()
	tokenManager := newHandlerTestTokenManager(t)
	authHandler := NewAuthHandler(service.NewAuthService(users, tokenManager))

	router := gin.New()
	router.GET("/me", middleware.Auth(tokenManager), authHandler.Me)

	rec := performJSON(router, http.MethodGet, "/me", "", "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func performJSON(router *gin.Engine, method string, path string, body string, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func newHandlerTestTokenManager(t *testing.T) *auth.TokenManager {
	t.Helper()

	tokenManager, err := auth.NewTokenManager("test-secret", time.Hour)
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}

	return tokenManager
}

type handlerFakeUserRepository struct {
	byID    map[uuid.UUID]*model.User
	byEmail map[string]*model.User
}

func newHandlerFakeUserRepository() *handlerFakeUserRepository {
	return &handlerFakeUserRepository{
		byID:    map[uuid.UUID]*model.User{},
		byEmail: map[string]*model.User{},
	}
}

func (r *handlerFakeUserRepository) Create(ctx context.Context, user *model.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	copy := *user
	r.byID[user.ID] = &copy
	r.byEmail[user.Email] = &copy
	return nil
}

func (r *handlerFakeUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, ok := r.byID[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copy := *user
	return &copy, nil
}

func (r *handlerFakeUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user, ok := r.byEmail[email]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	copy := *user
	return &copy, nil
}

func (r *handlerFakeUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	_, ok := r.byEmail[email]
	return ok, nil
}
