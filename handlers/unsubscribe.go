package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	database "github.com/pixperk/newsletter/db"
)

type UnsubscribeRequest struct {
	Email string `json:"email"`
}

type UnsubscribeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func sendJSONUnsubscribeResponse(w http.ResponseWriter, statusCode int, success bool, message string, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := UnsubscribeResponse{
		Success: success,
		Message: message,
	}
	if errorMsg != "" {
		response.Error = errorMsg
	}

	json.NewEncoder(w).Encode(response)
}

func UnsubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONUnsubscribeResponse(w, http.StatusMethodNotAllowed, false, "", "Method not allowed")
		return
	}

	var req UnsubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONUnsubscribeResponse(w, http.StatusBadRequest, false, "", "Invalid JSON request")
		return
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		sendJSONUnsubscribeResponse(w, http.StatusBadRequest, false, "", "Email is required")
		return
	}

	// Check if email exists in subscribers
	var exists bool
	err := database.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM subscribers WHERE email = $1)`, email).Scan(&exists)
	if err != nil {
		sendJSONUnsubscribeResponse(w, http.StatusInternalServerError, false, "", "Database error occurred")
		return
	}

	if !exists {
		sendJSONUnsubscribeResponse(w, http.StatusNotFound, false, "", "Email not found in subscription list")
		return
	}

	// Remove subscriber
	_, err = database.DB.Exec(`DELETE FROM subscribers WHERE email = $1`, email)
	if err != nil {
		sendJSONUnsubscribeResponse(w, http.StatusInternalServerError, false, "", "Database error occurred")
		return
	}

	sendJSONUnsubscribeResponse(w, http.StatusOK, true, "Successfully unsubscribed from newsletter 👋", "")
}
