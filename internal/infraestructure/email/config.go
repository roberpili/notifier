package email

import (
	"github.com/rober0xf/notifier/internal/database"
)

type Config struct {
	Provider string
	From     string
	FromName string
	APIToken string

	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string

	SkipSending bool
}

func NewConfigFromEnv() *Config {
	return &Config{
		Provider:     database.GetEnvOrFatal[string]("EMAIL_PROVIDER"),
		From:         database.GetEnvOrFatal[string]("MAIL_FROM"),
		FromName:     database.GetEnvOrFatal[string]("MAIL_FROM_NAME"),
		APIToken:     database.GetEnvOrFatal[string]("EMAIL_API_TOKEN"),
		SMTPHost:     database.GetEnvOrFatal[string]("SMTP_HOST"),
		SMTPPort:     database.GetEnvOrFatal[int]("SMTP_PORT"),
		SMTPUsername: database.GetEnvOrFatal[string]("SMTP_USERNAME"),
		SMTPPassword: database.GetEnvOrFatal[string]("SMTP_PASSWORD"),
		SkipSending:  database.GetEnvOrFatal[bool]("SKIP_EMAIL_SENDING"),
	}
}
