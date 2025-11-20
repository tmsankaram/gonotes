package files

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

type Service struct {
	mu    sync.RWMutex
	store map[string]File
}

func NewService() *Service {
	return &Service{
		store: make(map[string]File),
	}
}

func (s *Service) Save(header *multipart.FileHeader, file multipart.File) (File, error) {
	id := uuid.New().String()

	dstDir := "uploads"
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return File{}, err
	}

	// secure name
	name := filepath.Base(header.Filename)
	dstPath := filepath.Join(dstDir, id+"_"+name)

	out, err := os.Create(dstPath)
	if err != nil {
		return File{}, err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return File{}, err
	}

	f := File{
		ID:       id,
		Name:     name,
		Path:     dstPath,
		Size:     header.Size,
		MimeType: header.Header.Get("Content-Type"),
	}

	s.mu.Lock()
	s.store[id] = f
	s.mu.Unlock()

	return f, nil
}

func (s *Service) Get(id string) (File, error) {
	s.mu.RUnlock()
	defer s.mu.RUnlock()

	f, ok := s.store[id]
	if !ok {
		return File{}, errors.New("file not found")
	}
	return f, nil
}

func (s *Service) List() []File {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]File, 0, len(s.store))

	for _, f := range s.store {
		out = append(out, f)
	}
	return out
}
