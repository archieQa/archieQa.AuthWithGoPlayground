package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	reqBody := bytes.NewBufferString(`{"email": "test@example.com", "password": "password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", reqBody)
	res := httptest.NewRecorder()

	RegisterHandler(res, req)

	if res.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, res.Code)
	}

	expected := "User registered successfully"
	if res.Body.String() != expected {
		t.Errorf("Expected body %q, got %q", expected, res.Body.String())
	}
}
