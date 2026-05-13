package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/model"
	"gorm.io/gorm"
)

type gormPaymentRecordRepository struct {
	db *gorm.DB
}

func NewPaymentRecordRepository(db *gorm.DB) PaymentRecordRepository {
	return &gormPaymentRecordRepository{db: db}
}

func (r *gormPaymentRecordRepository) Create(ctx context.Context, paymentRecord *model.PaymentRecord) error {
	return r.db.WithContext(ctx).Create(paymentRecord).Error
}

func (r *gormPaymentRecordRepository) ListBySubscriptionForUser(ctx context.Context, subscriptionID uuid.UUID, userID uuid.UUID, limit int, offset int) ([]model.PaymentRecord, error) {
	var paymentRecords []model.PaymentRecord
	err := r.db.WithContext(ctx).
		Where("subscription_id = ? AND user_id = ?", subscriptionID, userID).
		Order("paid_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&paymentRecords).Error
	return paymentRecords, err
}
