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

type VerifyResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func generateVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func sendJSONVerifyResponse(w http.ResponseWriter, statusCode int, success bool, message string, errorMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := VerifyResponse{
		Success: success,
		Message: message,
	}
	if errorMsg != "" {
		response.Error = errorMsg
	}

	json.NewEncoder(w).Encode(response)
}

// VerifySubscribeHandler - Step 1: Send verification email
func VerifySubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONVerifyResponse(w, http.StatusMethodNotAllowed, false, "", "Method not allowed")
		return
	}

	var req VerifySubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONVerifyResponse(w, http.StatusBadRequest, false, "", "Invalid JSON request")
		return
	}

	email := strings.TrimSpace(req.Email)
	if email == "" {
		sendJSONVerifyResponse(w, http.StatusBadRequest, false, "", "Email is required")
		return
	}

	// Check if email is already subscribed
	var existsInSubscribers bool
	err := database.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM subscribers WHERE email = $1)`, email).Scan(&existsInSubscribers)
	if err != nil {
		sendJSONVerifyResponse(w, http.StatusInternalServerError, false, "", "Database error occurred")
		return
	}

	if existsInSubscribers {
		sendJSONVerifyResponse(w, http.StatusConflict, false, "", "Email is already subscribed to the newsletter")
		return
	}

	// Generate verification token
	token, err := generateVerificationToken()
	if err != nil {
		sendJSONVerifyResponse(w, http.StatusInternalServerError, false, "", "Failed to generate verification token")
		return
	}

	// Set expiration time (24 hours from now)
	expiresAt := time.Now().Add(24 * time.Hour)

	// Delete any existing verification for this email
	_, err = database.DB.Exec(`DELETE FROM email_verifications WHERE email = $1`, email)
	if err != nil {
		sendJSONVerifyResponse(w, http.StatusInternalServerError, false, "", "Database error occurred")
		return
	}

	// Insert verification record
	_, err = database.DB.Exec(`
		INSERT INTO email_verifications (email, verification_token, expires_at) 
		VALUES ($1, $2, $3)
	`, email, token, expiresAt)
	if err != nil {
		sendJSONVerifyResponse(w, http.StatusInternalServerError, false, "", "Database error occurred")
		return
	}

	// Create verification email content
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

	// Convert markdown to HTML and add footer
	htmlBody := utils.MarkdownToHTML(emailBody)

	// Load and append footer
	footerHTML := ""
	if f, err := os.ReadFile("footer.md"); err == nil {
		footerHTML = utils.MarkdownToHTML(string(f))
	}
	if footerHTML != "" {
		htmlBody += "<br>" + footerHTML
	}

	// Send verification email
	if err := utils.SendEmail(email, emailSubject, htmlBody); err != nil {
		sendJSONVerifyResponse(w, http.StatusInternalServerError, false, "", "Failed to send verification email: "+err.Error())
		return
	}

	sendJSONVerifyResponse(w, http.StatusOK, true, "Verification email sent! Please check your inbox and click the verification link.", "")
}

// VerifyConfirmHandler - Step 2: Confirm verification and add to subscribers
func VerifyConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONVerifyResponse(w, http.StatusMethodNotAllowed, false, "", "Method not allowed")
		return
	}

	var req VerifyConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONVerifyResponse(w, http.StatusBadRequest, false, "", "Invalid JSON request")
		return
	}

	token := strings.TrimSpace(req.Token)
	if token == "" {
		sendJSONVerifyResponse(w, http.StatusBadRequest, false, "", "Verification token is required")
		return
	}

	// Check if token exists and is not expired
	var email string
	var expiresAt time.Time
	var verified bool
	err := database.DB.QueryRow(`
		SELECT email, expires_at, verified 
		FROM email_verifications 
		WHERE verification_token = $1
	`, token).Scan(&email, &expiresAt, &verified)

	if err != nil {
		sendJSONVerifyResponse(w, http.StatusNotFound, false, "", "Invalid or expired verification token")
		return
	}

	// Check if token has expired
	if time.Now().After(expiresAt) {
		sendJSONVerifyResponse(w, http.StatusGone, false, "", "Verification token has expired")
		return
	}

	// Check if already verified
	if verified {
		sendJSONVerifyResponse(w, http.StatusConflict, false, "", "Email has already been verified")
		return
	}

	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		sendJSONVerifyResponse(w, http.StatusInternalServerError, false, "", "Database error occurred")
		return
	}
	defer tx.Rollback()

	// Mark as verified
	_, err = tx.Exec(`
		UPDATE email_verifications 
		SET verified = TRUE 
		WHERE verification_token = $1
	`, token)
	if err != nil {
		sendJSONVerifyResponse(w, http.StatusInternalServerError, false, "", "Database error occurred")
		return
	}

	// Add to subscribers table
	_, err = tx.Exec(`
		INSERT INTO subscribers (email) 
		VALUES ($1) 
		ON CONFLICT (email) DO NOTHING
	`, email)
	if err != nil {
		sendJSONVerifyResponse(w, http.StatusInternalServerError, false, "", "Database error occurred")
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		sendJSONVerifyResponse(w, http.StatusInternalServerError, false, "", "Database error occurred")
		return
	}

	sendJSONVerifyResponse(w, http.StatusOK, true, "Email verified successfully! You're now subscribed to the newsletter. 🎉", "")
}
