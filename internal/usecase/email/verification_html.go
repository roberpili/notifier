package email

import (
	"fmt"
)

func VerificationEmailHTML(token string, baseURL string) string {
	verificationURL := fmt.Sprintf("%s/v1/users/email_verification/%s", baseURL, token)

	return fmt.Sprintf(`
		<html>
			<body>
				<h1>Email verification</h1>
				<p>click the link to verify:</p>
				<a href="%s">verify account</a>
			</body>
		</html>
	`, verificationURL)
}
