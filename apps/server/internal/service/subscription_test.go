package service

import (
	"context"
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/dto"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/model"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/repository"
	"gorm.io/gorm"
)

func TestSubscriptionServiceCreateListDetailUpdateAndDelete(t *testing.T) {
	subscriptions := newFakeSubscriptionRepository()
	categories := newFakeCategoryRepository()
	service := NewSubscriptionService(subscriptions, categories)
	userID := uuid.New()
	category := createTestCategory(t, categories, userID, "Streaming")

	created, err := service.Create(context.Background(), userID, dto.SubscriptionRequest{
		CategoryID:      &category.ID,
		Name:            " Netflix ",
		Amount:          15.99,
		Currency:        "usd",
		BillingCycle:    "monthly",
		NextPaymentDate: "2026-06-01",
		Notes:           " Primary plan ",
	})
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}
	if created.Name != "Netflix" || created.Currency != "USD" || created.Status != "active" {
		t.Fatalf("unexpected created subscription: %+v", created)
	}
	if created.CategoryID == nil || *created.CategoryID != category.ID {
		t.Fatalf("expected category id %s, got %+v", category.ID, created.CategoryID)
	}

	listed, err := service.List(context.Background(), userID, SubscriptionQuery{Page: 1, PageSize: 10, Status: "active", BillingCycle: "monthly", CategoryID: &category.ID})
	if err != nil {
		t.Fatalf("list subscriptions: %v", err)
	}
	if listed.Total != 1 || len(listed.Items) != 1 || listed.Items[0].ID != created.ID {
		t.Fatalf("unexpected list response: %+v", listed)
	}

	detail, err := service.Detail(context.Background(), userID, created.ID)
	if err != nil {
		t.Fatalf("subscription detail: %v", err)
	}
	if detail.ID != created.ID {
		t.Fatalf("unexpected detail response: %+v", detail)
	}

	updated, err := service.Update(context.Background(), userID, created.ID, dto.SubscriptionRequest{
		Name:            "Netflix Premium",
		Amount:          19.99,
		Currency:        "USD",
		BillingCycle:    "yearly",
		NextPaymentDate: "2026-07-01",
		Status:          "paused",
	})
	if err != nil {
		t.Fatalf("update subscription: %v", err)
	}
	if updated.Name != "Netflix Premium" || updated.Status != "paused" || updated.CategoryID != nil {
		t.Fatalf("unexpected updated subscription: %+v", updated)
	}

	deleted, err := service.Delete(context.Background(), userID, created.ID)
	if err != nil {
		t.Fatalf("delete subscription: %v", err)
	}
	if deleted.ID != created.ID {
		t.Fatalf("unexpected deleted subscription: %+v", deleted)
	}
}

func TestSubscriptionServiceRejectsForeignCategory(t *testing.T) {
	subscriptions := newFakeSubscriptionRepository()
	categories := newFakeCategoryRepository()
	service := NewSubscriptionService(subscriptions, categories)
	foreignCategory := createTestCategory(t, categories, uuid.New(), "Streaming")

	_, err := service.Create(context.Background(), uuid.New(), validSubscriptionRequest(&foreignCategory.ID))
	if !errors.Is(err, ErrSubscriptionCategory) {
		t.Fatalf("expected category ownership error, got %v", err)
	}
}

func TestSubscriptionServiceRejectsCrossUserAccess(t *testing.T) {
	subscriptions := newFakeSubscriptionRepository()
	service := NewSubscriptionService(subscriptions, newFakeCategoryRepository())
	ownerID := uuid.New()
	otherUserID := uuid.New()

	created, err := service.Create(context.Background(), ownerID, validSubscriptionRequest(nil))
	if err != nil {
		t.Fatalf("create subscription: %v", err)
	}

	if _, err := service.Detail(context.Background(), otherUserID, created.ID); !errors.Is(err, ErrSubscriptionNotFound) {
		t.Fatalf("expected not found for cross-user detail, got %v", err)
	}
	if _, err := service.Update(context.Background(), otherUserID, created.ID, validSubscriptionRequest(nil)); !errors.Is(err, ErrSubscriptionNotFound) {
		t.Fatalf("expected not found for cross-user update, got %v", err)
	}
	if _, err := service.Delete(context.Background(), otherUserID, created.ID); !errors.Is(err, ErrSubscriptionNotFound) {
		t.Fatalf("expected not found for cross-user delete, got %v", err)
	}
}

func TestSubscriptionServicePaginationAndFilters(t *testing.T) {
	subscriptions := newFakeSubscriptionRepository()
	service := NewSubscriptionService(subscriptions, newFakeCategoryRepository())
	userID := uuid.New()

	firstReq := validSubscriptionRequest(nil)
	firstReq.Name = "A"
	firstReq.Status = "active"
	firstReq.BillingCycle = "monthly"
	if _, err := service.Create(context.Background(), userID, firstReq); err != nil {
		t.Fatalf("create first subscription: %v", err)
	}
	secondReq := validSubscriptionRequest(nil)
	secondReq.Name = "B"
	secondReq.Status = "paused"
	secondReq.BillingCycle = "yearly"
	if _, err := service.Create(context.Background(), userID, secondReq); err != nil {
		t.Fatalf("create second subscription: %v", err)
	}

	listed, err := service.List(context.Background(), userID, SubscriptionQuery{Page: 1, PageSize: 1, Status: "paused"})
	if err != nil {
		t.Fatalf("list subscriptions: %v", err)
	}
	if listed.Total != 1 || len(listed.Items) != 1 || listed.Items[0].Name != "B" {
		t.Fatalf("unexpected filtered list: %+v", listed)
	}
}

func TestSubscriptionServiceRejectsInvalidInput(t *testing.T) {
	service := NewSubscriptionService(newFakeSubscriptionRepository(), newFakeCategoryRepository())

	req := validSubscriptionRequest(nil)
	req.Amount = 0
	if _, err := service.Create(context.Background(), uuid.New(), req); !errors.Is(err, ErrInvalidSubscription) {
		t.Fatalf("expected invalid subscription error, got %v", err)
	}

	if _, err := service.List(context.Background(), uuid.New(), SubscriptionQuery{Page: 1, PageSize: 101}); !errors.Is(err, ErrInvalidSubscriptionQry) {
		t.Fatalf("expected invalid subscription query error, got %v", err)
	}
}

func validSubscriptionRequest(categoryID *uuid.UUID) dto.SubscriptionRequest {
	return dto.SubscriptionRequest{
		CategoryID:      categoryID,
		Name:            "Netflix",
		Amount:          15.99,
		Currency:        "USD",
		BillingCycle:    "monthly",
		NextPaymentDate: "2026-06-01",
		Status:          "active",
	}
}

func createTestCategory(t *testing.T, categories *fakeCategoryRepository, userID uuid.UUID, name string) *model.Category {
	t.Helper()

	category := &model.Category{UserID: userID, Name: name, Color: defaultCategoryColor}
	if err := categories.Create(context.Background(), category); err != nil {
		t.Fatalf("create category fixture: %v", err)
	}

	return category
}

type fakeSubscriptionRepository struct {
	byID map[uuid.UUID]*model.Subscription
}

func newFakeSubscriptionRepository() *fakeSubscriptionRepository {
	return &fakeSubscriptionRepository{byID: map[uuid.UUID]*model.Subscription{}}
}

func (r *fakeSubscriptionRepository) Create(ctx context.Context, subscription *model.Subscription) error {
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

func (r *fakeSubscriptionRepository) FindByIDForUser(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Subscription, error) {
	subscription, ok := r.byID[id]
	if !ok || subscription.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	copy := *subscription
	return &copy, nil
}

func (r *fakeSubscriptionRepository) ListByUser(ctx context.Context, filter repository.SubscriptionFilter) ([]model.Subscription, int64, error) {
	matches := []model.Subscription{}
	for _, subscription := range r.byID {
		if subscription.UserID != filter.UserID {
			continue
		}
		if filter.Status != "" && subscription.Status != filter.Status {
			continue
		}
		if filter.CategoryID != nil && (subscription.CategoryID == nil || *subscription.CategoryID != *filter.CategoryID) {
			continue
		}
		if filter.BillingCycle != "" && subscription.BillingCycle != filter.BillingCycle {
			continue
		}
		matches = append(matches, *subscription)
	}
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].NextPaymentDate.Before(matches[j].NextPaymentDate)
	})

	total := int64(len(matches))
	if filter.Offset >= len(matches) {
		return []model.Subscription{}, total, nil
	}
	end := filter.Offset + filter.Limit
	if end > len(matches) {
		end = len(matches)
	}
	return matches[filter.Offset:end], total, nil
}

func (r *fakeSubscriptionRepository) Update(ctx context.Context, subscription *model.Subscription) error {
	copy := *subscription
	copy.UpdatedAt = time.Now().UTC()
	r.byID[subscription.ID] = &copy
	return nil
}

func (r *fakeSubscriptionRepository) Delete(ctx context.Context, subscription *model.Subscription) error {
	delete(r.byID, subscription.ID)
	return nil
}
