package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/pixperk/newsletter/utils"
)

type SendRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

const maxSendWorkers = 20

func SendHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			SendJSON(w, http.StatusMethodNotAllowed, JSONResponse{Error: "Method not allowed"})
			return
		}

		if r.Header.Get("X-Secret") != os.Getenv("SEND_SECRET") {
			SendJSON(w, http.StatusUnauthorized, JSONResponse{Error: "Unauthorized - invalid or missing X-Secret header"})
			return
		}

		var req SendRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			SendJSON(w, http.StatusBadRequest, JSONResponse{Error: "Invalid JSON request"})
			return
		}

		if req.Subject == "" || req.Body == "" {
			SendJSON(w, http.StatusBadRequest, JSONResponse{Error: "Subject and body are required"})
			return
		}

		htmlBody := utils.MarkdownToHTML(req.Body)

		footerHTML := loadFooterHTML()
		if footerHTML != "" {
			htmlBody += footerHTML
		}

		rows, err := db.Query(`SELECT email FROM subscribers`)
		if err != nil {
			SendJSON(w, http.StatusInternalServerError, JSONResponse{Error: "Database error occurred"})
			return
		}
		defer rows.Close()

		var emails []string
		for rows.Next() {
			var email string
			if err := rows.Scan(&email); err != nil {
				continue
			}
			emails = append(emails, email)
		}

		subscriberCount := len(emails)
		if subscriberCount == 0 {
			SendJSON(w, http.StatusOK, JSONResponse{Success: true, Message: "No subscribers to send emails to"})
			return
		}

		// Worker pool to bound concurrency
		type result struct {
			email string
			err   error
		}

		jobs := make(chan string, len(emails))
		results := make(chan result, len(emails))

		var wg sync.WaitGroup
		workers := maxSendWorkers
		if subscriberCount < workers {
			workers = subscriberCount
		}
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for email := range jobs {
					err := utils.SendEmail(email, req.Subject, htmlBody)
					results <- result{email: email, err: err}
				}
			}()
		}

		for _, email := range emails {
			jobs <- email
		}
		close(jobs)

		go func() {
			wg.Wait()
			close(results)
		}()

		successCount := 0
		errorCount := 0
		var errors []string

		for res := range results {
			if res.err != nil {
				errorCount++
				errors = append(errors, fmt.Sprintf("Failed to send to %s: %v", res.email, res.err))
			} else {
				successCount++
			}
		}

		if errorCount == 0 {
			SendJSON(w, http.StatusOK, JSONResponse{Success: true, Message: "Newsletter sent successfully to all subscribers! 🚀", EmailsSent: successCount, SubscriberCount: subscriberCount})
		} else if successCount > 0 {
			errorMessage := fmt.Sprintf("Partial success: %d sent, %d failed. Errors: %s", successCount, errorCount, strings.Join(errors, "; "))
			SendJSON(w, http.StatusPartialContent, JSONResponse{Success: true, Message: "Newsletter partially sent", Error: errorMessage, EmailsSent: successCount, SubscriberCount: subscriberCount})
		} else {
			errorMessage := fmt.Sprintf("All emails failed to send. Errors: %s", strings.Join(errors, "; "))
			SendJSON(w, http.StatusInternalServerError, JSONResponse{Error: errorMessage, EmailsSent: successCount, SubscriberCount: subscriberCount})
		}
	}
}

func loadFooterHTML() string {
	footerContent, err := os.ReadFile("footer.md")
	if err != nil {
		return ""
	}
	return "<br>" + utils.MarkdownToHTML(string(footerContent))
}
