package users

import (
	"errors"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(u User) (User, error) {
	if err := s.db.Create(&u).Error; err != nil {
		return User{}, err
	}
	return u, nil
}

func (s *Service) GetByEmail(email string) (User, error) {
	var u User
	err := s.db.Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, nil
	}
	return u, err
}

func (s *Service) GetByID(id uint) (User, error) {
	var u User
	err := s.db.First(&u, id).Error
	return u, err
}

func (s *Service) Update(u User) error {
	return s.db.Save(&u).Error
}

func (s *Service) Delete(id uint) error {
	return s.db.Delete(&User{}, id).Error
}

func (s *Service) SetTOTP(userID uint, secret string) error {
	return s.db.Model(&User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"totp_secret":  secret,
		"totp_enabled": true,
	}).Error
}

func (s *Service) GetByOauth(provider, oauthID string) (User, error) {
	var u User
	err := s.db.Where("oauth_provider = ? AND oauth_id = ?", provider, oauthID).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, nil
	}
	return u, err
}

func (s *Service) CreateOAuthUser(email, provider, oauthID string) (User, error) {
	u := User{
		Email:         email,
		OAuthProvider: provider,
		OAuthID:       oauthID,
	}
	if err := s.db.Create(&u).Error; err != nil {
		return User{}, err
	}
	return u, nil
}
