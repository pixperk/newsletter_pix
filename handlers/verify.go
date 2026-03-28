package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
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
	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head>
<body style="margin:0;padding:0;background-color:#f4f4f7;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f7;padding:40px 0;">
    <tr><td align="center">
      <table role="presentation" width="560" cellpadding="0" cellspacing="0" style="background-color:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 4px 24px rgba(0,0,0,0.06);">

        <!-- Header -->
        <tr>
          <td style="background:linear-gradient(135deg,#0066cc 0%%,#004999 100%%);padding:40px 48px;text-align:center;">
            <img src="https://www.pixperk.tech/assets/avatar.jpg" alt="PixPerk" width="64" height="64" style="border-radius:50%%;border:3px solid rgba(255,255,255,0.3);margin-bottom:16px;" />
            <h1 style="margin:0;color:#ffffff;font-size:22px;font-weight:600;letter-spacing:-0.3px;">Hey, it's Yashaswi (aka PixPerk)</h1>
            <p style="margin:8px 0 0;color:rgba(255,255,255,0.8);font-size:14px;">One last step to join the newsletter</p>
          </td>
        </tr>

        <!-- Body -->
        <tr>
          <td style="padding:40px 48px;">
            <p style="margin:0 0 20px;color:#333;font-size:16px;line-height:1.6;">Hi there! 👋</p>
            <p style="margin:0 0 28px;color:#555;font-size:15px;line-height:1.7;">
              Thanks for subscribing to the <strong style="color:#333;">PixPerk</strong> newsletter. Click the button below to verify your email and start receiving updates on backend engineering, dev insights, and more.
            </p>

            <!-- CTA Button -->
            <table role="presentation" cellpadding="0" cellspacing="0" style="margin:0 auto 28px;">
              <tr>
                <td style="border-radius:8px;background-color:#0066cc;">
                  <a href="%s" target="_blank" style="display:inline-block;padding:14px 36px;color:#ffffff;font-size:15px;font-weight:600;text-decoration:none;letter-spacing:0.3px;">
                    Verify Email Address
                  </a>
                </td>
              </tr>
            </table>

            <p style="margin:0 0 8px;color:#888;font-size:13px;line-height:1.6;text-align:center;">This link expires in 24 hours.</p>
            <p style="margin:0;color:#888;font-size:13px;line-height:1.6;text-align:center;">If you didn't request this, you can safely ignore this email.</p>
          </td>
        </tr>

        <!-- Divider -->
        <tr><td style="padding:0 48px;"><hr style="border:none;border-top:1px solid #eee;margin:0;" /></td></tr>

        <!-- Footer -->
        <tr>
          <td style="padding:28px 48px 36px;text-align:center;">
            <p style="margin:0 0 6px;color:#333;font-size:14px;font-weight:600;">Yashaswi</p>
            <p style="margin:0 0 16px;color:#888;font-size:13px;">Backend Developer &middot; <a href="https://www.pixperk.tech" target="_blank" style="color:#0066cc;text-decoration:none;">pixperk.tech</a></p>
            <p style="margin:0;color:#bbb;font-size:11px;">&copy; 2025 Yashaswi &mdash; All bytes reserved.</p>
          </td>
        </tr>

      </table>
    </td></tr>
  </table>
</body>
</html>`, verificationURL)

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
