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

	// 管理页面
	r.GET("/admin", func(c *gin.Context) {
		c.File("./static/admin.html")
	})

	// 工具页面（公开）
	r.GET("/tools", func(c *gin.Context) {
		c.File("./static/tools.html")
	})

	// 游戏页面（公开）
	r.GET("/games", func(c *gin.Context) {
		c.File("./static/games.html")
	})

	// 关于页面（公开）
	r.GET("/about", func(c *gin.Context) {
		c.File("./static/about.html")
	})

	// handler
	h := &handler.PostHandler{}
	r.POST("/login", h.Login)
	r.GET("/posts", h.GetPosts)
	r.GET("/posts/:id", h.GetPostByID)

	auth := r.Group("/")
	auth.Use(middleware.Auth())

	auth.POST("/posts", h.CreatePost)
	auth.DELETE("/posts/:id", h.DeletePost)
	auth.PUT("/posts/:id", h.UpdatePost)

	return r
}
