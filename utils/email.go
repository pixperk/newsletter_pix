package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Sender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Recipient struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type BrevoEmail struct {
	Sender      Sender      `json:"sender"`
	To          []Recipient `json:"to"`
	Subject     string      `json:"subject"`
	HtmlContent string      `json:"htmlContent"`
}

func SendEmail(to, subject, htmlBody string) error {
	email := BrevoEmail{
		Sender: Sender{
			Name:  "Yashaswi",
			Email: os.Getenv("BREVO_SENDER_EMAIL"),
		},
		To: []Recipient{
			{Email: to, Name: to},
		},
		Subject:     subject,
		HtmlContent: htmlBody,
	}

	payload, _ := json.Marshal(email)

	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("api-key", os.Getenv("BREVO_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to Brevo: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for detailed error information
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read Brevo response: %w", err)
	}

	if resp.StatusCode >= 400 {
		// Try to parse Brevo error response
		var brevoError struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		}

		if json.Unmarshal(body, &brevoError) == nil && brevoError.Message != "" {
			return fmt.Errorf("brevo API error (%s): %s", resp.Status, brevoError.Message)
		}

		// Fallback to status and raw body if parsing fails
		return fmt.Errorf("brevo API error (%s): %s", resp.Status, string(body))
	}

	return nil
}
