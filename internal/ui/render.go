package ui

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Renderer struct {
	T *template.Template
}

func NewRenderer(t *template.Template) *Renderer {
	return &Renderer{T: t}
}

func (r *Renderer) Page(c *gin.Context, name string, data gin.H) {
	// auto-add user + flash from context
	if user, exists := c.Get("user"); exists {
		data["user"] = user
	}
	if flash, exists := c.Get("flash"); exists {
		data["Flash"] = flash
	}

	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := r.T.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}
