package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		traceId := uuid.New().String()[:8]

		c.Set("traceId", traceId)

		c.Next() // 继续往下走

		cost := time.Since(start)

		log.Printf("[REQ][%s] %s %s %d %v",
			traceId,
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			cost,
		)
	}
}
