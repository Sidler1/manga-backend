package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		log.Printf("| %3d | %13v | %15s | %s |",
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			path,
		)
	}
}
