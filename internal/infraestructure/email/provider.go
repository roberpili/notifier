package email

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/resend/resend-go/v4"
)

type SenderProvider struct {
	client *resend.Client
	from   string
	logger *slog.Logger
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
		s.logger.Error("failed to send email", "error", err, "to", to, "subject", subject)
		return fmt.Errorf("resend send failed: %w", err)
	}
	s.logger.Info("email sent successfully", "to", to, "subject", subject, "email_id", result.Id)

	return nil
}
