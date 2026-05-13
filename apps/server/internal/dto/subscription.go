package dto

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionRequest struct {
	CategoryID      *uuid.UUID `json:"category_id"`
	Name            string     `json:"name" binding:"required,max=160"`
	Amount          float64    `json:"amount" binding:"required"`
	Currency        string     `json:"currency" binding:"required,len=3"`
	BillingCycle    string     `json:"billing_cycle" binding:"required"`
	NextPaymentDate string     `json:"next_payment_date" binding:"required"`
	Status          string     `json:"status"`
	Notes           string     `json:"notes" binding:"max=2000"`
}

type SubscriptionResponse struct {
	ID              uuid.UUID         `json:"id"`
	CategoryID      *uuid.UUID        `json:"category_id,omitempty"`
	Category        *CategoryResponse `json:"category,omitempty"`
	Name            string            `json:"name"`
	Amount          float64           `json:"amount"`
	Currency        string            `json:"currency"`
	BillingCycle    string            `json:"billing_cycle"`
	NextPaymentDate time.Time         `json:"next_payment_date"`
	Status          string            `json:"status"`
	Notes           string            `json:"notes"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type SubscriptionListResponse struct {
	Items    []SubscriptionResponse `json:"items"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
	Total    int64                  `json:"total"`
}
