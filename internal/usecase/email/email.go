package email

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/rober0xf/notifier/internal/infraestructure/email"
)

type EmailService struct {
	sender email.EmailSender
	logger *slog.Logger
}

func NewEmailService(sender email.EmailSender, logger *slog.Logger) *EmailService {
	return &EmailService{
		sender: sender,
		logger: logger,
	}
}

func (s *EmailService) SendWithRetry(ctx context.Context, to []string, subject, body string) error {
	var err error

	for i := range 5 {
		attemptCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err = s.sender.Send(attemptCtx, to, subject, body)
		cancel()

		if err == nil {
			return nil
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		s.logger.WarnContext(ctx, "failed to send email",
			"attempt", i+1,
			"error", err)

		if i < 4 {
			select {
			case <-time.After(time.Duration(1<<i) * time.Second):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return fmt.Errorf("send email after retries: %w", err)
}
