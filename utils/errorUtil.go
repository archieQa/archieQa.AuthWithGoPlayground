package utils

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a standardized error response structure
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// NewErrorResponse creates a new ErrorResponse
func NewErrorResponse(status int, message string) *ErrorResponse {
	return &ErrorResponse{
		Status:  status,
		Message: message,
	}
}

// RespondWithError sends a JSON error response
func RespondWithError(w http.ResponseWriter, status int, message string) {
	errorResponse := NewErrorResponse(status, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse)
}

// Common error messages
const (
	ErrInvalidRequest     = "Invalid request"
	ErrInternalServer     = "Internal server error"
	ErrUnauthorized       = "Unauthorized"
	ErrNotFound           = "Resource not found"
	ErrAlreadyExists      = "Resource already exists"
	ErrInvalidCredentials = "Invalid credentials"
)

// HandleInternalServerError logs the error and sends a 500 Internal Server Error response
func HandleInternalServerError(w http.ResponseWriter, err error) {
	// TODO: Add proper logging here
	RespondWithError(w, http.StatusInternalServerError, ErrInternalServer)
}

// HandleBadRequest sends a 400 Bad Request response with the given message
func HandleBadRequest(w http.ResponseWriter, message string) {
	RespondWithError(w, http.StatusBadRequest, message)
}

// HandleUnauthorized sends a 401 Unauthorized response
func HandleUnauthorized(w http.ResponseWriter) {
	RespondWithError(w, http.StatusUnauthorized, ErrUnauthorized)
}

// HandleNotFound sends a 404 Not Found response
func HandleNotFound(w http.ResponseWriter) {
	RespondWithError(w, http.StatusNotFound, ErrNotFound)
}

// HandleConflict sends a 409 Conflict response
func HandleConflict(w http.ResponseWriter, message string) {
	RespondWithError(w, http.StatusConflict, message)
}
