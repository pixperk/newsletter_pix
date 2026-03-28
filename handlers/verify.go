package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	database "github.com/pixperk/newsletter/db"
	"github.com/pixperk/newsletter/utils"
)

type VerifySubscribeRequest struct {
	Email string `json:"email"`
}

type VerifyConfirmRequest struct {
	Token string `json:"token"`
}

func generateVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// VerifySubscribeHandler - Step 1: Send verification email
func VerifySubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendJSON(w, http.StatusMethodNotAllowed, JSONResponse{Error: "Method not allowed"})
		return
	}

	var req VerifySubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSON(w, http.StatusBadRequest, JSONResponse{Error: "Invalid JSON request"})
		return
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		SendJSON(w, http.StatusBadRequest, JSONResponse{Error: "Email is required"})
		return
	}

	var existsInSubscribers bool
	err := database.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM subscribers WHERE email = $1)`, email).Scan(&existsInSubscribers)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, JSONResponse{Error: "Database error occurred"})
		return
	}

	if existsInSubscribers {
		SendJSON(w, http.StatusConflict, JSONResponse{Error: "Email is already subscribed to the newsletter"})
		return
	}

	token, err := generateVerificationToken()
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, JSONResponse{Error: "Failed to generate verification token"})
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = database.DB.Exec(`DELETE FROM email_verifications WHERE email = $1`, email)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, JSONResponse{Error: "Database error occurred"})
		return
	}

	_, err = database.DB.Exec(`
		INSERT INTO email_verifications (email, verification_token, expires_at)
		VALUES ($1, $2, $3)
	`, email, token, expiresAt)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, JSONResponse{Error: "Database error occurred"})
		return
	}

	verificationURL := fmt.Sprintf("https://pixperk.tech?verify=%s", token)

	emailSubject := "Verify your newsletter subscription"
	emailBody := fmt.Sprintf(`
# Email Verification Required

Hi there! 👋

Thanks for subscribing to the PixPerk newsletter. To complete your subscription, please click the link below to verify your email address:

**[Verify Email Address](%s)**

This link will expire in 24 hours. If you didn't request this subscription, you can safely ignore this email.

---

Best regards,
**yashaswi.**
`, verificationURL)

	htmlBody := utils.MarkdownToHTML(emailBody)

	footerHTML := ""
	if f, err := os.ReadFile("footer.md"); err == nil {
		footerHTML = utils.MarkdownToHTML(string(f))
	}
	if footerHTML != "" {
		htmlBody += "<br>" + footerHTML
	}

	if err := utils.SendEmail(email, emailSubject, htmlBody); err != nil {
		SendJSON(w, http.StatusInternalServerError, JSONResponse{Error: "Failed to send verification email: " + err.Error()})
		return
	}

	SendJSON(w, http.StatusOK, JSONResponse{Success: true, Message: "Verification email sent! Please check your inbox and click the verification link."})
}

// VerifyConfirmHandler - Step 2: Confirm verification and add to subscribers
func VerifyConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendJSON(w, http.StatusMethodNotAllowed, JSONResponse{Error: "Method not allowed"})
		return
	}

	var req VerifyConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSON(w, http.StatusBadRequest, JSONResponse{Error: "Invalid JSON request"})
		return
	}

	token := strings.TrimSpace(req.Token)
	if token == "" {
		SendJSON(w, http.StatusBadRequest, JSONResponse{Error: "Verification token is required"})
		return
	}

	var email string
	var expiresAt time.Time
	var verified bool
	err := database.DB.QueryRow(`
		SELECT email, expires_at, verified
		FROM email_verifications
		WHERE verification_token = $1
	`, token).Scan(&email, &expiresAt, &verified)

	if err != nil {
		SendJSON(w, http.StatusNotFound, JSONResponse{Error: "Invalid or expired verification token"})
		return
	}

	if time.Now().After(expiresAt) {
		SendJSON(w, http.StatusGone, JSONResponse{Error: "Verification token has expired"})
		return
	}

	if verified {
		SendJSON(w, http.StatusConflict, JSONResponse{Error: "Email has already been verified"})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, JSONResponse{Error: "Database error occurred"})
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE email_verifications
		SET verified = TRUE
		WHERE verification_token = $1
	`, token)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, JSONResponse{Error: "Database error occurred"})
		return
	}

	_, err = tx.Exec(`
		INSERT INTO subscribers (email)
		VALUES ($1)
		ON CONFLICT (email) DO NOTHING
	`, email)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, JSONResponse{Error: "Database error occurred"})
		return
	}

	if err = tx.Commit(); err != nil {
		SendJSON(w, http.StatusInternalServerError, JSONResponse{Error: "Database error occurred"})
		return
	}

	SendJSON(w, http.StatusOK, JSONResponse{Success: true, Message: "Email verified successfully! You're now subscribed to the newsletter. 🎉"})
}
