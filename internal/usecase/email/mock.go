package email

import (
	"context"
	"sync"
)

type SentEmail struct {
	To      []string
	Subject string
	Body    string
}

type MockSender struct {
	SentEmails []SentEmail
	Err        error
	wg         sync.WaitGroup
	expecting  bool
}

func NewMockSender() *MockSender {
	return &MockSender{
		SentEmails: []SentEmail{},
	}
}

func (m *MockSender) ExpectSends(n int) {
	m.expecting = true
	m.wg.Add(n)
}

func (m *MockSender) Wait() {
	m.wg.Wait()
}

func (m *MockSender) Send(ctx context.Context, to []string, subject, htmlBody string) error {
	if m.expecting {
		defer m.wg.Done()
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	m.SentEmails = append(m.SentEmails, SentEmail{
		To:      to,
		Subject: subject,
		Body:    htmlBody,
	})

	return m.Err
}
