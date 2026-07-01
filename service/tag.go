package service

import (
	"errors"

	"github.com/Signal-zxh/signalzxh-blog/db"
	"github.com/Signal-zxh/signalzxh-blog/model"
)

type TagService interface {
	GetTags() ([]model.Tag, error)
	GetTagByID(id int) (model.Tag, error)
	CreateTag(name string) (int64, error)
	UpdateTag(id int, name string) error
	DeleteTag(id int) error
}

type tagService struct {
	repo db.TagRepo
}

func NewTagService(repo db.TagRepo) TagService {
	return &tagService{repo: repo}
}

func (s *tagService) GetTags() ([]model.Tag, error) {
	return s.repo.GetTags()
}

func (s *tagService) GetTagByID(id int) (model.Tag, error) {
	if id <= 0 {
		return model.Tag{}, ErrInvalidInput
	}
	tag, err := s.repo.GetTagByID(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return model.Tag{}, ErrNotFound
		}
		return model.Tag{}, err
	}
	return tag, nil
}

func (s *tagService) CreateTag(name string) (int64, error) {
	if name == "" || len(name) > 100 {
		return 0, ErrInvalidInput
	}
	return s.repo.CreateTag(name)
}

func (s *tagService) UpdateTag(id int, name string) error {
	if id <= 0 || name == "" || len(name) > 100 {
		return ErrInvalidInput
	}
	err := s.repo.UpdateTag(id, name)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsAffected) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *tagService) DeleteTag(id int) error {
	if id <= 0 {
		return ErrInvalidInput
	}
	err := s.repo.DeleteTag(id)
	if err != nil {
		if errors.Is(err, db.ErrNoRowsAffected) {
			return ErrNotFound
		}
		return err
	}
	return nil
}