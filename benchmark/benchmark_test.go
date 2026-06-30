package benchmark_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/Signal-zxh/signalzxh-blog/middleware"
	"github.com/Signal-zxh/signalzxh-blog/model"
	"github.com/Signal-zxh/signalzxh-blog/utils"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "benchmark-secret-key")
}

// JWT Token 生成性能测试
func BenchmarkGenerateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = utils.GenerateToken(i % 1000)
	}
}

// JWT Token 解析性能测试
func BenchmarkParseToken(b *testing.B) {
	token, _ := utils.GenerateToken(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = utils.ParseToken(token)
	}
}

// JWT Token 生成+解析性能测试
func BenchmarkGenerateAndParseToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		token, _ := utils.GenerateToken(i % 1000)
		_, _ = utils.ParseToken(token)
	}
}

// Post JSON 序列化性能测试
func BenchmarkPostMarshal(b *testing.B) {
	post := model.Post{
		ID:      1,
		Title:   "Benchmark Test Post Title",
		Content: "This is a benchmark test post content with some longer text to simulate real content.",
		UserID:  1,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(post)
	}
}

// Post JSON 反序列化性能测试
func BenchmarkPostUnmarshal(b *testing.B) {
	post := model.Post{
		ID:      1,
		Title:   "Benchmark Test Post Title",
		Content: "This is a benchmark test post content with some longer text to simulate real content.",
		UserID:  1,
	}
	data, _ := json.Marshal(post)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p model.Post
		_ = json.Unmarshal(data, &p)
	}
}

// Posts 列表序列化性能测试
func BenchmarkPostsMarshal(b *testing.B) {
	posts := make([]model.Post, 100)
	for i := 0; i < 100; i++ {
		posts[i] = model.Post{
			ID:      i + 1,
			Title:   "Benchmark Test Post Title",
			Content: "This is a benchmark test post content with some longer text to simulate real content.",
			UserID:  1,
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(posts)
	}
}

// Posts 列表反序列化性能测试
func BenchmarkPostsUnmarshal(b *testing.B) {
	posts := make([]model.Post, 100)
	for i := 0; i < 100; i++ {
		posts[i] = model.Post{
			ID:      i + 1,
			Title:   "Benchmark Test Post Title",
			Content: "This is a benchmark test post content with some longer text to simulate real content.",
			UserID:  1,
		}
	}
	data, _ := json.Marshal(posts)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p []model.Post
		_ = json.Unmarshal(data, &p)
	}
}

// Response 序列化性能测试
func BenchmarkResponseMarshal(b *testing.B) {
	resp := model.Success(model.Post{
		ID:      1,
		Title:   "Test",
		Content: "Content",
		UserID:  1,
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(resp)
	}
}

// Auth 中间件性能测试（有效 token）
func BenchmarkAuthMiddleware_ValidToken(b *testing.B) {
	r := gin.New()
	r.GET("/protected", middleware.Auth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	token, _ := utils.GenerateToken(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

// Auth 中间件性能测试（无 token）
func BenchmarkAuthMiddleware_NoToken(b *testing.B) {
	r := gin.New()
	r.GET("/protected", middleware.Auth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

// Logger 中间件性能测试
func BenchmarkLoggerMiddleware(b *testing.B) {
	r := gin.New()
	r.GET("/test", middleware.Logger(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}

// 并发 Token 生成测试
func BenchmarkGenerateToken_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = utils.GenerateToken(i % 1000)
			i++
		}
	})
}

// 并发 Token 解析测试
func BenchmarkParseToken_Parallel(b *testing.B) {
	token, _ := utils.GenerateToken(1)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = utils.ParseToken(token)
		}
	})
}

// 并发 JSON 序列化测试
func BenchmarkPostMarshal_Parallel(b *testing.B) {
	post := model.Post{
		ID:      1,
		Title:   "Benchmark Test Post Title",
		Content: "This is a benchmark test post content with some longer text to simulate real content.",
		UserID:  1,
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = json.Marshal(post)
		}
	})
}

// 并发 HTTP 请求测试
func BenchmarkHTTPHandler_Parallel(b *testing.B) {
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.Success("pong"))
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/ping", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
		}
	})
}

// 并发 Auth 中间件测试
func BenchmarkAuthMiddleware_Parallel(b *testing.B) {
	r := gin.New()
	r.GET("/protected", middleware.Auth(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	token, _ := utils.GenerateToken(1)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
		}
	})
}

// sync.Mutex 并发测试
func BenchmarkMutex(b *testing.B) {
	var mu sync.Mutex
	var count int
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			count++
			mu.Unlock()
		}
	})
}

// sync.RWMutex 读并发测试
func BenchmarkRWMutex_Read(b *testing.B) {
	var mu sync.RWMutex
	var count int
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.RLock()
			_ = count
			mu.RUnlock()
		}
	})
}

// sync.RWMutex 写并发测试
func BenchmarkRWMutex_Write(b *testing.B) {
	var mu sync.RWMutex
	var count int
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			count++
			mu.Unlock()
		}
	})
}

// sync.Map 并发测试
func BenchmarkSyncMap(b *testing.B) {
	var m sync.Map
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Store(i%100, i)
			_, _ = m.Load(i % 100)
			i++
		}
	})
}