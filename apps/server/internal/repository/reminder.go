package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/model"
	"gorm.io/gorm"
)

type gormReminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) ReminderRepository {
	return &gormReminderRepository{db: db}
}

func (r *gormReminderRepository) Create(ctx context.Context, reminder *model.Reminder) error {
	return r.db.WithContext(ctx).Create(reminder).Error
}

func (r *gormReminderRepository) FindBySubscriptionForUser(ctx context.Context, subscriptionID uuid.UUID, userID uuid.UUID) (*model.Reminder, error) {
	var reminder model.Reminder
	if err := r.db.WithContext(ctx).First(&reminder, "subscription_id = ? AND user_id = ?", subscriptionID, userID).Error; err != nil {
		return nil, err
	}

	return &reminder, nil
}
