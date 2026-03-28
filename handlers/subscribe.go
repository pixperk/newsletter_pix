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

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendJSON(w, http.StatusMethodNotAllowed, JSONResponse{Error: "Method not allowed"})
		return
	}

	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSON(w, http.StatusBadRequest, JSONResponse{Error: "Invalid JSON request"})
		return
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		SendJSON(w, http.StatusBadRequest, JSONResponse{Error: "Email is required"})
		return
	}

	var exists bool
	err := database.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM subscribers WHERE email = $1)`, email).Scan(&exists)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, JSONResponse{Error: "Database error occurred"})
		return
	}

	if exists {
		SendJSON(w, http.StatusConflict, JSONResponse{Error: "Email is already subscribed to the newsletter"})
		return
	}

	_, err = database.DB.Exec(`INSERT INTO subscribers (email) VALUES ($1)`, email)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, JSONResponse{Error: "Database error occurred"})
		return
	}

	SendJSON(w, http.StatusOK, JSONResponse{Success: true, Message: "Successfully subscribed to newsletter! 🎉"})
}
