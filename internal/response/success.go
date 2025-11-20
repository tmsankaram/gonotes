package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}
type ListResponse struct {
	Items interface{} `json:"items"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Total int         `json:"total"`
}

func Success(c *gin.Context, message string, details interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Message: message,
		Details: details,
	})
}

func Created(c *gin.Context, message string, details interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Message: message,
		Details: details,
	})
}
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Accepted(c *gin.Context, message string, details interface{}) {
	c.JSON(http.StatusAccepted, SuccessResponse{
		Message: message,
		Details: details,
	})
}

func List(c *gin.Context, items interface{}, page, limit, total int) {
	c.JSON(http.StatusOK, ListResponse{
		Items: items,
		Page:  page,
		Limit: limit,
		Total: total,
	})
}
