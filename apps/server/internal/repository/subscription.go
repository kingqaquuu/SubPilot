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
	if err := r.db.WithContext(ctx).First(&subscription, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (r *gormSubscriptionRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]model.Subscription, error) {
	var subscriptions []model.Subscription
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("next_payment_date ASC").
		Limit(limit).
		Offset(offset).
		Find(&subscriptions).Error
	return subscriptions, err
}
