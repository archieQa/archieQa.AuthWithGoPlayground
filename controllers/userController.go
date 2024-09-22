package controllers

import (
	"encoding/json"
	"net/http"

	"auth_go/services"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var profileUpdate struct {
		UserID      int    `json:"user_id"`
		DisplayName string `json:"display_name"`
		Bio         string `json:"bio"`
	}

	if err := json.NewDecoder(r.Body).Decode(&profileUpdate); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := uc.userService.UpdateProfile(profileUpdate.UserID, profileUpdate.DisplayName, profileUpdate.Bio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}

func (uc *UserController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var passwordChange struct {
		UserID      int    `json:"user_id"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&passwordChange); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := uc.userService.ChangePassword(passwordChange.UserID, passwordChange.OldPassword, passwordChange.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password changed successfully"})
}

func (uc *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	profile, err := uc.userService.GetProfile(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}
