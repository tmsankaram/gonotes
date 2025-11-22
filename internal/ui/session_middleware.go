package ui

import (
	"github.com/gin-gonic/gin"

	"github.com/tmsankram/gonotes/internal/auth"
	"github.com/tmsankram/gonotes/internal/users"
)

func SessionMiddleware(usersSvc *users.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(JWT_COOKIE)
		if err != nil || token == "" {
			c.Next()
			return
		}
		claims, err := auth.ValidateToken(token)
		if err != nil {
			// If invalid token, clear cookie
			c.SetCookie(JWT_COOKIE, "", -1, "/", "", false, true)
			c.Next()
			return
		}
		u, err := usersSvc.GetByID(claims.UserID)
		if err != nil {
			c.Next()
			return
		}
		c.Set("user", u) // templates can use .User
		c.Next()
	}
}
