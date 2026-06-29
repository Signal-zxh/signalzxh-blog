package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Signal-zxh/signal-zxh/handler"
	"github.com/Signal-zxh/signal-zxh/router"
	"github.com/Signal-zxh/signal-zxh/service"
)

func TestLogin(t *testing.T) {
	os.Setenv("ADMIN_USERNAME", "admin")
	os.Setenv("ADMIN_PASSWORD", "123456")

	postHandler := handler.NewPostHandler(&service.PostService{})
	r := router.SetupRouter(postHandler)

	body := `{"username":"admin","password":"123"}`
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("login failed, got %d", w.Code)
	}

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("解析响应失败：%v", err)
	}

	tokenData, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatal("响应中没有data字段")
	}
	token, ok := tokenData["token"].(string)
	if !ok || token == "" {
		t.Fatal("未获取token")
	}
}
