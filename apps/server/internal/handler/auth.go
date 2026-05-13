package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/auth"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/dto"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/response"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_request", "invalid register request")
		return
	}

	authResponse, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			response.Error(c, http.StatusConflict, "email_exists", "email already exists")
			return
		}
		response.Error(c, http.StatusInternalServerError, "register_failed", "register failed")
		return
	}

	response.Success(c, authResponse)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_request", "invalid login request")
		return
	}

	authResponse, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Error(c, http.StatusUnauthorized, "invalid_credentials", "invalid credentials")
			return
		}
		response.Error(c, http.StatusInternalServerError, "login_failed", "login failed")
		return
	}

	response.Success(c, authResponse)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, err := auth.UserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "missing_user", "missing authenticated user")
		return
	}

	user, err := h.authService.CurrentUser(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "user_not_found", "user not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "current_user_failed", "current user lookup failed")
		return
	}

	response.Success(c, user)
}
