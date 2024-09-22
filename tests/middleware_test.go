package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"your-project/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.AuthMiddleware())

	router.GET("/protected", func(c *gin.Context) {
		c.String(http.StatusOK, "protected")
	})

	t.Run("Valid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer valid_token")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "protected", w.Body.String())
	})

	t.Run("Invalid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Missing token", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RateLimitMiddleware(2, time.Second))

	router.GET("/limited", func(c *gin.Context) {
		c.String(http.StatusOK, "limited")
	})

	t.Run("Within rate limit", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/limited", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "limited", w.Body.String())
		}
	})

	t.Run("Exceeding rate limit", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/limited", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
	})

	t.Run("Rate limit reset", func(t *testing.T) {
		time.Sleep(time.Second)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/limited", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "limited", w.Body.String())
	})
}

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.CORSMiddleware())

	router.GET("/cors", func(c *gin.Context) {
		c.String(http.StatusOK, "cors")
	})

	t.Run("CORS headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/cors", nil)
		req.Header.Set("Origin", "http://example.com")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "cors", w.Body.String())
		assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET,POST,PUT,DELETE,OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Authorization,Content-Type", w.Header().Get("Access-Control-Allow-Headers"))
	})

	t.Run("OPTIONS request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/cors", nil)
		req.Header.Set("Origin", "http://example.com")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET,POST,PUT,DELETE,OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Authorization,Content-Type", w.Header().Get("Access-Control-Allow-Headers"))
	})
}
