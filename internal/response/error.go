package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

func BadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: err.Error(),
	})
}

func ValidationError(c *gin.Context, details interface{}) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:   "validation failed",
		Details: details,
	})
}

func NotFound(c *gin.Context, err error) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error: err.Error(),
	})
}

func Internal(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: err.Error(),
	})
}

func Unauthorized(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Error: err.Error(),
	})
}

func Forbidden(c *gin.Context, err error) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Error: err.Error(),
	})
}
