package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pixperk/newsletter/utils"
)

type SendRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"` // Markdown input
}

type SendResponse struct {
	Success         bool   `json:"success"`
	Message         string `json:"message"`
	Error           string `json:"error,omitempty"`
	EmailsSent      int    `json:"emails_sent,omitempty"`
	SubscriberCount int    `json:"subscriber_count,omitempty"`
}

func sendJSONSendResponse(w http.ResponseWriter, statusCode int, success bool, message string, errorMsg string, emailsSent int, subscriberCount int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := SendResponse{
		Success: success,
		Message: message,
	}
	if errorMsg != "" {
		response.Error = errorMsg
	}
	if emailsSent > 0 {
		response.EmailsSent = emailsSent
	}
	if subscriberCount > 0 {
		response.SubscriberCount = subscriberCount
	}

	json.NewEncoder(w).Encode(response)
}

// loadFooterHTML reads footer.md and converts it to HTML
func loadFooterHTML() string {
	// Try to read footer.md from current directory
	footerPath := "footer.md"
	if _, err := os.Stat(footerPath); os.IsNotExist(err) {
		// Try relative path from handlers directory
		footerPath = filepath.Join("..", "footer.md")
		if _, err := os.Stat(footerPath); os.IsNotExist(err) {
			// Try absolute path to project root
			footerPath = filepath.Join(".", "footer.md")
			if _, err := os.Stat(footerPath); os.IsNotExist(err) {
				return ""
			}
		}
	}

	footerContent, err := os.ReadFile(footerPath)
	if err != nil {
		return ""
	}

	// Convert markdown to HTML
	footerHTML := utils.MarkdownToHTML(string(footerContent))

	// Add some spacing for the footer
	return "<br>" + footerHTML
}

func SendHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			sendJSONSendResponse(w, http.StatusMethodNotAllowed, false, "", "Method not allowed", 0, 0)
			return
		}

		if r.Header.Get("X-Secret") != os.Getenv("SEND_SECRET") {
			sendJSONSendResponse(w, http.StatusUnauthorized, false, "", "Unauthorized - invalid or missing X-Secret header", 0, 0)
			return
		}

		var req SendRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendJSONSendResponse(w, http.StatusBadRequest, false, "", "Invalid JSON request", 0, 0)
			return
		}

		if req.Subject == "" || req.Body == "" {
			sendJSONSendResponse(w, http.StatusBadRequest, false, "", "Subject and body are required", 0, 0)
			return
		}

		// Convert main content to HTML
		htmlBody := utils.MarkdownToHTML(req.Body)

		// Append footer if available
		footerHTML := loadFooterHTML()
		if footerHTML != "" {
			htmlBody += footerHTML
		}

		rows, err := db.Query(`SELECT email FROM subscribers`)
		if err != nil {
			sendJSONSendResponse(w, http.StatusInternalServerError, false, "", "Database error occurred", 0, 0)
			return
		}
		defer rows.Close()

		var emails []string
		for rows.Next() {
			var email string
			rows.Scan(&email)
			emails = append(emails, email)
		}

		subscriberCount := len(emails)
		if subscriberCount == 0 {
			sendJSONSendResponse(w, http.StatusOK, true, "No subscribers to send emails to", "", 0, 0)
			return
		}

		// Send emails asynchronously
		for _, email := range emails {
			go utils.SendEmail(email, req.Subject, htmlBody)
		}

		sendJSONSendResponse(w, http.StatusOK, true, "Newsletter sent successfully! 🚀", "", subscriberCount, subscriberCount)
	}
}
