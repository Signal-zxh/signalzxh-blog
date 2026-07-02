package service

import (
	"testing"

	"github.com/Signal-zxh/signalzxh-blog/db"
	"github.com/Signal-zxh/signalzxh-blog/model"
)

type fakeTagRepo struct {
	tags      []model.Tag
	getErr    error
	createErr error
	updateErr error
	deleteErr error
}

func (f *fakeTagRepo) GetTags() ([]model.Tag, error) {
	return f.tags, f.getErr
}

func (f *fakeTagRepo) GetTagByID(id int) (model.Tag, error) {
	for _, t := range f.tags {
		if t.ID == id {
			return t, nil
		}
	}
	return model.Tag{}, db.ErrNotFound
}

func (f *fakeTagRepo) CreateTag(name string) (int64, error) {
	return 1, f.createErr
}

func (f *fakeTagRepo) UpdateTag(id int, name string) error {
	return f.updateErr
}

func (f *fakeTagRepo) DeleteTag(id int) error {
	return f.deleteErr
}

func (f *fakeTagRepo) GetOrCreateTag(name string) (int64, error) {
	return 1, nil
}

func (f *fakeTagRepo) GetTagsByPostID(postID int) ([]model.Tag, error) {
	return nil, nil
}

func (f *fakeTagRepo) AddTagsToPost(postID int, tagIDs []int) error {
	return nil
}

func (f *fakeTagRepo) RemoveTagsFromPost(postID int) error {
	return nil
}

func TestTagService_CreateTag_RejectsInvalidName(t *testing.T) {
	svc := NewTagService(&fakeTagRepo{})
	_, err := svc.CreateTag("")
	if err != ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestTagService_UpdateTag_MapsNotFound(t *testing.T) {
	repo := &fakeTagRepo{updateErr: db.ErrNoRowsAffected}
	svc := NewTagService(repo)
	err := svc.UpdateTag(1, "new")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestTagService_DeleteTag_MapsNotFound(t *testing.T) {
	repo := &fakeTagRepo{deleteErr: db.ErrNoRowsAffected}
	svc := NewTagService(repo)
	err := svc.DeleteTag(1)
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
