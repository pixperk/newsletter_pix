package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	year := strconv.Itoa(time.Now().Year())
	emailSubject := "Verify your newsletter subscription"
	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1.0">
  <style media="all" type="text/css">
    @media only screen and (max-width: 640px) {
      .container { padding: 8px !important; width: 100%% !important; }
      .card { border-radius: 12px !important; }
      .card-body { padding: 28px 24px !important; }
      .card-footer { padding: 20px 24px 28px !important; }
    }
  </style>
</head>
<body style="margin:0;padding:0;background-color:#0a0a0a;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Inter,Roboto,Helvetica,Arial,sans-serif;-webkit-font-smoothing:antialiased;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#0a0a0a;padding:40px 0;">
    <tr><td align="center">

      <!-- Header -->
      <table role="presentation" width="560" cellpadding="0" cellspacing="0" class="container" style="margin-bottom:20px;">
        <tr>
          <td style="padding:0 8px;text-align:left;" align="left">
            <span style="font-size:15px;font-weight:600;color:#e0e0e0;letter-spacing:-0.3px;">pixperk</span>
            <span style="font-size:15px;font-weight:300;color:#333;">&middot;</span>
            <span style="font-size:13px;font-weight:400;color:#484848;">verify</span>
          </td>
        </tr>
      </table>

      <!-- Card -->
      <table role="presentation" width="560" cellpadding="0" cellspacing="0" class="container card" style="background-color:#111111;border:1px solid rgba(255,255,255,0.06);border-radius:14px;overflow:hidden;">

        <!-- Avatar + Greeting -->
        <tr>
          <td class="card-body" style="padding:40px 40px 0;">
            <table role="presentation" cellpadding="0" cellspacing="0" style="margin-bottom:28px;">
              <tr>
                <td style="vertical-align:middle;width:48px;" width="48" valign="middle">
                  <img src="https://www.pixperk.tech/assets/avatar.jpg" alt="PixPerk" width="44" height="44" style="border-radius:50%%;border:1px solid rgba(255,255,255,0.08);display:block;" />
                </td>
                <td style="vertical-align:middle;padding-left:14px;" valign="middle">
                  <span style="font-size:16px;font-weight:600;color:#ececec;letter-spacing:-0.3px;">hey, it's yashaswi</span><br/>
                  <span style="font-size:13px;color:#555;">one last step to join the newsletter</span>
                </td>
              </tr>
            </table>
          </td>
        </tr>

        <!-- Body -->
        <tr>
          <td class="card-body" style="padding:0 40px 36px;">
            <p style="margin:0 0 20px;color:#9a9a9a;font-size:15px;line-height:1.75;">
              thanks for subscribing to the <strong style="color:#d4d4d4;">pixperk</strong> newsletter. click below to verify your email and start receiving updates on backend engineering, dev insights, and more.
            </p>

            <!-- CTA Button -->
            <table role="presentation" cellpadding="0" cellspacing="0" style="margin:0 auto 28px;">
              <tr>
                <td style="border-radius:8px;background-color:#e0e0e0;">
                  <a href="%s" target="_blank" style="display:inline-block;padding:12px 32px;color:#0a0a0a;font-size:14px;font-weight:600;text-decoration:none;letter-spacing:0.2px;">
                    Verify Email Address
                  </a>
                </td>
              </tr>
            </table>

            <p style="margin:0 0 6px;color:#555;font-size:12px;line-height:1.6;text-align:center;">This link expires in 24 hours.</p>
            <p style="margin:0;color:#444;font-size:12px;line-height:1.6;text-align:center;">If you didn't request this, you can safely ignore this email.</p>
          </td>
        </tr>

        <!-- Divider -->
        <tr><td style="padding:0 40px;"><hr style="border:none;border-top:1px solid rgba(255,255,255,0.06);margin:0;" /></td></tr>

        <!-- Footer -->
        <tr>
          <td class="card-footer" style="padding:24px 40px 32px;text-align:center;">
            <p style="margin:0 0 4px;color:#d4d4d4;font-size:13px;font-weight:600;">Yashaswi</p>
            <p style="margin:0 0 14px;color:#555;font-size:12px;">Backend Developer &middot; <a href="https://www.pixperk.tech" target="_blank" style="color:#999;text-decoration:none;">pixperk.tech</a></p>
            <p style="margin:0;color:#333;font-size:11px;">&copy; %s Yashaswi &mdash; All bytes reserved.</p>
          </td>
        </tr>

      </table>
    </td></tr>
  </table>
</body>
</html>`, verificationURL, year)

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
