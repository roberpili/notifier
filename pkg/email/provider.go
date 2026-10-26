package email

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/resend/resend-go/v4"
)

type SenderProvider struct {
	client *resend.Client
	from   string
}

func NewSenderProvider(token, fromEmail, fromName string) *SenderProvider {
	from := fromEmail

	if fromName != "" {
		from = fmt.Sprintf("%s <%s>", fromName, fromEmail)
	}

	return &SenderProvider{
		client: resend.NewClient(token),
		from:   from,
	}
}

func (s *SenderProvider) Send(ctx context.Context, to []string, subject string, htmlBody string) error {
	// skip in tests
	if os.Getenv("SKIP_EMAIL_SENDING") == "true" {
		log.Printf("Test mode: skipping email send to %v (subject: %s)", to, subject)
		return nil
	}

	if len(to) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      to,
		Subject: subject,
		Html:    htmlBody,
	}

	result, err := s.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		log.Printf("EMAIL FAILED: to=%v subject=%q error=%v", to, subject, err)
		return fmt.Errorf("resend send failed: %w", err)
	}

	log.Printf(
		"EMAIL SENT: to=%v subject=%q email_id=%s",
		to,
		subject,
		result.Id,
	)

	return nil
}
