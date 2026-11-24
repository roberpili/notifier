package email

import _ "embed"

//go:embed disposable_emails.txt
var disposableEmails []byte

func MustDisposableEmail() []string {
	emails := make([]string, 0)
	for _, domain := range stringLines(disposableEmails) {
		emails = append(emails, domain)
	}
	return emails
}
