package services

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/your-project/models"
	"gorm.io/gorm"
)

var (
	secretKey = []byte("your-secret-key") // Replace with a secure secret key
)

type TokenService struct {
	DB *gorm.DB
}

type TokenClaims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func NewTokenService(db *gorm.DB) *TokenService {
	return &TokenService{DB: db}
}

func (s *TokenService) GenerateTokenPair(userID uint) (string, string, error) {
	// Generate access token
	accessToken, err := s.generateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *TokenService) generateAccessToken(userID uint) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func (s *TokenService) generateRefreshToken(userID uint) (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["user_id"] = userID
	rtClaims["exp"] = time.Now().Add(7 * 24 * time.Hour).Unix()

	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	// Store refresh token in the database
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err = models.CreateToken(s.DB, userID, refreshTokenString, expiresAt)
	if err != nil {
		return "", err
	}

	return refreshTokenString, nil
}

func (s *TokenService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *TokenService) RefreshToken(refreshToken string) (string, string, error) {
	// Validate the refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid token claims")
	}

	userID := uint(claims["user_id"].(float64))

	// Check if the refresh token exists in the database
	dbToken, err := models.GetTokenByString(s.DB, refreshToken)
	if err != nil {
		return "", "", errors.New("refresh token not found")
	}

	if dbToken.IsExpired() {
		return "", "", errors.New("refresh token has expired")
	}

	// Generate new token pair
	newAccessToken, newRefreshToken, err := s.GenerateTokenPair(userID)
	if err != nil {
		return "", "", err
	}

	// Revoke the old refresh token
	err = models.DeleteToken(s.DB, dbToken.ID)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *TokenService) RevokeToken(tokenString string) error {
	token, err := models.GetTokenByString(s.DB, tokenString)
	if err != nil {
		return err
	}

	return models.DeleteToken(s.DB, token.ID)
}
