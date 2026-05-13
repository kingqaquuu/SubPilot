package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/dto"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/model"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/repository"
	"gorm.io/gorm"
)

const defaultCategoryColor = "#64748b"

var (
	ErrCategoryNameExists = errors.New("category name already exists")
	ErrCategoryNotFound   = errors.New("category not found")
	ErrInvalidCategory    = errors.New("invalid category")
)

type CategoryService struct {
	categories repository.CategoryRepository
}

func NewCategoryService(categories repository.CategoryRepository) *CategoryService {
	return &CategoryService{categories: categories}
}

func (s *CategoryService) Create(ctx context.Context, userID uuid.UUID, req dto.CategoryRequest) (*dto.CategoryResponse, error) {
	name, color, err := normalizeCategoryRequest(req)
	if err != nil {
		return nil, err
	}

	exists, err := s.categoryNameExists(ctx, name, userID, uuid.Nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCategoryNameExists
	}

	category := &model.Category{
		UserID: userID,
		Name:   name,
		Color:  color,
	}
	if err := s.categories.Create(ctx, category); err != nil {
		if errors.Is(err, repository.ErrDuplicateCategoryName) {
			return nil, ErrCategoryNameExists
		}
		return nil, fmt.Errorf("create category: %w", err)
	}

	resp := categoryResponse(category)
	return &resp, nil
}

func (s *CategoryService) List(ctx context.Context, userID uuid.UUID) ([]dto.CategoryResponse, error) {
	categories, err := s.categories.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}

	resp := make([]dto.CategoryResponse, 0, len(categories))
	for i := range categories {
		resp = append(resp, categoryResponse(&categories[i]))
	}

	return resp, nil
}

func (s *CategoryService) Update(ctx context.Context, userID uuid.UUID, id uuid.UUID, req dto.CategoryRequest) (*dto.CategoryResponse, error) {
	name, color, err := normalizeCategoryRequest(req)
	if err != nil {
		return nil, err
	}

	category, err := s.categories.FindByIDForUser(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("find category: %w", err)
	}

	exists, err := s.categoryNameExists(ctx, name, userID, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrCategoryNameExists
	}

	category.Name = name
	category.Color = color
	if err := s.categories.Update(ctx, category); err != nil {
		if errors.Is(err, repository.ErrDuplicateCategoryName) {
			return nil, ErrCategoryNameExists
		}
		return nil, fmt.Errorf("update category: %w", err)
	}

	resp := categoryResponse(category)
	return &resp, nil
}

func (s *CategoryService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*dto.CategoryResponse, error) {
	category, err := s.categories.FindByIDForUser(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("find category: %w", err)
	}

	if err := s.categories.Delete(ctx, category); err != nil {
		return nil, fmt.Errorf("delete category: %w", err)
	}

	resp := categoryResponse(category)
	return &resp, nil
}

func (s *CategoryService) categoryNameExists(ctx context.Context, name string, userID uuid.UUID, exceptID uuid.UUID) (bool, error) {
	category, err := s.categories.FindByNameForUser(ctx, name, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("find category by name: %w", err)
	}

	return exceptID == uuid.Nil || category.ID != exceptID, nil
}

func normalizeCategoryRequest(req dto.CategoryRequest) (string, string, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > 120 {
		return "", "", ErrInvalidCategory
	}

	color := strings.TrimSpace(req.Color)
	if color == "" {
		color = defaultCategoryColor
	}
	if len(color) > 32 {
		return "", "", ErrInvalidCategory
	}

	return name, color, nil
}

func categoryResponse(category *model.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Color:     category.Color,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}
