package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	database "github.com/pixperk/newsletter/db"
)

type SubscribeRequest struct {
	Email string `json:"email"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, success bool, message string, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: success,
		Message: message,
	}
	if errorMsg != "" {
		response.Error = errorMsg
	}

	json.NewEncoder(w).Encode(response)
}

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONResponse(w, http.StatusMethodNotAllowed, false, "", "Method not allowed")
		return
	}

	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, false, "", "Invalid JSON request")
		return
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		sendJSONResponse(w, http.StatusBadRequest, false, "", "Email is required")
		return
	}

	// Check if email already exists
	var exists bool
	err := database.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM subscribers WHERE email = $1)`, email).Scan(&exists)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", "Database error occurred")
		return
	}

	if exists {
		sendJSONResponse(w, http.StatusConflict, false, "", "Email is already subscribed to the newsletter")
		return
	}

	// Insert new subscriber
	_, err = database.DB.Exec(`INSERT INTO subscribers (email) VALUES ($1)`, email)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, false, "", "Database error occurred")
		return
	}

	sendJSONResponse(w, http.StatusOK, true, "Successfully subscribed to newsletter! 🎉", "")
}
