package service

import (
	"errors"

	"github.com/Signal-zxh/signalzxh-blog/db"
	"github.com/Signal-zxh/signalzxh-blog/model"
)

type CategoryService interface {
	GetCategories() ([]model.Category, error)
	GetCategoryByID(id int) (model.Category, error)
	CreateCategory(name string) (int64, error)
	UpdateCategory(id int, name string) error
	DeleteCategory(id int) error
}

type categoryService struct {
	repo db.CategoryRepo
}

func NewCategoryService(repo db.CategoryRepo) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetCategories() ([]model.Category, error) {
	return s.repo.GetCategories()
}

func (s *categoryService) GetCategoryByID(id int) (model.Category, error) {
	if id <= 0 {
		return model.Category{}, ErrInvalidInput
	}
	cat, err := s.repo.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return model.Category{}, ErrNotFound
		}
		return model.Category{}, err
	}
	return cat, nil
}

func (s *categoryService) CreateCategory(name string) (int64, error) {
	if name == "" || len(name) > 100 {
		return 0, ErrInvalidInput
	}
	return s.repo.CreateCategory(name)
}

func (s *categoryService) UpdateCategory(id int, name string) error {
	if id <= 0 || name == "" || len(name) > 100 {
		return ErrInvalidInput
	}
	err := s.repo.UpdateCategory(id, name)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsAffected) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *categoryService) DeleteCategory(id int) error {
	if id <= 0 {
		return ErrInvalidInput
	}
	err := s.repo.DeleteCategory(id)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsAffected) {
			return ErrNotFound
		}
		return err
	}
	return nil
}