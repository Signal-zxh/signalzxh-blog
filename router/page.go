package router

import "github.com/gin-gonic/gin"

func RegisterPage(r *gin.Engine) {
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

	// 文章详情页（公开）
	r.GET("/post-detail.html", func(c *gin.Context) {
		c.File("./static/post-detail.html")
	})
}
