package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDkey = "reqID"

func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := uuid.New().String()
		ctx.Set(RequestIDkey, id)
		ctx.Writer.Header().Set("X-Request-ID", id)
		ctx.Next()
	}
}
