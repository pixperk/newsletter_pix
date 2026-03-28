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

const testRecipient = "mishrayashaswikumar@gmail.com"

func TestSendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendJSON(w, http.StatusMethodNotAllowed, JSONResponse{Error: "Method not allowed"})
		return
	}

	if r.Header.Get("X-Secret") != os.Getenv("SEND_SECRET") {
		SendJSON(w, http.StatusUnauthorized, JSONResponse{Error: "Unauthorized - invalid or missing X-Secret header"})
		return
	}

	var req TestSendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSON(w, http.StatusBadRequest, JSONResponse{Error: "Invalid JSON request"})
		return
	}

	if req.Subject == "" || req.Body == "" {
		SendJSON(w, http.StatusBadRequest, JSONResponse{Error: "Subject and body are required"})
		return
	}

	htmlBody := utils.MarkdownToHTML(req.Body)

	footerHTML := ""
	if f, err := os.ReadFile("footer.md"); err == nil {
		footerHTML = utils.MarkdownToHTML(string(f))
	}
	if footerHTML != "" {
		htmlBody += "<br>" + footerHTML
	}

	if err := utils.SendEmail(testRecipient, req.Subject, htmlBody); err != nil {
		SendJSON(w, http.StatusInternalServerError, JSONResponse{Error: "Failed to send test email: " + err.Error()})
		return
	}

	SendJSON(w, http.StatusOK, JSONResponse{Success: true, Message: "Test email sent to " + testRecipient})
}
