package email

import (
	"context"
	"fmt"
)

type EmailSender interface {
	Send(ctx context.Context, to []string, subject, htmlBody string) error
}

func NewSender(cfg Config) (EmailSender, error) {
	if cfg.SkipSending {
		return NoopSender{}, nil
	}

	if cfg.From == "" {
		return nil, fmt.Errorf("email sender address is required")
	}

	switch cfg.Provider {
	case "resend":
		if cfg.APIToken == "" {
			return nil, fmt.Errorf("email API token is required")
		}

		return NewSenderProvider(cfg.APIToken, cfg.From, cfg.FromName), nil

	case "smtp", "":
		if cfg.SMTPHost == "" || cfg.SMTPPort == 0 {
			return nil, fmt.Errorf("SMTP host and port is required")
		}

		return NewSMTPSender(
			cfg.SMTPHost,
			cfg.SMTPPort,
			cfg.SMTPUsername,
			cfg.SMTPPassword,
			cfg.From,
			cfg.FromName,
		), nil

	default:
		return nil, fmt.Errorf("unknown email provider %q", cfg.Provider)
	}
}

// logs
type NoopSender struct{}

func (NoopSender) Send(_ context.Context,
	_ []string,
	_ string,
	_ string) error {
	return nil
}
