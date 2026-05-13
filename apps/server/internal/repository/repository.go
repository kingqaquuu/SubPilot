package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	FindByIDForUser(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Category, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Category, error)
}

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *model.Subscription) error
	FindByIDForUser(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Subscription, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]model.Subscription, error)
}

type ReminderRepository interface {
	Create(ctx context.Context, reminder *model.Reminder) error
	FindBySubscriptionForUser(ctx context.Context, subscriptionID uuid.UUID, userID uuid.UUID) (*model.Reminder, error)
}

type PaymentRecordRepository interface {
	Create(ctx context.Context, paymentRecord *model.PaymentRecord) error
	ListBySubscriptionForUser(ctx context.Context, subscriptionID uuid.UUID, userID uuid.UUID, limit int, offset int) ([]model.PaymentRecord, error)
}

type Repositories struct {
	Users          UserRepository
	Categories     CategoryRepository
	Subscriptions  SubscriptionRepository
	Reminders      ReminderRepository
	PaymentRecords PaymentRecordRepository
}

func New(db *gorm.DB) Repositories {
	return Repositories{
		Users:          NewUserRepository(db),
		Categories:     NewCategoryRepository(db),
		Subscriptions:  NewSubscriptionRepository(db),
		Reminders:      NewReminderRepository(db),
		PaymentRecords: NewPaymentRecordRepository(db),
	}
}
