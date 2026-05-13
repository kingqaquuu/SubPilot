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

func TestSubscriptionHandlerCreateListDetailUpdateAndDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, token := newSubscriptionTestRouter(t)

	createRec := performJSON(router, http.MethodPost, "/subscriptions", `{"name":"Netflix","amount":15.99,"currency":"USD","billing_cycle":"monthly","next_payment_date":"2026-06-01"}`, token)
	if createRec.Code != http.StatusOK {
		t.Fatalf("expected create status %d, got %d: %s", http.StatusOK, createRec.Code, createRec.Body.String())
	}
	var createResp struct {
		Data dto.SubscriptionResponse `json:"data"`
	}
	if err := json.Unmarshal(createRec.Body.Bytes(), &createResp); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if createResp.Data.ID == uuid.Nil {
		t.Fatal("expected subscription id")
	}

	listRec := performJSON(router, http.MethodGet, "/subscriptions?page=1&page_size=10&status=active", "", token)
	if listRec.Code != http.StatusOK {
		t.Fatalf("expected list status %d, got %d: %s", http.StatusOK, listRec.Code, listRec.Body.String())
	}

	detailRec := performJSON(router, http.MethodGet, "/subscriptions/"+createResp.Data.ID.String(), "", token)
	if detailRec.Code != http.StatusOK {
		t.Fatalf("expected detail status %d, got %d: %s", http.StatusOK, detailRec.Code, detailRec.Body.String())
	}

	updateRec := performJSON(router, http.MethodPut, "/subscriptions/"+createResp.Data.ID.String(), `{"name":"Netflix Premium","amount":19.99,"currency":"USD","billing_cycle":"yearly","next_payment_date":"2026-07-01","status":"paused"}`, token)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("expected update status %d, got %d: %s", http.StatusOK, updateRec.Code, updateRec.Body.String())
	}

	deleteRec := performJSON(router, http.MethodDelete, "/subscriptions/"+createResp.Data.ID.String(), "", token)
	if deleteRec.Code != http.StatusOK {
		t.Fatalf("expected delete status %d, got %d: %s", http.StatusOK, deleteRec.Code, deleteRec.Body.String())
	}
}

func TestSubscriptionHandlerRejectsMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, _ := newSubscriptionTestRouter(t)

	rec := performJSON(router, http.MethodGet, "/subscriptions", "", "")
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestSubscriptionHandlerRejectsInvalidSubscriptionID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, token := newSubscriptionTestRouter(t)

	rec := performJSON(router, http.MethodGet, "/subscriptions/not-a-uuid", "", token)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid id status %d, got %d: %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}

func TestSubscriptionHandlerRejectsInvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, token := newSubscriptionTestRouter(t)

	rec := performJSON(router, http.MethodPost, "/subscriptions", `{"name":"","amount":0,"currency":"USD","billing_cycle":"monthly","next_payment_date":"2026-06-01"}`, token)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid request status %d, got %d: %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}
}

func newSubscriptionTestRouter(t *testing.T) (*gin.Engine, string) {
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

	subscriptionService := service.NewSubscriptionService(newHandlerFakeSubscriptionRepository(), newHandlerFakeCategoryRepository())
	subscriptionHandler := NewSubscriptionHandler(subscriptionService)

	router := gin.New()
	subscriptions := router.Group("/subscriptions", middleware.Auth(tokenManager))
	subscriptions.POST("", subscriptionHandler.Create)
	subscriptions.GET("", subscriptionHandler.List)
	subscriptions.GET("/:id", subscriptionHandler.Detail)
	subscriptions.PUT("/:id", subscriptionHandler.Update)
	subscriptions.DELETE("/:id", subscriptionHandler.Delete)

	return router, token
}

type handlerFakeSubscriptionRepository struct {
	byID map[uuid.UUID]*model.Subscription
}

func newHandlerFakeSubscriptionRepository() *handlerFakeSubscriptionRepository {
	return &handlerFakeSubscriptionRepository{byID: map[uuid.UUID]*model.Subscription{}}
}

func (r *handlerFakeSubscriptionRepository) Create(ctx context.Context, subscription *model.Subscription) error {
	if subscription.ID == uuid.Nil {
		subscription.ID = uuid.New()
	}
	now := time.Now().UTC()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now
	copy := *subscription
	r.byID[subscription.ID] = &copy
	return nil
}

func (r *handlerFakeSubscriptionRepository) FindByIDForUser(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Subscription, error) {
	subscription, ok := r.byID[id]
	if !ok || subscription.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	copy := *subscription
	return &copy, nil
}

func (r *handlerFakeSubscriptionRepository) ListByUser(ctx context.Context, filter repository.SubscriptionFilter) ([]model.Subscription, int64, error) {
	matches := []model.Subscription{}
	for _, subscription := range r.byID {
		if subscription.UserID == filter.UserID {
			matches = append(matches, *subscription)
		}
	}
	return matches, int64(len(matches)), nil
}

func (r *handlerFakeSubscriptionRepository) Update(ctx context.Context, subscription *model.Subscription) error {
	copy := *subscription
	copy.UpdatedAt = time.Now().UTC()
	r.byID[subscription.ID] = &copy
	return nil
}

func (r *handlerFakeSubscriptionRepository) Delete(ctx context.Context, subscription *model.Subscription) error {
	delete(r.byID, subscription.ID)
	return nil
}
