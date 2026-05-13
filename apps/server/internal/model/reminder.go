package model

import "github.com/google/uuid"

type Reminder struct {
	Base

	UserID           uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	SubscriptionID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_reminders_subscription" json:"subscription_id"`
	RemindBeforeDays int       `gorm:"not null;default:3" json:"remind_before_days"`
	Enabled          bool      `gorm:"not null;default:true" json:"enabled"`

	User         User         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Subscription Subscription `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
