package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/model"
	"gorm.io/gorm"
)

var ErrDuplicateCategoryName = errors.New("duplicate category name")

type gormCategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &gormCategoryRepository{db: db}
}

func (r *gormCategoryRepository) Create(ctx context.Context, category *model.Category) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		if isUniqueViolation(err) {
			return ErrDuplicateCategoryName
		}
		return err
	}

	return nil
}

func (r *gormCategoryRepository) FindByIDForUser(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).First(&category, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *gormCategoryRepository) FindByNameForUser(ctx context.Context, name string, userID uuid.UUID) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).First(&category, "name = ? AND user_id = ?", name, userID).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *gormCategoryRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *gormCategoryRepository) Update(ctx context.Context, category *model.Category) error {
	if err := r.db.WithContext(ctx).Save(category).Error; err != nil {
		if isUniqueViolation(err) {
			return ErrDuplicateCategoryName
		}
		return err
	}

	return nil
}

func (r *gormCategoryRepository) Delete(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Delete(category).Error
}
