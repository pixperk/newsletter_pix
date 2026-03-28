package handlers

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Success         bool   `json:"success"`
	Message         string `json:"message,omitempty"`
	Error           string `json:"error,omitempty"`
	EmailsSent      int    `json:"emails_sent,omitempty"`
	SubscriberCount int    `json:"subscriber_count,omitempty"`
}

func SendJSON(w http.ResponseWriter, statusCode int, resp JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
