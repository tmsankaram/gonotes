package middleware

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/tmsankram/gonotes/internal/response"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] %v", r)
				response.Internal(c, errors.New("internal server error"))
				c.Abort()
			}
		}()
		c.Next()
	}
}
