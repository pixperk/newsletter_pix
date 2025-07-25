package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("brevo failed: %s", resp.Status)
	}

	return nil
}
