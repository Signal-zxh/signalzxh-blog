package router

import (
	"github.com/Signal-zxh/signal-zxh/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.Logger())

	RegisterPage(r)
	RegisterAuth(r)
	RegisterAPI(r)

	return r
}
