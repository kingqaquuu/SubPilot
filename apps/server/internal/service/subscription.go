package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/dto"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/model"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/repository"
	"gorm.io/gorm"
)

const (
	defaultSubscriptionPage     = 1
	defaultSubscriptionPageSize = 20
	maxSubscriptionPageSize     = 100
)

var (
	ErrInvalidSubscription    = errors.New("invalid subscription")
	ErrSubscriptionNotFound   = errors.New("subscription not found")
	ErrSubscriptionCategory   = errors.New("subscription category not found")
	ErrInvalidSubscriptionQry = errors.New("invalid subscription query")
)

type SubscriptionService struct {
	subscriptions repository.SubscriptionRepository
	categories    repository.CategoryRepository
}

type SubscriptionQuery struct {
	Page         int
	PageSize     int
	Status       string
	CategoryID   *uuid.UUID
	BillingCycle string
}

func NewSubscriptionService(subscriptions repository.SubscriptionRepository, categories repository.CategoryRepository) *SubscriptionService {
	return &SubscriptionService{subscriptions: subscriptions, categories: categories}
}

func (s *SubscriptionService) Create(ctx context.Context, userID uuid.UUID, req dto.SubscriptionRequest) (*dto.SubscriptionResponse, error) {
	normalized, err := s.normalizeRequest(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	subscription := &model.Subscription{
		UserID:          userID,
		CategoryID:      normalized.CategoryID,
		Name:            normalized.Name,
		Amount:          normalized.Amount,
		Currency:        normalized.Currency,
		BillingCycle:    normalized.BillingCycle,
		NextPaymentDate: normalized.NextPaymentDate,
		Status:          normalized.Status,
		Notes:           normalized.Notes,
	}
	if err := s.subscriptions.Create(ctx, subscription); err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}

	return s.Detail(ctx, userID, subscription.ID)
}

func (s *SubscriptionService) List(ctx context.Context, userID uuid.UUID, query SubscriptionQuery) (*dto.SubscriptionListResponse, error) {
	filter, page, pageSize, err := normalizeSubscriptionQuery(userID, query)
	if err != nil {
		return nil, err
	}

	subscriptions, total, err := s.subscriptions.ListByUser(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}

	items := make([]dto.SubscriptionResponse, 0, len(subscriptions))
	for i := range subscriptions {
		items = append(items, subscriptionResponse(&subscriptions[i]))
	}

	return &dto.SubscriptionListResponse{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func (s *SubscriptionService) Detail(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*dto.SubscriptionResponse, error) {
	subscription, err := s.subscriptions.FindByIDForUser(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("find subscription: %w", err)
	}

	resp := subscriptionResponse(subscription)
	return &resp, nil
}

func (s *SubscriptionService) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.SubscriptionRequest) (*dto.SubscriptionResponse, error) {
	subscription, err := s.subscriptions.FindByIDForUser(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("find subscription: %w", err)
	}

	normalized, err := s.normalizeRequest(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	subscription.CategoryID = normalized.CategoryID
	subscription.Name = normalized.Name
	subscription.Amount = normalized.Amount
	subscription.Currency = normalized.Currency
	subscription.BillingCycle = normalized.BillingCycle
	subscription.NextPaymentDate = normalized.NextPaymentDate
	subscription.Status = normalized.Status
	subscription.Notes = normalized.Notes
	if err := s.subscriptions.Update(ctx, subscription); err != nil {
		return nil, fmt.Errorf("update subscription: %w", err)
	}

	return s.Detail(ctx, userID, subscription.ID)
}

func (s *SubscriptionService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*dto.SubscriptionResponse, error) {
	subscription, err := s.subscriptions.FindByIDForUser(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, fmt.Errorf("find subscription: %w", err)
	}

	resp := subscriptionResponse(subscription)
	if err := s.subscriptions.Delete(ctx, subscription); err != nil {
		return nil, fmt.Errorf("delete subscription: %w", err)
	}

	return &resp, nil
}

func (s *SubscriptionService) normalizeRequest(ctx context.Context, userID uuid.UUID, req dto.SubscriptionRequest) (*normalizedSubscription, error) {
	normalized, err := normalizeSubscriptionRequest(req)
	if err != nil {
		return nil, err
	}

	if normalized.CategoryID != nil {
		if _, err := s.categories.FindByIDForUser(ctx, *normalized.CategoryID, userID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrSubscriptionCategory
			}
			return nil, fmt.Errorf("find category: %w", err)
		}
	}

	return normalized, nil
}

type normalizedSubscription struct {
	CategoryID      *uuid.UUID
	Name            string
	Amount          float64
	Currency        string
	BillingCycle    model.BillingCycle
	NextPaymentDate time.Time
	Status          model.SubscriptionStatus
	Notes           string
}

func normalizeSubscriptionRequest(req dto.SubscriptionRequest) (*normalizedSubscription, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 160 || req.Amount <= 0 {
		return nil, ErrInvalidSubscription
	}

	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if len(currency) != 3 || !isASCIIAlpha(currency) {
		return nil, ErrInvalidSubscription
	}

	billingCycle := model.BillingCycle(strings.TrimSpace(req.BillingCycle))
	if !validBillingCycle(billingCycle) {
		return nil, ErrInvalidSubscription
	}

	status := model.SubscriptionStatus(strings.TrimSpace(req.Status))
	if status == "" {
		status = model.SubscriptionStatusActive
	}
	if !validSubscriptionStatus(status) {
		return nil, ErrInvalidSubscription
	}

	nextPaymentDate, err := parseSubscriptionDate(req.NextPaymentDate)
	if err != nil {
		return nil, ErrInvalidSubscription
	}

	notes := strings.TrimSpace(req.Notes)
	if len(notes) > 2000 {
		return nil, ErrInvalidSubscription
	}

	return &normalizedSubscription{
		CategoryID:      req.CategoryID,
		Name:            name,
		Amount:          req.Amount,
		Currency:        currency,
		BillingCycle:    billingCycle,
		NextPaymentDate: nextPaymentDate,
		Status:          status,
		Notes:           notes,
	}, nil
}

func normalizeSubscriptionQuery(userID uuid.UUID, query SubscriptionQuery) (repository.SubscriptionFilter, int, int, error) {
	page := query.Page
	if page == 0 {
		page = defaultSubscriptionPage
	}
	pageSize := query.PageSize
	if pageSize == 0 {
		pageSize = defaultSubscriptionPageSize
	}
	if page < 1 || pageSize < 1 || pageSize > maxSubscriptionPageSize {
		return repository.SubscriptionFilter{}, 0, 0, ErrInvalidSubscriptionQry
	}

	status := model.SubscriptionStatus(strings.TrimSpace(query.Status))
	if status != "" && !validSubscriptionStatus(status) {
		return repository.SubscriptionFilter{}, 0, 0, ErrInvalidSubscriptionQry
	}

	billingCycle := model.BillingCycle(strings.TrimSpace(query.BillingCycle))
	if billingCycle != "" && !validBillingCycle(billingCycle) {
		return repository.SubscriptionFilter{}, 0, 0, ErrInvalidSubscriptionQry
	}

	return repository.SubscriptionFilter{
		UserID:       userID,
		Status:       status,
		CategoryID:   query.CategoryID,
		BillingCycle: billingCycle,
		Limit:        pageSize,
		Offset:       (page - 1) * pageSize,
	}, page, pageSize, nil
}

func subscriptionResponse(subscription *model.Subscription) dto.SubscriptionResponse {
	resp := dto.SubscriptionResponse{
		ID:              subscription.ID,
		CategoryID:      subscription.CategoryID,
		Name:            subscription.Name,
		Amount:          subscription.Amount,
		Currency:        subscription.Currency,
		BillingCycle:    string(subscription.BillingCycle),
		NextPaymentDate: subscription.NextPaymentDate,
		Status:          string(subscription.Status),
		Notes:           subscription.Notes,
		CreatedAt:       subscription.CreatedAt,
		UpdatedAt:       subscription.UpdatedAt,
	}
	if subscription.Category != nil {
		category := categoryResponse(subscription.Category)
		resp.Category = &category
	}

	return resp
}

func parseSubscriptionDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, ErrInvalidSubscription
	}
	if parsed, err := time.Parse("2006-01-02", value); err == nil {
		return parsed, nil
	}

	return time.Parse(time.RFC3339, value)
}

func validBillingCycle(value model.BillingCycle) bool {
	return value == model.BillingCycleMonthly || value == model.BillingCycleYearly || value == model.BillingCycleCustom
}

func validSubscriptionStatus(value model.SubscriptionStatus) bool {
	return value == model.SubscriptionStatusActive || value == model.SubscriptionStatusPaused || value == model.SubscriptionStatusCanceled
}

func isASCIIAlpha(value string) bool {
	for _, r := range value {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}
