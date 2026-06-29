package router

import (
	"github.com/Signal-zxh/signal-zxh/handler"
	"github.com/Signal-zxh/signal-zxh/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(postHandler *handler.PostHandler) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Logger())

	RegisterPage(r)
	RegisterAuth(r, postHandler)
	RegisterAPI(r, postHandler)

	return r
}
