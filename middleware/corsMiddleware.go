package middleware

import (
	"net/http"
)

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS) settings
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set allowed origins
		// In production, replace "*" with your specific domain(s)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Set allowed methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Set allowed headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	}
}
