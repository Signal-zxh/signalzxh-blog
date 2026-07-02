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

type mockTagService struct{}

func (m *mockTagService) GetTags() ([]model.Tag, error) {
	return []model.Tag{{ID: 1, Name: "Go"}}, nil
}

func (m *mockTagService) GetTagByID(id int) (model.Tag, error) {
	if id == 999 {
		return model.Tag{}, service.ErrNotFound
	}
	return model.Tag{ID: id, Name: "Go"}, nil
}

func (m *mockTagService) CreateTag(name string) (int64, error) {
	if name == "" {
		return 0, service.ErrInvalidInput
	}
	return 42, nil
}

func (m *mockTagService) UpdateTag(id int, name string) error {
	if id == 999 {
		return service.ErrNotFound
	}
	if name == "" {
		return service.ErrInvalidInput
	}
	return nil
}

func (m *mockTagService) DeleteTag(id int) error {
	if id == 999 {
		return service.ErrNotFound
	}
	if id <= 0 {
		return service.ErrInvalidInput
	}
	return nil
}

func TestCreateTagHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	tagHandler := handler.NewTagHandler(&mockTagService{})
	r.POST("/tags", tagHandler.CreateTag)

	body := `{"name":"Go"}`
	req := httptest.NewRequest(http.MethodPost, "/tags", strings.NewReader(body))
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

func TestUpdateTagHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	tagHandler := handler.NewTagHandler(&mockTagService{})
	r.PUT("/tags/:id", tagHandler.UpdateTag)

	body := `{"name":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/tags/999", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteTagHandler_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	tagHandler := handler.NewTagHandler(&mockTagService{})
	r.DELETE("/tags/:id", tagHandler.DeleteTag)

	req := httptest.NewRequest(http.MethodDelete, "/tags/0", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
