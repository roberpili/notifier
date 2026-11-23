package email

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
