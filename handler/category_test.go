package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Signal-zxh/signalzxh-blog/handler"
	"github.com/Signal-zxh/signalzxh-blog/model"
	"github.com/Signal-zxh/signalzxh-blog/service"
	"github.com/gin-gonic/gin"
)

type mockCategoryService struct{}

func (m *mockCategoryService) GetCategories() ([]model.Category, error) {
	return []model.Category{{ID: 1, Name: "Go"}}, nil
}

func (m *mockCategoryService) GetCategoryByID(id int) (model.Category, error) {
	if id == 999 {
		return model.Category{}, service.ErrNotFound
	}
	return model.Category{ID: id, Name: "Go"}, nil
}

func (m *mockCategoryService) CreateCategory(name string) (int64, error) {
	if name == "" {
		return 0, service.ErrInvalidInput
	}
	return 42, nil
}

func (m *mockCategoryService) UpdateCategory(id int, name string) error {
	if id == 999 {
		return service.ErrNotFound
	}
	if name == "" {
		return service.ErrInvalidInput
	}
	return nil
}

func (m *mockCategoryService) DeleteCategory(id int) error {
	if id == 999 {
		return service.ErrNotFound
	}
	if id <= 0 {
		return service.ErrInvalidInput
	}
	return nil
}

func TestCreateCategoryHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	categoryHandler := handler.NewCategoryHandler(&mockCategoryService{})
	r.POST("/categories", categoryHandler.CreateCategory)

	body := `{"name":"Go"}`
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp model.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("expected success code 0, got %v", resp.Code)
	}
}

func TestUpdateCategoryHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	categoryHandler := handler.NewCategoryHandler(&mockCategoryService{})
	r.PUT("/categories/:id", categoryHandler.UpdateCategory)

	body := `{"name":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/categories/999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteCategoryHandler_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	categoryHandler := handler.NewCategoryHandler(&mockCategoryService{})
	r.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	req := httptest.NewRequest(http.MethodDelete, "/categories/0", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
