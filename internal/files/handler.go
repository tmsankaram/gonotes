package files

import (
	"errors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tmsankram/gonotes/internal/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	g := r.Group("/files")
	{
		g.POST("/upload", h.upload)
		g.GET("/", h.list)
		g.GET("/:id/download", h.download)
	}
}

// UploadFile godoc
// @Summary Upload file
// @Description Upload a file
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "file to upload"
// @Success 201 {object} files.File
// @Failure 400 {object} response.ErrorResponse
// @Security ApiKeyAuth
// @Router /files/upload [post]
func (h *Handler) upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, errors.New("file required"))
		return
	}
	f, err := fileHeader.Open()

	if err != nil {
		response.BadRequest(c, errors.New("cannot read file"))
		return
	}

	defer f.Close()

	result, err := h.svc.Save(fileHeader, f)
	if err != nil {
		response.Internal(c, errors.New(err.Error()))
		return
	}

	response.Created(c, "file uploaded successfully", result)
}

func (h *Handler) list(c *gin.Context) {
	response.Success(c, "files fetched successfully", h.svc.List())
}

func (h *Handler) download(c *gin.Context) {
	id := c.Param("id")
	f, err := h.svc.Get(id)
	if err != nil {
		response.NotFound(c, errors.New("file not found"))
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+f.Name)
	c.Header("Content-Type", f.MimeType)

	fp, err := os.Open(f.Path)
	if err != nil {
		response.Internal(c, errors.New("cannot read file"))
		return
	}
	defer fp.Close()
	http.ServeContent(c.Writer, c.Request, f.Name, os.FileInfo.ModTime(nil), fp)
}
