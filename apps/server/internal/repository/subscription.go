package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/model"
	"gorm.io/gorm"
)

type gormSubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &gormSubscriptionRepository{db: db}
}

func (r *gormSubscriptionRepository) Create(ctx context.Context, subscription *model.Subscription) error {
	return r.db.WithContext(ctx).Create(subscription).Error
}

func (r *gormSubscriptionRepository) FindByIDForUser(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Subscription, error) {
	var subscription model.Subscription
	if err := r.db.WithContext(ctx).Preload("Category").First(&subscription, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (r *gormSubscriptionRepository) ListByUser(ctx context.Context, filter SubscriptionFilter) ([]model.Subscription, int64, error) {
	var subscriptions []model.Subscription
	query := r.db.WithContext(ctx).Model(&model.Subscription{}).Where("user_id = ?", filter.UserID)
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.BillingCycle != "" {
		query = query.Where("billing_cycle = ?", filter.BillingCycle)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Category").
		Order("next_payment_date ASC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&subscriptions).Error
	return subscriptions, total, err
}

func (r *gormSubscriptionRepository) Update(ctx context.Context, subscription *model.Subscription) error {
	return r.db.WithContext(ctx).Save(subscription).Error
}

func (r *gormSubscriptionRepository) Delete(ctx context.Context, subscription *model.Subscription) error {
	return r.db.WithContext(ctx).Delete(subscription).Error
}
