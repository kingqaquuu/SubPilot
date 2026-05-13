package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/auth"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/dto"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/response"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/service"
)

type CategoryHandler struct {
	categories *service.CategoryService
}

func NewCategoryHandler(categories *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categories: categories}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_request", "invalid category request")
		return
	}

	category, err := h.categories.Create(c.Request.Context(), userID, req)
	if err != nil {
		writeCategoryError(c, err, "create_category_failed")
		return
	}

	response.Success(c, category)
}

func (h *CategoryHandler) List(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	categories, err := h.categories.List(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "list_categories_failed", "list categories failed")
		return
	}

	response.Success(c, categories)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	id, ok := categoryID(c)
	if !ok {
		return
	}

	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_request", "invalid category request")
		return
	}

	category, err := h.categories.Update(c.Request.Context(), userID, id, req)
	if err != nil {
		writeCategoryError(c, err, "update_category_failed")
		return
	}

	response.Success(c, category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}

	id, ok := categoryID(c)
	if !ok {
		return
	}

	category, err := h.categories.Delete(c.Request.Context(), userID, id)
	if err != nil {
		writeCategoryError(c, err, "delete_category_failed")
		return
	}

	response.Success(c, category)
}

func authenticatedUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, err := auth.UserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "missing_user", "missing authenticated user")
		return uuid.Nil, false
	}

	return userID, true
}

func categoryID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid_category_id", "invalid category id")
		return uuid.Nil, false
	}

	return id, true
}

func writeCategoryError(c *gin.Context, err error, fallbackCode string) {
	switch {
	case errors.Is(err, service.ErrInvalidCategory):
		response.Error(c, http.StatusBadRequest, "invalid_request", "invalid category request")
	case errors.Is(err, service.ErrCategoryNameExists):
		response.Error(c, http.StatusConflict, "category_name_exists", "category name already exists")
	case errors.Is(err, service.ErrCategoryNotFound):
		response.Error(c, http.StatusNotFound, "category_not_found", "category not found")
	default:
		response.Error(c, http.StatusInternalServerError, fallbackCode, "category operation failed")
	}
}
