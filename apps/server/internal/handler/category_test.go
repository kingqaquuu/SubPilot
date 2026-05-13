package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/auth"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/dto"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/middleware"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/model"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/repository"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/service"
	"gorm.io/gorm"
)

func TestCategoryHandlerCreateListUpdateAndDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, token := newCategoryTestRouter(t)

	createRec := performJSON(router, http.MethodPost, "/categories", `{"name":"Streaming","color":"#0ea5e9"}`, token)
	if createRec.Code != http.StatusOK {
		t.Fatalf("expected create status %d, got %d: %s", http.StatusOK, createRec.Code, createRec.Body.String())
	}
	var createResp struct {
		Data dto.CategoryResponse `json:"data"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if createResp.Data.ID == uuid.Nil {
		t.Fatal("expected category id")
	}

	listRec := performJSON(router, http.MethodGet, "/categories", "", token)
	if listRec.Code != http.StatusOK {
		t.Fatalf("expected list status %d, got %d: %s", http.StatusOK, listRec.Code, listRec.Body.String())
	}
	var listResp struct {
		Data []dto.CategoryResponse `json:"data"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(listResp.Data) != 1 || listResp.Data[0].ID != createResp.Data.ID {
		t.Fatalf("unexpected list response: %+v", listResp.Data)
	}

	updateRec := performJSON(router, http.MethodPut, "/categories/"+createResp.Data.ID.String(), `{"name":"Music","color":"#a855f7"}`, token)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected update status %d, got %d: %s", http.StatusOK, updateRec.Code, updateRec.Body.String())
	}

	deleteRec := performJSON(router, http.MethodDelete, "/categories/"+createResp.Data.ID.String(), "", token)
	if deleteRec.Code != http.StatusOK {
		t.Fatalf("expected delete status %d, got %d: %s", http.StatusOK, deleteRec.Code, deleteRec.Body.String())
	}
}

func TestCategoryHandlerRejectsMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, _ := newCategoryTestRouter(t)

	rec := performJSON(router, http.MethodGet, "/categories", "", "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestCategoryHandlerRejectsDuplicateName(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, token := newCategoryTestRouter(t)

	if rec := performJSON(router, http.MethodPost, "/categories", `{"name":"Streaming"}`, token); rec.Code != http.StatusOK {
		t.Fatalf("expected first create status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	rec := performJSON(router, http.MethodPost, "/categories", `{"name":"Streaming"}`, token)
	if rec.Code != http.StatusConflict {
		t.Fatalf("expected duplicate status %d, got %d: %s", http.StatusConflict, rec.Code, rec.Body.String())
	}
}

func TestCategoryHandlerRejectsInvalidCategoryID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, token := newCategoryTestRouter(t)

	rec := performJSON(router, http.MethodPut, "/categories/not-a-uuid", `{"name":"Music"}`, token)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid id status %d, got %d: %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}

func newCategoryTestRouter(t *testing.T) (*gin.Engine, string) {
	t.Helper()

	tokenManager, err := auth.NewTokenManager("test-secret", time.Hour)
	if err != nil {
		t.Fatalf("new token manager: %v", err)
	}
	userID := uuid.New()
	token, err := tokenManager.Generate(userID)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	categoryService := service.NewCategoryService(newHandlerFakeCategoryRepository())
	categoryHandler := NewCategoryHandler(categoryService)

	router := gin.New()
	categories := router.Group("/categories", middleware.Auth(tokenManager))
	categories.POST("", categoryHandler.Create)
	categories.GET("", categoryHandler.List)
	categories.PUT("/:id", categoryHandler.Update)
	categories.DELETE("/:id", categoryHandler.Delete)

	return router, token
}

type handlerFakeCategoryRepository struct {
	byID map[uuid.UUID]*model.Category
}

func newHandlerFakeCategoryRepository() *handlerFakeCategoryRepository {
	return &handlerFakeCategoryRepository{byID: map[uuid.UUID]*model.Category{}}
}

func (r *handlerFakeCategoryRepository) Create(ctx context.Context, category *model.Category) error {
	for _, existing := range r.byID {
		if existing.UserID == category.UserID && existing.Name == category.Name {
			return repository.ErrDuplicateCategoryName
		}
	}
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	now := time.Now().UTC()
	category.CreatedAt = now
	category.UpdatedAt = now
	copy := *category
	r.byID[category.ID] = &copy
	return nil
}

func (r *handlerFakeCategoryRepository) FindByIDForUser(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Category, error) {
	category, ok := r.byID[id]
	if !ok || category.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	copy := *category
	return &copy, nil
}

func (r *handlerFakeCategoryRepository) FindByNameForUser(ctx context.Context, name string, userID uuid.UUID) (*model.Category, error) {
	for _, category := range r.byID {
		if category.UserID == userID && category.Name == name {
			copy := *category
			return &copy, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *handlerFakeCategoryRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Category, error) {
	categories := []model.Category{}
	for _, category := range r.byID {
		if category.UserID == userID {
			categories = append(categories, *category)
		}
	}
	return categories, nil
}

func (r *handlerFakeCategoryRepository) Update(ctx context.Context, category *model.Category) error {
	for id, existing := range r.byID {
		if id != category.ID && existing.UserID == category.UserID && existing.Name == category.Name {
			return repository.ErrDuplicateCategoryName
		}
	}
	copy := *category
	copy.UpdatedAt = time.Now().UTC()
	r.byID[category.ID] = &copy
	return nil
}

func (r *handlerFakeCategoryRepository) Delete(ctx context.Context, category *model.Category) error {
	delete(r.byID, category.ID)
	return nil
}
