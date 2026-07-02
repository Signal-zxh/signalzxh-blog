package service

import (
	"testing"

	"github.com/Signal-zxh/signalzxh-blog/db"
	"github.com/Signal-zxh/signalzxh-blog/model"
)

type fakeCategoryRepo struct {
	categories []model.Category
	getErr     error
	createErr  error
	updateErr  error
	deleteErr  error
}

func (f *fakeCategoryRepo) GetCategories() ([]model.Category, error) {
	return f.categories, f.getErr
}

func (f *fakeCategoryRepo) GetCategoryByID(id int) (model.Category, error) {
	for _, c := range f.categories {
		if c.ID == id {
			return c, nil
		}
	}
	return model.Category{}, db.ErrNotFound
}

func (f *fakeCategoryRepo) CreateCategory(name string) (int64, error) {
	return 1, f.createErr
}

func (f *fakeCategoryRepo) UpdateCategory(id int, name string) error {
	return f.updateErr
}

func (f *fakeCategoryRepo) DeleteCategory(id int) error {
	return f.deleteErr
}

func TestCategoryService_CreateCategory_RejectsInvalidName(t *testing.T) {
	svc := NewCategoryService(&fakeCategoryRepo{})
	_, err := svc.CreateCategory("")
	if err != ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestCategoryService_UpdateCategory_MapsNotFound(t *testing.T) {
	repo := &fakeCategoryRepo{updateErr: db.ErrNoRowsAffected}
	svc := NewCategoryService(repo)
	err := svc.UpdateCategory(1, "new")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCategoryService_DeleteCategory_MapsNotFound(t *testing.T) {
	repo := &fakeCategoryRepo{deleteErr: db.ErrNoRowsAffected}
	svc := NewCategoryService(repo)
	err := svc.DeleteCategory(1)
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
