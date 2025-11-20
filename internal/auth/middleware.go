package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tmsankram/gonotes/internal/response"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			response.Unauthorized(c, Err("missing or invalid token"))
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(h, "Bearer ")
		claims, err := ValidateToken(tokenStr)
		if err != nil {
			response.Unauthorized(c, Err("invalid token"))
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func Err(msg string) error {
	return &AuthError{msg}
}

type AuthError struct {
	msg string
}

func (e *AuthError) Error() string {
	return e.msg
}
