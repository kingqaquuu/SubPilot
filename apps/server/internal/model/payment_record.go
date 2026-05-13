package model

import (
	"time"

	"github.com/google/uuid"
)

type PaymentRecord struct {
	Base

	UserID         uuid.UUID `gorm:"type:uuid;not null;index;index:idx_payment_records_user_paid_at" json:"user_id"`
	SubscriptionID uuid.UUID `gorm:"type:uuid;not null;index" json:"subscription_id"`
	Amount         float64   `gorm:"type:numeric(12,2);not null" json:"amount"`
	Currency       string    `gorm:"type:varchar(3);not null" json:"currency"`
	PaidAt         time.Time `gorm:"not null;index:idx_payment_records_user_paid_at" json:"paid_at"`
	Note           string    `gorm:"type:text" json:"note"`

	User         User         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Subscription Subscription `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
