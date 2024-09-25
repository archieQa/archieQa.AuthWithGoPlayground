package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Temporary Store for users
var (
	users = make(map[string]User)
	mu    sync.Mutex
)

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Error decoding JSON: %v\n", err)
		http.Error(w, "Invalid JSON", http.StatusUnprocessableEntity)
		return
	}

	if user.Email == "" || user.Password == "" {
		log.Println("Missing email or password")
		http.Error(w, "Email and password are required", http.StatusUnprocessableEntity)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	user.Password = string(hashedPassword)

	mu.Lock()
	defer mu.Unlock()
	if _, exists := users[user.Email]; exists {
		log.Printf("User already exists: %s\n", user.Email)
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	users[user.Email] = user

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}
