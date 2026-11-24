package email

import _ "embed"

//go:embed disposable_emails.txt
var disposableEmails []byte

func MustDisposableEmail() []string {
	emails := make([]string, 0)
	emails = append(emails, stringLines(disposableEmails)...)

	return emails
}
