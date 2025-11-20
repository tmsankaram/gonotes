package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		reqID, _ := c.Get(RequestIDkey)

		log.Printf(
			"[REQ %v] %s %s -> %d {%s}",
			reqID,
			c.Request.Method,
			c.Request.URL.Path,
			status,
			latency,
		)
	}
}
