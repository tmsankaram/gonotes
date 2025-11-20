package notes

import (
	"errors"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetByID(id int64) (Note, error) {
	var n Note
	err := s.db.First(&n, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return n, err
	}
	return n, err
}

func (s *Service) Create(n Note) (Note, error) {
	if err := s.db.Create(&n).Error; err != nil {
		return Note{}, err
	}
	return n, nil
}

func (s *Service) Update(id int64, data Note) (Note, error) {
	var n Note
	if err := s.db.First(&n, id).Error; err != nil {
		return Note{}, err
	}

	n.Title = data.Title
	n.Content = data.Content

	if err := s.db.Save(&n).Error; err != nil {
		return Note{}, err
	}

	return n, nil
}

func (s *Service) Delete(id int64) error {
	if err := s.db.Delete(&Note{}, id).Error; err != nil {
		return err
	}
	return nil
}
