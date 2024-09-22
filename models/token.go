package models

import (
	"time"

	"gorm.io/gorm"
)

// Token represents a refresh token stored in the database
type Token struct {
	gorm.Model
	UserID      uint      `gorm:"not null"`
	TokenString string    `gorm:"unique;not null"`
	ExpiresAt   time.Time `gorm:"not null"`
}

// CreateToken creates a new token for the given user
func CreateToken(db *gorm.DB, userID uint, tokenString string, expiresAt time.Time) (*Token, error) {
	token := &Token{
		UserID:      userID,
		TokenString: tokenString,
		ExpiresAt:   expiresAt,
	}

	result := db.Create(token)
	if result.Error != nil {
		return nil, result.Error
	}

	return token, nil
}

// GetTokenByString retrieves a token by its string value
func GetTokenByString(db *gorm.DB, tokenString string) (*Token, error) {
	var token Token
	result := db.Where("token_string = ?", tokenString).First(&token)
	if result.Error != nil {
		return nil, result.Error
	}

	return &token, nil
}

// DeleteToken removes a token from the database
func DeleteToken(db *gorm.DB, tokenID uint) error {
	result := db.Delete(&Token{}, tokenID)
	return result.Error
}

// IsExpired checks if the token has expired
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
