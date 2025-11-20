package users

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`

	TOTPSecret  string `json:"-"` // base32 encoded TOTP secret
	TOTPEnabled bool   `json:"totp_enabled"`

	OAuthProvider string `json:"-"`
	OAuthID       string `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
