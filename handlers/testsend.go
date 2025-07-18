package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/pixperk/newsletter/utils"
)

type TestSendRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type TestSendResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

const testRecipient = "mishrayashaswikumar@gmail.com" // Change to your email

func sendJSONTestSendResponse(w http.ResponseWriter, statusCode int, success bool, message string, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := TestSendResponse{
		Success: success,
		Message: message,
	}
	if errorMsg != "" {
		response.Error = errorMsg
	}
	json.NewEncoder(w).Encode(response)
}

func TestSendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONTestSendResponse(w, http.StatusMethodNotAllowed, false, "", "Method not allowed")
		return
	}

	if r.Header.Get("X-Secret") != os.Getenv("SEND_SECRET") {
		sendJSONTestSendResponse(w, http.StatusUnauthorized, false, "", "Unauthorized - invalid or missing X-Secret header")
		return
	}

	var req TestSendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONTestSendResponse(w, http.StatusBadRequest, false, "", "Invalid JSON request")
		return
	}

	if req.Subject == "" || req.Body == "" {
		sendJSONTestSendResponse(w, http.StatusBadRequest, false, "", "Subject and body are required")
		return
	}

	// Convert markdown to HTML
	htmlBody := utils.MarkdownToHTML(req.Body)

	// Optionally append the footer
	footerHTML := ""
	if f, err := os.ReadFile("footer.md"); err == nil {
		footerHTML = utils.MarkdownToHTML(string(f))
	}
	if footerHTML != "" {
		htmlBody += "<br>" + footerHTML
	}

	// Send email to test recipient
	if err := utils.SendEmail(testRecipient, req.Subject, htmlBody); err != nil {
		sendJSONTestSendResponse(w, http.StatusInternalServerError, false, "", "Failed to send test email")
		return
	}

	sendJSONTestSendResponse(w, http.StatusOK, true, "Test email sent to "+testRecipient, "")
}
