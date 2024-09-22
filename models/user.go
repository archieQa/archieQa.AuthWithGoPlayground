package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// User represents the user model in the database
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password  string    `json:"-" gorm:"not null" validate:"required,min=8"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Validate checks if the User struct is valid
func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
