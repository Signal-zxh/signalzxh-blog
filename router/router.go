package router

import (
	"github.com/Signal-zxh/signalzxh-blog/db"
	"github.com/Signal-zxh/signalzxh-blog/handler"
	"github.com/Signal-zxh/signalzxh-blog/middleware"
	_ "github.com/Signal-zxh/signalzxh-blog/docs"  // swagger docs
	"github.com/Signal-zxh/signalzxh-blog/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(postHandler *handler.PostHandler) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Logger())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	categoryService := service.NewCategoryService(db.CategoryRepoImpl)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	tagService := service.NewTagService(db.TagRepoImpl)
	tagHandler := handler.NewTagHandler(tagService)

	RegisterPage(r)
	RegisterAuth(r, postHandler)
	RegisterAPI(r, postHandler, categoryHandler, tagHandler)

	return r
}
