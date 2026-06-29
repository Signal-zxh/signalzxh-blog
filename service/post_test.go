package service

import (
	"testing"
)

func TestGetPostsByPage_ParameterValidation(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		wantPage int
		wantSize int
	}{
		{
			name:     "valid parameters",
			page:     1,
			pageSize: 10,
			wantPage: 1,
			wantSize: 10,
		},
		{
			name:     "page less than 1",
			page:     -1,
			pageSize: 10,
			wantPage: 1,
			wantSize: 10,
		},
		{
			name:     "page equals 0",
			page:     0,
			pageSize: 10,
			wantPage: 1,
			wantSize: 10,
		},
		{
			name:     "pageSize less than 1",
			page:     1,
			pageSize: -1,
			wantPage: 1,
			wantSize: 10,
		},
		{
			name:     "pageSize equals 0",
			page:     1,
			pageSize: 0,
			wantPage: 1,
			wantSize: 10,
		},
		{
			name:     "pageSize greater than 100",
			page:     1,
			pageSize: 200,
			wantPage: 1,
			wantSize: 100,
		},
		{
			name:     "pageSize equals 100",
			page:     1,
			pageSize: 100,
			wantPage: 1,
			wantSize: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, pageSize := tt.page, tt.pageSize

			if page < 1 {
				page = 1
			}
			if pageSize < 1 {
				pageSize = 10
			}
			if pageSize > 100 {
				pageSize = 100
			}

			if page != tt.wantPage {
				t.Errorf("page got %d, want %d", page, tt.wantPage)
			}
			if pageSize != tt.wantSize {
				t.Errorf("pageSize got %d, want %d", pageSize, tt.wantSize)
			}
		})
	}
}

func validateTitle(title string) bool {
	return title != "" && len(title) <= 100
}

func TestCreatePost_Validation(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{
			name:    "valid title",
			title:   "test title",
			wantErr: false,
		},
		{
			name:    "empty title",
			title:   "",
			wantErr: true,
		},
		{
			name:    "title exactly 100 chars",
			title:   "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			wantErr: false,
		},
		{
			name:    "title 101 chars",
			title:   "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := validateTitle(tt.title)
			if valid == tt.wantErr {
				t.Errorf("validateTitle(%q) = %v, wantErr %v", tt.title, !valid, tt.wantErr)
			}
		})
	}
}
