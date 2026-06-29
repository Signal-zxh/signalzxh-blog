package service

import (
	"testing"
	"time"

	"github.com/Signal-zxh/signal-zxh/model"
)

type mockPostRepo struct{}

func (m *mockPostRepo) GetPosts() ([]model.Post, error)                         { return nil, nil }
func (m *mockPostRepo) GetPostsByPage(page, pageSize int) ([]model.Post, error) { return nil, nil }
func (m *mockPostRepo) GetPostsCount() (int, error)                             { return 0, nil }
func (m *mockPostRepo) CreatePost(post model.Post) (int64, error)               { return 1, nil }
func (m *mockPostRepo) UpdatePost(post model.Post) error                        { return nil }
func (m *mockPostRepo) DeletePost(id int) error                                 { return nil }
func (m *mockPostRepo) GetPostByID(id int) (model.Post, error)                  { return model.Post{}, nil }

type mockPostCache struct{}

func (m *mockPostCache) GetPostByID(id int) (model.Post, bool, error) {
	return model.Post{}, false, nil
}
func (m *mockPostCache) SetPost(post model.Post, ttl time.Duration) error     { return nil }
func (m *mockPostCache) SetNilPost(id int, ttl time.Duration) error           { return nil }
func (m *mockPostCache) GetPosts() ([]model.Post, bool, error)                { return nil, false, nil }
func (m *mockPostCache) SetPosts(posts []model.Post, ttl time.Duration) error { return nil }
func (m *mockPostCache) GetPostsByPage(page, pageSize int) ([]model.Post, bool, error) {
	return nil, false, nil
}
func (m *mockPostCache) SetPostsByPage(posts []model.Post, page, pageSize int, ttl time.Duration) error {
	return nil
}
func (m *mockPostCache) InvalidatePost(id int) error { return nil }
func (m *mockPostCache) InvalidatePosts() error      { return nil }

func TestGetPostsByPage_ParameterValidation(t *testing.T) {
	service := NewPostService(&mockPostRepo{}, &mockPostCache{})

	tests := []struct {
		name     string
		page     int
		pageSize int
		wantPage int
		wantSize int
	}{
		{"valid_parameters", 1, 10, 1, 10},
		{"page_less_than_1", 0, 10, 1, 10},
		{"page_equals_0", -1, 10, 1, 10},
		{"pageSize_less_than_1", 1, 0, 1, 10},
		{"pageSize_equals_0", 1, -1, 1, 10},
		{"pageSize_greater_than_100", 1, 200, 1, 100},
		{"pageSize_equals_100", 1, 100, 1, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.GetPostsByPage(tt.page, tt.pageSize)
			if err != nil {
				t.Errorf("GetPostsByPage() error = %v", err)
			}
		})
	}
}

func TestCreatePost_Validation(t *testing.T) {
	service := NewPostService(&mockPostRepo{}, &mockPostCache{})

	tests := []struct {
		name    string
		title   string
		content string
		userID  int
		wantErr bool
	}{
		{"valid_title", "Test Title", "Test Content", 1, false},
		{"empty_title", "", "Test Content", 1, true},
		{"title_exactly_100_chars", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "content", 1, false},
		{"title_101_chars", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "content", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CreatePost(tt.title, tt.content, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
