package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Secret Key for JWT signing
var jwtKey = []byte("my_secret_key")

type Claims struct {
	Email string `json: "email"`
	jwt.StandardClaims
}

// LoginHandler handles user login

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds User
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Retrieve the stored User
	storedUser, exists := users[creds.Email]
	if !exists || bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(creds.Password)) != nil {
		http.Error(w, "Invalid Email or Password", http.StatusUnauthorized)
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

	log.Printf("This Email is successfully logged in: %s\n", creds.Email)

	// Send the JWt response

	w.Write([]byte(tokenString))

}
