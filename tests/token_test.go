package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/your-username/your-project/internal/token"
)

func TestTokenManagement(t *testing.T) {
	// Initialize token manager
	tokenManager := token.NewManager("your-secret-key")

	t.Run("Token Generation", func(t *testing.T) {
		userID := "user123"
		tokenString, err := tokenManager.GenerateToken(userID)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)
	})

	t.Run("Token Validation", func(t *testing.T) {
		userID := "user123"
		tokenString, _ := tokenManager.GenerateToken(userID)

		claims, err := tokenManager.ValidateToken(tokenString)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
	})

	t.Run("Token Expiration", func(t *testing.T) {
		userID := "user123"
		tokenString, _ := tokenManager.GenerateToken(userID)

		// Wait for token to expire (assuming short expiration for testing)
		time.Sleep(2 * time.Second)

		_, err := tokenManager.ValidateToken(tokenString)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})

	t.Run("Token Revocation", func(t *testing.T) {
		userID := "user123"
		tokenString, _ := tokenManager.GenerateToken(userID)

		err := tokenManager.RevokeToken(tokenString)
		assert.NoError(t, err)

		_, err = tokenManager.ValidateToken(tokenString)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "token has been revoked")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		invalidToken := "invalid.token.string"
		_, err := tokenManager.ValidateToken(invalidToken)
		assert.Error(t, err)
	})
}
