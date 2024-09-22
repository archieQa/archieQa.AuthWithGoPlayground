package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/your-username/your-project/database"
	"github.com/your-username/your-project/handlers"
	"github.com/your-username/your-project/models"
	"github.com/your-username/your-project/utils"
)

func TestUpdateProfile(t *testing.T) {
	// Set up
	router := gin.Default()
	router.PUT("/profile", handlers.UpdateProfile)

	// Create a test user
	testUser := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	hashedPassword, _ := utils.HashPassword(testUser.Password)
	result := database.GetDB().Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id", testUser.Username, testUser.Email, hashedPassword)
	var userID uint
	result.Scan(&userID)

	t.Run("Valid Profile Update", func(t *testing.T) {
		updateData := map[string]string{
			"username": "updateduser",
			"email":    "updated@example.com",
		}
		jsonValue, _ := json.Marshal(updateData)
		req, _ := http.NewRequest("PUT", "/profile", bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", "Bearer "+generateTestToken(userID))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Profile updated successfully")

		// Verify database update
		var updatedUser models.User
		database.GetDB().QueryRow("SELECT username, email FROM users WHERE id = $1", userID).Scan(&updatedUser.Username, &updatedUser.Email)
		assert.Equal(t, "updateduser", updatedUser.Username)
		assert.Equal(t, "updated@example.com", updatedUser.Email)
	})

	t.Run("Invalid Email Update", func(t *testing.T) {
		updateData := map[string]string{
			"email": "invalid-email",
		}
		jsonValue, _ := json.Marshal(updateData)
		req, _ := http.NewRequest("PUT", "/profile", bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", "Bearer "+generateTestToken(userID))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid email format")
	})
}

func TestChangePassword(t *testing.T) {
	// Set up
	router := gin.Default()
	router.POST("/change-password", handlers.ChangePassword)

	// Create a test user
	testUser := models.User{
		Username: "passworduser",
		Email:    "password@example.com",
		Password: "OldPassword123!",
	}
	hashedPassword, _ := utils.HashPassword(testUser.Password)
	result := database.GetDB().Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id", testUser.Username, testUser.Email, hashedPassword)
	var userID uint
	result.Scan(&userID)

	t.Run("Valid Password Change", func(t *testing.T) {
		changeData := map[string]string{
			"old_password": "OldPassword123!",
			"new_password": "NewPassword456!",
		}
		jsonValue, _ := json.Marshal(changeData)
		req, _ := http.NewRequest("POST", "/change-password", bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", "Bearer "+generateTestToken(userID))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Password changed successfully")

		// Verify password change
		var newHashedPassword string
		database.GetDB().QueryRow("SELECT password FROM users WHERE id = $1", userID).Scan(&newHashedPassword)
		assert.True(t, utils.CheckPasswordHash("NewPassword456!", newHashedPassword))
	})

	t.Run("Incorrect Old Password", func(t *testing.T) {
		changeData := map[string]string{
			"old_password": "WrongOldPassword",
			"new_password": "NewPassword789!",
		}
		jsonValue, _ := json.Marshal(changeData)
		req, _ := http.NewRequest("POST", "/change-password", bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", "Bearer "+generateTestToken(userID))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Incorrect old password")
	})

	t.Run("Weak New Password", func(t *testing.T) {
		changeData := map[string]string{
			"old_password": "OldPassword123!",
			"new_password": "weak",
		}
		jsonValue, _ := json.Marshal(changeData)
		req, _ := http.NewRequest("POST", "/change-password", bytes.NewBuffer(jsonValue))
		req.Header.Set("Authorization", "Bearer "+generateTestToken(userID))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Password does not meet security requirements")
	})
}

func generateTestToken(userID uint) string {
	token, _ := utils.GenerateToken(userID)
	return token
}
