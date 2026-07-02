package router

import (
	"net/http"

	"github.com/Signal-zxh/signalzxh-blog/agent"
	"github.com/Signal-zxh/signalzxh-blog/handler"
	"github.com/Signal-zxh/signalzxh-blog/middleware"
	"github.com/Signal-zxh/signalzxh-blog/model"
	"github.com/Signal-zxh/signalzxh-blog/service"
	"github.com/gin-gonic/gin"
)

func RegisterAPI(r *gin.Engine, h *handler.PostHandler, ch *handler.CategoryHandler, th *handler.TagHandler, postService service.PostService) {
	toolService := agent.NewToolService(postService)
	t := handler.NewToolHandler(toolService)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.Success(gin.H{
			"message": "pong",
		}))
	})

	r.POST("/login", h.Login)
	r.GET("/posts", h.GetPosts)
	r.GET("/posts/:id", h.GetPostByID)
	r.GET("/posts/detail", h.GetPostsWithCategoryTag)
	r.GET("/posts/:id/detail", h.GetPostWithCategoryTag)

	r.GET("/categories", ch.GetCategories)
	r.GET("/categories/:id", ch.GetCategory)
	r.GET("/categories/:id/posts", h.GetPostsByCategory)
	r.GET("/tags/:id/posts", h.GetPostsByTag)

	api := r.Group("/api/tools")
	api.POST("/http", t.HttpProbe)
	api.POST("/agent", t.Agent)

	auth := r.Group("/api")
	auth.Use(middleware.Auth())

	auth.POST("/posts", h.CreatePost)
	auth.PUT("/posts/:id", h.UpdatePost)
	auth.DELETE("/posts/:id", h.DeletePost)
	auth.POST("/posts/with-tag", h.CreatePostWithCategoryTag)
	auth.PUT("/posts/:id/with-tag", h.UpdatePostWithCategoryTag)

	auth.POST("/categories", ch.CreateCategory)
	auth.PUT("/categories/:id", ch.UpdateCategory)
	auth.DELETE("/categories/:id", ch.DeleteCategory)

	auth.POST("/tags", th.CreateTag)
	auth.PUT("/tags/:id", th.UpdateTag)
	auth.DELETE("/tags/:id", th.DeleteTag)
	r.GET("/tags", th.GetTags)
	r.GET("/tags/:id", th.GetTag)
}
