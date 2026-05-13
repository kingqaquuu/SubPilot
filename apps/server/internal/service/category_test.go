package service

import (
	"context"
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/dto"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/model"
	"github.com/kingqaquuu/SubPilot/apps/server/internal/repository"
	"gorm.io/gorm"
)

func TestCategoryServiceCreateListUpdateAndDelete(t *testing.T) {
	categories := newFakeCategoryRepository()
	service := NewCategoryService(categories)
	userID := uuid.New()

	created, err := service.Create(context.Background(), userID, dto.CategoryRequest{
		Name:  " Streaming ",
		Color: "",
	})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	if created.Name != "Streaming" {
		t.Fatalf("expected trimmed name, got %q", created.Name)
	}
	if created.Color != defaultCategoryColor {
		t.Fatalf("expected default color, got %q", created.Color)
	}

	listed, err := service.List(context.Background(), userID)
	if err != nil {
		t.Fatalf("list categories: %v", err)
	}
	if len(listed) != 1 || listed[0].ID != created.ID {
		t.Fatalf("unexpected category list: %+v", listed)
	}

	updated, err := service.Update(context.Background(), userID, created.ID, dto.CategoryRequest{
		Name:  "Music",
		Color: "#a855f7",
	})
	if err != nil {
		t.Fatalf("update category: %v", err)
	}
	if updated.Name != "Music" || updated.Color != "#a855f7" {
		t.Fatalf("unexpected updated category: %+v", updated)
	}

	deleted, err := service.Delete(context.Background(), userID, created.ID)
	if err != nil {
		t.Fatalf("delete category: %v", err)
	}
	if deleted.ID != created.ID {
		t.Fatalf("unexpected deleted category: %+v", deleted)
	}
	if !categories.clearedCategoryRefs[categoryRefKey{userID: userID, categoryID: created.ID}] {
		t.Fatal("expected delete to clear same-user subscription category references")
	}

	listed, err = service.List(context.Background(), userID)
	if err != nil {
		t.Fatalf("list categories after delete: %v", err)
	}
	if len(listed) != 0 {
		t.Fatalf("expected empty category list, got %+v", listed)
	}
}

func TestCategoryServiceRejectsDuplicateNamesForSameUser(t *testing.T) {
	categories := newFakeCategoryRepository()
	service := NewCategoryService(categories)
	userID := uuid.New()

	if _, err := service.Create(context.Background(), userID, dto.CategoryRequest{Name: "Streaming"}); err != nil {
		t.Fatalf("create first category: %v", err)
	}
	if _, err := service.Create(context.Background(), userID, dto.CategoryRequest{Name: "Streaming"}); !errors.Is(err, ErrCategoryNameExists) {
		t.Fatalf("expected duplicate category error, got %v", err)
	}
}

func TestCategoryServiceAllowsSameNameForDifferentUsers(t *testing.T) {
	categories := newFakeCategoryRepository()
	service := NewCategoryService(categories)

	if _, err := service.Create(context.Background(), uuid.New(), dto.CategoryRequest{Name: "Streaming"}); err != nil {
		t.Fatalf("create first user category: %v", err)
	}
	if _, err := service.Create(context.Background(), uuid.New(), dto.CategoryRequest{Name: "Streaming"}); err != nil {
		t.Fatalf("expected same category name for another user to be allowed, got %v", err)
	}
}

func TestCategoryServiceRejectsCrossUserAccess(t *testing.T) {
	categories := newFakeCategoryRepository()
	service := NewCategoryService(categories)
	ownerID := uuid.New()
	otherUserID := uuid.New()

	created, err := service.Create(context.Background(), ownerID, dto.CategoryRequest{Name: "Streaming"})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}

	if _, err := service.Update(context.Background(), otherUserID, created.ID, dto.CategoryRequest{Name: "Music"}); !errors.Is(err, ErrCategoryNotFound) {
		t.Fatalf("expected not found for cross-user update, got %v", err)
	}
	if _, err := service.Delete(context.Background(), otherUserID, created.ID); !errors.Is(err, ErrCategoryNotFound) {
		t.Fatalf("expected not found for cross-user delete, got %v", err)
	}
	if categories.clearedCategoryRefs[categoryRefKey{userID: otherUserID, categoryID: created.ID}] {
		t.Fatal("did not expect cross-user delete to clear subscription category references")
	}
}

func TestCategoryServiceMapsRepositoryDuplicateName(t *testing.T) {
	categories := newFakeCategoryRepository()
	categories.createErr = repository.ErrDuplicateCategoryName
	service := NewCategoryService(categories)

	if _, err := service.Create(context.Background(), uuid.New(), dto.CategoryRequest{Name: "Streaming"}); !errors.Is(err, ErrCategoryNameExists) {
		t.Fatalf("expected duplicate category error, got %v", err)
	}
}

func TestCategoryServiceRejectsInvalidInput(t *testing.T) {
	service := NewCategoryService(newFakeCategoryRepository())

	if _, err := service.Create(context.Background(), uuid.New(), dto.CategoryRequest{Name: "   "}); !errors.Is(err, ErrInvalidCategory) {
		t.Fatalf("expected invalid category error, got %v", err)
	}
}

type fakeCategoryRepository struct {
	byID                map[uuid.UUID]*model.Category
	clearedCategoryRefs map[categoryRefKey]bool
	createErr           error
	updateErr           error
}

func newFakeCategoryRepository() *fakeCategoryRepository {
	return &fakeCategoryRepository{
		byID:                map[uuid.UUID]*model.Category{},
		clearedCategoryRefs: map[categoryRefKey]bool{},
	}
}

type categoryRefKey struct {
	userID     uuid.UUID
	categoryID uuid.UUID
}

func (r *fakeCategoryRepository) Create(ctx context.Context, category *model.Category) error {
	if r.createErr != nil {
		return r.createErr
	}
	for _, existing := range r.byID {
		if existing.UserID == category.UserID && existing.Name == category.Name {
			return repository.ErrDuplicateCategoryName
		}
	}
	if category.ID == uuid.Nil {
		category.ID = uuid.New()
	}
	now := time.Now().UTC()
	category.CreatedAt = now
	category.UpdatedAt = now
	copy := *category
	r.byID[category.ID] = &copy
	return nil
}

func (r *fakeCategoryRepository) FindByIDForUser(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*model.Category, error) {
	category, ok := r.byID[id]
	if !ok || category.UserID != userID {
		return nil, gorm.ErrRecordNotFound
	}
	copy := *category
	return &copy, nil
}

func (r *fakeCategoryRepository) FindByNameForUser(ctx context.Context, name string, userID uuid.UUID) (*model.Category, error) {
	for _, category := range r.byID {
		if category.UserID == userID && category.Name == name {
			copy := *category
			return &copy, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *fakeCategoryRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Category, error) {
	categories := []model.Category{}
	for _, category := range r.byID {
		if category.UserID == userID {
			categories = append(categories, *category)
		}
	}
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Name < categories[j].Name
	})
	return categories, nil
}

func (r *fakeCategoryRepository) Update(ctx context.Context, category *model.Category) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	for id, existing := range r.byID {
		if id != category.ID && existing.UserID == category.UserID && existing.Name == category.Name {
			return repository.ErrDuplicateCategoryName
		}
	}
	copy := *category
	copy.UpdatedAt = time.Now().UTC()
	r.byID[category.ID] = &copy
	return nil
}

func (r *fakeCategoryRepository) Delete(ctx context.Context, category *model.Category) error {
	r.clearedCategoryRefs[categoryRefKey{userID: category.UserID, categoryID: category.ID}] = true
	delete(r.byID, category.ID)
	return nil
}
