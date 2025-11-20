package notes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tmsankram/gonotes/internal/pagination"
	"github.com/tmsankram/gonotes/internal/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	notes := r.Group("/notes")
	{
		notes.GET("/", h.getAll)
		notes.GET("/:id", h.getByID)
		notes.POST("/", h.create)
		notes.PUT("/:id", h.update)
		notes.DELETE("/:id", h.delete)
	}
}

type NoteQuery struct {
	Title   string `form:"title"`
	Content string `form:"content"`
}

func (s *Service) Paginated(page, limit int) ([]Note, int64, error) {
	var notes []Note
	var total int64

	s.db.Model(&Note{}).Count(&total)

	offset := (page - 1) * limit

	err := s.db.Limit(limit).Offset(offset).Order("id DESC").Find(&notes).Error

	return notes, total, err
}

// ListNotes godoc
// @Summary List notes
// @Description Get paginated notes
// @Tags notes
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} response.ListResponse
// @Security ApiKeyAuth
// @Router /notes [get]
func (h *Handler) getAll(c *gin.Context) {
	var q NoteQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	var page pagination.Pagination

	if err := c.ShouldBindQuery(&page); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	page.Normalize()
	items, total, err := h.svc.Paginated(page.Page, page.Limit)
	if err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	response.List(c, items, page.Page, page.Limit, int(total))
}

func (h *Handler) getByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "invalid id")
		return
	}

	n, err := h.svc.GetByID(id)
	if err != nil {
		response.NotFound(c, errors.New("Note not found"))
		return
	}

	response.Success(c, "Note retrieved successfully", n)
}

type createUpdateReq struct {
	Title   string `json:"title" binding:"required,min=3,max=100,notest"`
	Content string `json:"content" binding:"required,min=5,max=5000"`
}

func (h *Handler) create(c *gin.Context) {
	var req createUpdateReq

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	n, err := h.svc.Create(Note{
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	response.Created(c, "Note created successfully", n)
}

func (h *Handler) update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "invalid id")
		return
	}
	var req createUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	n, err := h.svc.Update(id, Note{
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	response.Success(c, "Note updated successfully", n)
}

func (h *Handler) delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "invalid id")
		return
	}

	if err := h.svc.Delete(id); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
