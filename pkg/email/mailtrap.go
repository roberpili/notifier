package email

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mailtrap/mailtrap-go"
)

type MailtrapSender struct {
	client *mailtrap.Client
	from   mailtrap.Address
}

func NewMailtrapSender(token, fromEmail, fromName string) (*MailtrapSender, error) {
	c, err := mailtrap.NewClient(token)
	if err != nil {
		return nil, fmt.Errorf("mailtrap client: %w", err)
	}

	return &MailtrapSender{
		client: c,
		from: mailtrap.Address{
			Email: fromEmail,
			Name:  fromName,
		},
	}, nil
}

// with dialer we can log where it fails and the req doesnt get stuck
func (m *MailtrapSender) Send(ctx context.Context, to []string, subject string, htmlBody string) error {
	// skip in tests
	if os.Getenv("SKIP_EMAIL_SENDING") == "true" {
		log.Printf("Test mode: skipping email send to %v (subject: %s)", to, subject)
		return nil
	}

	if len(to) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	rcpts := make([]mailtrap.Address, 0, len(to))
	for _, addr := range to {
		rcpts = append(rcpts, mailtrap.Address{Email: addr})
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, _, err := m.client.Send(ctx, &mailtrap.SendRequest{
		From:    m.from,
		To:      rcpts,
		Subject: subject,
		HTML:    htmlBody,
	})
	if err != nil {
		return fmt.Errorf("mailtrap send failed: %w", err)
	}

	return nil
}
