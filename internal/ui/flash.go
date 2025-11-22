package ui

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

func FlashMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		flash, err := ctx.Cookie("flash")
		if err == nil && flash != "" {
			decoded, _ := url.QueryUnescape(flash)
			ctx.Set("flash", decoded)

			// delete it
			ctx.SetCookie("flash", "", -1, "/", "", false, true)
		}
		ctx.Next()
	}
}

func Flash(ctx *gin.Context, msg string) {
	ctx.SetCookie("flash", url.QueryEscape(msg), 5, "/", "", false, true)
}
