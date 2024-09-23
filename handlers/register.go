package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Temporary Store for users
var users = make(map[string]User)

// ResgisterHandler handels user Registration

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var user User

	if r.Body == nil {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}
	// This is a json Decoder that reads from r.Body the incoming request to the server for registration
	// .Decode(&user) this part will fill the User struct with data
	// It is err because if something goes wrong err will hold the info
	err := json.NewDecoder(r.Body).Decode(&user)

	log.Printf("Decoded User: %+v\n", user)
	log.Printf("Plaintext password: %s\n", user.Password)
	// This is where we check what went wrong
	if err != nil || user.Email == "" || user.Password == "" {
		log.Println("Error decoding JSON or missing fields:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Now we  Hash the User's password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// store the users with hashed passwords
	user.Password = string(hashedPassword)
	users[user.Email] = user

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User Registerd Succesfully"))
}
