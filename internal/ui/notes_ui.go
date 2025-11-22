package ui

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/tmsankram/gonotes/internal/notes"
)

type NotesUI struct {
	Notes    *notes.Service
	Renderer *Renderer
}

func NewNotesUI(n *notes.Service, r *Renderer) *NotesUI {
	return &NotesUI{
		Notes:    n,
		Renderer: r,
	}
}

// GET /notes
func (h *NotesUI) NotesPage(c *gin.Context) {
	notes, _, err := h.Notes.Paginated(1, 9999) // load all for UI
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	h.Renderer.Page(c, "notes/list.html", gin.H{
		"Title": "Notes",
		"Notes": notes,
	})
}

// GET /notes/create-form
func (h *NotesUI) CreateForm(c *gin.Context) {
	h.Renderer.Page(c, "notes/create.html", gin.H{})
}

// POST /notes/create
func (h *NotesUI) CreatePost(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")

	n, err := h.Notes.Create(notes.Note{
		Title:   title,
		Content: content,
	})
	if err != nil {
		h.Renderer.Page(c, "notes/create.html", gin.H{
			"Flash": "Error creating note",
		})
		return
	}

	// return new item partial for HTMX prepend
	h.Renderer.Page(c, "notes/item.html", gin.H{
		"Note": n,
	})
}

// GET /notes/:id/edit
func (h *NotesUI) EditForm(c *gin.Context) {
	id := c.Param("id")
	nid, _ := strconv.Atoi(id)

	n, err := h.Notes.GetByID(int64(nid))
	if err != nil {
		c.String(404, "Note not found")
		return
	}

	h.Renderer.Page(c, "notes/edit.html", gin.H{
		"Note": n,
	})
}

// POST /notes/:id/edit
func (h *NotesUI) EditPost(c *gin.Context) {
	id := c.Param("id")
	nid, _ := strconv.Atoi(id)

	title := c.PostForm("title")
	content := c.PostForm("content")

	n, err := h.Notes.Update(int64(nid), notes.Note{
		Title:   title,
		Content: content,
	})
	if err != nil {
		h.Renderer.Page(c, "notes/edit.html", gin.H{
			"Flash": "Update failed",
			"Note":  n,
		})
		return
	}

	// return updated list item partial
	h.Renderer.Page(c, "notes/item.html", gin.H{
		"Note": n,
	})
}

// DELETE /notes/:id
func (h *NotesUI) Delete(c *gin.Context) {
	id := c.Param("id")
	nid, _ := strconv.Atoi(id)

	err := h.Notes.Delete(int64(nid))
	if err != nil {
		c.String(400, "Delete failed")
		return
	}

	c.Status(http.StatusOK)
}
