package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

const smtpSessionTimeout = 30 * time.Second

type SMTPSender struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSMTPSender(host, port, username, password, from string) *SMTPSender {
	if host == "" {
		host = "smtp.gmail.com"
	}
	if port == "" {
		port = "587"
	}
	if from == "" {
		from = username
	}
	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SMTPSender) Send(ctx context.Context, to []string, subject, htmlBody string) error {
	if len(to) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	from, err := mail.ParseAddress(s.from)
	if err != nil {
		return fmt.Errorf("invalid from address: %w", err)
	}
	rcpts := make([]string, 0, len(to))
	for _, addr := range to {
		a, err := mail.ParseAddress(addr)
		if err != nil {
			return fmt.Errorf("invalid recipient %q: %w", addr, err)
		}
		rcpts = append(rcpts, a.Address)
	}
	if strings.ContainsAny(subject, "\r\n") {
		return fmt.Errorf("subject contains line breaks")
	}

	address := net.JoinHostPort(s.host, s.port)
	dialer := &net.Dialer{Timeout: 5 * time.Second}

	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}

	// bound the whole session
	deadline := time.Now().Add(smtpSessionTimeout)
	if dl, ok := ctx.Deadline(); ok && dl.Before(deadline) {
		deadline = dl
	}
	if err := conn.SetDeadline(deadline); err != nil {
		conn.Close()
		return fmt.Errorf("set deadline failed: %w", err)
	}

	// abort immediately if context is cancelled
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			conn.Close()
		case <-done:
		}
	}()

	c, err := smtp.NewClient(conn, s.host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("smtp client failed: %w", err)
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		cfg := &tls.Config{ServerName: s.host, MinVersion: tls.VersionTLS12}
		if err = c.StartTLS(cfg); err != nil {
			return fmt.Errorf("starttls failed: %w", err)
		}
	} else if s.username != "" {
		return fmt.Errorf("server %s does not support STARTTLS; refusing to send credentials", address)
	}

	if s.username != "" {
		auth := smtp.PlainAuth("", s.username, s.password, s.host)
		if err = c.Auth(auth); err != nil {
			return fmt.Errorf("auth failed: %w", err)
		}
	}

	if err = c.Mail(from.Address); err != nil {
		return fmt.Errorf("mail from failed: %w", err)
	}
	for _, addr := range rcpts {
		if err = c.Rcpt(addr); err != nil {
			return fmt.Errorf("rcpt %s failed: %w", addr, err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("data failed: %w", err)
	}

	msg := "From: " + from.String() + "\r\n" +
		"To: " + strings.Join(rcpts, ", ") + "\r\n" +
		"Subject: " + mime.QEncoding.Encode("utf-8", subject) + "\r\n" +
		"Date: " + time.Now().Format(time.RFC1123Z) + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		htmlBody

	if _, err = w.Write([]byte(msg)); err != nil {
		w.Close()
		return fmt.Errorf("write body failed: %w", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("finish body failed: %w", err)
	}

	return c.Quit()
}
