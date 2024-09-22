package services

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"your-project/models"
)

type AuthService struct {
	userRepository UserRepository
}

func NewAuthService(userRepo UserRepository) *AuthService {
	return &AuthService{
		userRepository: userRepo,
	}
}

func (s *AuthService) Register(email, password string) error {
	// Check if user already exists
	if _, err := s.userRepository.GetByEmail(email); err == nil {
		return errors.New("user already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create new user
	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	// Save user to database
	return s.userRepository.Create(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	// Get user by email
	user, err := s.userRepository.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte("your-secret-key")) // Replace with your actual secret key
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
}
