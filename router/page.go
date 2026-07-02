package router

import "github.com/gin-gonic/gin"

func RegisterPage(r *gin.Engine) {
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	r.GET("/admin", func(c *gin.Context) {
		c.File("./static/admin.html")
	})
	r.GET("/admin/categories", func(c *gin.Context) {
		c.File("./static/admin-categories.html")
	})
	r.GET("/tools", func(c *gin.Context) {
		c.File("./static/tools.html")
	})
	r.GET("/games", func(c *gin.Context) {
		c.File("./static/games.html")
	})
	r.GET("/about", func(c *gin.Context) {
		c.File("./static/about.html")
	})
	r.GET("/post-detail.html", func(c *gin.Context) {
		c.File("./static/post-detail.html")
	})
}
