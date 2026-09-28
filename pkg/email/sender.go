package email

import (
	"context"
	"fmt"
	"log"
	"os"
)

type EmailSender interface {
	Send(ctx context.Context, to []string, subject, htmlBody string) error
}

// logs
type NoopSender struct{}

func (NoopSender) Send(_ context.Context, to []string, subject, _ string) error {
	log.Printf("email skipped: to=%v subject=%q", to, subject)
	return nil
}

func NewFromEnv() (EmailSender, error) {
	if os.Getenv("SKIP_EMAIL_SENDING") == "true" {
		return NoopSender{}, nil
	}

	from := os.Getenv("MAIL_FROM")
	if from == "" {
		return nil, fmt.Errorf("MAIL_FROM is required")
	}

	switch os.Getenv("EMAIL_PROVIDER") {
	case "mailtrap":
		token := os.Getenv("MAILTRAP_API_TOKEN")
		if token == "" {
			return nil, fmt.Errorf("MAILTRAP_API_TOKEN is required")
		}
		return NewMailtrapSender(token, from, os.Getenv("MAIL_FROM_NAME"))

	case "smtp", "":
		return NewSMTPSender(
			os.Getenv("SMTP_HOST"),
			os.Getenv("SMTP_PORT"),
			os.Getenv("SMTP_USERNAME"),
			os.Getenv("SMTP_PASSWORD"),
			from,
		), nil

	default:
		return nil, fmt.Errorf("unknown EMAIL_PROVIDER %q", os.Getenv("EMAIL_PROVIDER"))
	}
}
