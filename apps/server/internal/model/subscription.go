package model

import (
	"time"

	"github.com/google/uuid"
)

type BillingCycle string

const (
	BillingCycleMonthly BillingCycle = "monthly"
	BillingCycleYearly  BillingCycle = "yearly"
	BillingCycleCustom  BillingCycle = "custom"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusPaused   SubscriptionStatus = "paused"
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
)

type Subscription struct {
	Base

	UserID          uuid.UUID          `gorm:"type:uuid;not null;index;index:idx_subscriptions_user_status;index:idx_subscriptions_user_next_payment" json:"user_id"`
	CategoryID      *uuid.UUID         `gorm:"type:uuid;index" json:"category_id,omitempty"`
	Name            string             `gorm:"type:varchar(160);not null" json:"name"`
	Amount          float64            `gorm:"type:numeric(12,2);not null" json:"amount"`
	Currency        string             `gorm:"type:varchar(3);not null" json:"currency"`
	BillingCycle    BillingCycle       `gorm:"type:varchar(32);not null" json:"billing_cycle"`
	NextPaymentDate time.Time          `gorm:"not null;index:idx_subscriptions_user_next_payment" json:"next_payment_date"`
	Status          SubscriptionStatus `gorm:"type:varchar(32);not null;index:idx_subscriptions_user_status" json:"status"`
	Notes           string             `gorm:"type:text" json:"notes"`

	User     User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Category *Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category,omitempty"`
}
