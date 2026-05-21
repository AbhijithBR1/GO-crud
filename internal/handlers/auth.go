package handlers

import (
	"encoding/json"
	"net/http"

	"bookmanagement/internal/models"
)

// HandleRegister corresponds to POST /register
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := models.RegisterUser(reqBody.Username, reqBody.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict) // 409 Conflict
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// HandleLogin corresponds to POST /login
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := models.LoginUser(reqBody.Username, reqBody.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized) // 401 Unauthorized
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// We return the token as JSON: { "token": "abc123def456" }
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
