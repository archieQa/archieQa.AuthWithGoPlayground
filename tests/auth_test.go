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

func TestMain(m *testing.M) {
	// Set up test database connection
	err := database.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer database.CloseDB()

	// Run tests
	m.Run()
}

func TestRegister(t *testing.T) {
	// Set up
	router := gin.Default()
	router.POST("/register", handlers.Register)

	// Test case 1: Valid registration
	t.Run("Valid Registration", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "Password123!",
		}
		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "User registered successfully")
	})

	// Test case 2: Invalid email
	t.Run("Invalid Email", func(t *testing.T) {
		user := models.User{
			Username: "testuser2",
			Email:    "invalidemail",
			Password: "Password123!",
		}
		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid email format")
	})

	// Add more test cases for weak password, existing username, etc.
}

func TestLogin(t *testing.T) {
	// Set up
	router := gin.Default()
	router.POST("/login", handlers.Login)

	// Test case 1: Valid login
	t.Run("Valid Login", func(t *testing.T) {
		// First, register a user
		user := models.User{
			Username: "loginuser",
			Email:    "login@example.com",
			Password: "Password123!",
		}
		hashedPassword, _ := utils.HashPassword(user.Password)
		database.GetDB().Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, hashedPassword)

		// Now, try to login
		loginData := map[string]string{
			"email":    user.Email,
			"password": user.Password,
		}
		jsonValue, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token")
	})

	// Test case 2: Invalid credentials
	t.Run("Invalid Credentials", func(t *testing.T) {
		loginData := map[string]string{
			"email":    "nonexistent@example.com",
			"password": "WrongPassword123!",
		}
		jsonValue, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")
	})
}

func TestJWTGeneration(t *testing.T) {
	userID := uint(1)
	token, err := utils.GenerateToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the token
	claims, err := utils.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}

func TestJWTValidation(t *testing.T) {
	userID := uint(1)
	token, _ := utils.GenerateToken(userID)

	// Test case 1: Valid token
	t.Run("Valid Token", func(t *testing.T) {
		err := utils.ValidateToken(token)
		assert.NoError(t, err)
	})

	// Test case 2: Invalid token
	t.Run("Invalid Token", func(t *testing.T) {
		err := utils.ValidateToken("invalid.token.here")
		assert.Error(t, err)
	})
}
