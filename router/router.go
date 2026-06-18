package router

import (
	"github.com/Signal-zxh/signal-zxh/handler"
	"github.com/Signal-zxh/signal-zxh/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Logger())

	// 静态页面
	r.Static("/static", "./static")

	// 首页
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// handler
	h := &handler.PostHandler{}
	r.POST("/login", h.Login)

	auth := r.Group("/")
	auth.Use(middleware.Auth())

	auth.GET("/posts", h.GetPosts)
	auth.GET("/posts/:id", h.GetPostByID)
	auth.POST("/posts", h.CreatePost)
	auth.DELETE("/posts/:id", h.DeletePost)
	auth.PUT("/posts/:id", h.UpdatePost)

	return r
}
