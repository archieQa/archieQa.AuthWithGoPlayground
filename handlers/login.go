package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Secret Key for JWT signing
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// Claims struct for JWT payload
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Retrieve the stored User
	storedUser, exists := users[creds.Email]
	if !exists {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare the stored hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: creds.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	log.Printf("User successfully logged in: %s\n", creds.Email)

	// Send the JWT response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"token": tokenString}); err != nil {
		http.Error(w, "Could not encode response", http.StatusInternalServerError)
	}
}
