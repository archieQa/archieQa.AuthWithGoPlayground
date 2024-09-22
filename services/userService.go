package services

import (
	"errors"

	"github.com/your-project/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetProfile(userID string) (*models.User, error) {
	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) UpdateProfile(userID int, displayName, bio string) error {
	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	user.DisplayName = displayName
	user.Bio = bio

	return s.userRepository.Update(user)
}

func (s *UserService) ChangePassword(userID int, oldPassword, newPassword string) error {
	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepository.Update(user)
}

func (s *UserService) DeleteAccount(userID int) error {
	return s.userRepository.Delete(userID)
}

type UserRepository interface {
	GetByID(id interface{}) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}
