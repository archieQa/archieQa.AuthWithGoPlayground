package main

import (
	"log"
	"net/http"

	"github.com/arsi/auth_service/handlers"
)

func main() {
	// Simple handler to test the server
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is up and running"))
	})

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/verify", handlers.VerifyHandler)
	// Start the server on port 8080
	log.Println("Starting Server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}

}
