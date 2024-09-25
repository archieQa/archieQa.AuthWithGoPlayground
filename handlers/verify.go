package handlers

// This code imports necessary packages for handling HTTP requests and JWT (JSON Web Token)
// authentication in a Go application. The "net/http" package is used for HTTP client and
// server implementations, while the "strings" package provides functions to manipulate
// UTF-8 encoded strings. The "github.com/dgrijalva/jwt-go" package is a third-party
// library used for creating and parsing JWTs, which are commonly used for securely
// transmitting information between parties as a JSON object.

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Ensure jwtKey is declared only once in the appropriate scope
// var jwtKey = []byte("your_secret_key")

// Claims struct to read the JWT claims
type JWTClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// VerifyHandler handles the verification of JWT tokens
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	// Split the bearer token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		http.Error(w, "Invalid Authorization format", http.StatusUnauthorized)
		return
	}

	// Extract the actual token
	tokenStr := tokenParts[1]

	claims := &JWTClaims{
		Email:          "",
		StandardClaims: jwt.StandardClaims{},
	}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return jwtKey, nil
	})

	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("Token is valid"))
}
