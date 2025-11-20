package files

type File struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"-"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}
