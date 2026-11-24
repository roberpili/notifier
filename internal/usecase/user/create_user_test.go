package user_test

import (
	"context"
	"errors"
	"testing"

	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	t.Run("succesfully creates a new user", func(t *testing.T) {
		s, _, mockEmailSender := setupCreateUserTest(t)

		user, err := s.Execute(context.Background(), "richard", "piedpiper@gmail.com", "password123#-!")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, 0, len(mockEmailSender.SentEmails)) // 0 because env=test
		assert.NotEmpty(t, user.ID)
	})

	t.Run("fails when username in use", func(t *testing.T) {
		s, _, _ := setupCreateUserTest(t)

		_, err := s.Execute(context.Background(), "richard", "piedpiper@gmail.com", "password123#-!")
		assert.NoError(t, err)

		_, err = s.Execute(context.Background(), "richard", "pied@gmail.com", "password123#-!")
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrUsernameAlreadyExists)
	})

	t.Run("fails when email in use", func(t *testing.T) {
		s, _, _ := setupCreateUserTest(t)

		_, err := s.Execute(context.Background(), "richard", "piedpiper@gmail.com", "password123#-!")
		assert.NoError(t, err)

		_, err = s.Execute(context.Background(), "richard2", "piedpiper@gmail.com", "password123#-!")
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrEmailAlreadyExists)
	})

	t.Run("fails with invalid email format", func(t *testing.T) {
		s, _, _ := setupCreateUserTest(t)

		_, err := s.Execute(context.Background(), "richard", "notanemail", "password123#-!")
		assert.ErrorIs(t, err, domainErr.ErrInvalidEmailFormat)
	})

	t.Run("fails when email domain cannot get emails", func(t *testing.T) {
		s, _, _ := setupCreateUserTest(t)

		_, err := s.Execute(context.Background(), "richard", "richard@mailnator.com", "password123#-!")
		assert.ErrorIs(t, err, domainErr.ErrInvalidDomain)
	})

	t.Run("fails with weak password", func(t *testing.T) {
		s, _, _ := setupCreateUserTest(t)

		_, err := s.Execute(context.Background(), "richard", "piedpiper@gmail.com", "123")
		assert.ErrorIs(t, err, domainErr.ErrInvalidPassword)
	})

	t.Run("stores hashed password not plaintext", func(t *testing.T) {
		s, mockRepo, _ := setupCreateUserTest(t)

		_, err := s.Execute(context.Background(), "richard", "piedpiper@gmail.com", "password123#-!")
		assert.NoError(t, err)

		storedUser := mockRepo.users[1]
		assert.NotEqual(t, "password123#-!", storedUser.PasswordHash)
		assert.NotEmpty(t, storedUser.PasswordHash)
	})

	t.Run("sends email when not in test mode", func(t *testing.T) {
		t.Setenv("ENV", "development")
		s, _, mockEmailSender := setupCreateUserTest(t)

		mockEmailSender.ExpectSends(1)
		_, err := s.Execute(context.Background(), "john", "john@example.com", "password123#-!")
		assert.NoError(t, err)

		mockEmailSender.Wait()
		assert.Equal(t, 1, len(mockEmailSender.SentEmails))

		sentEmail := mockEmailSender.SentEmails[0]
		assert.Equal(t, []string{"john@example.com"}, sentEmail.To)
		assert.Contains(t, sentEmail.Body, "verify account")
	})

	t.Run("creates user even when email sending fails", func(t *testing.T) {
		t.Setenv("ENV", "development")
		s, _, mockEmailSender := setupCreateUserTest(t)

		// simulate error
		mockEmailSender.Err = errors.New("SMTP connection failed")
		mockEmailSender.ExpectSends(1)

		user, err := s.Execute(context.Background(), "john", "john@example.com", "password123#-!")
		mockEmailSender.Wait()

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, 1, len(mockEmailSender.SentEmails))
	})

	t.Run("handles repository error", func(t *testing.T) {
		s, mockRepo, _ := setupCreateUserTest(t)

		mockRepo.err = errors.New("database connection failed")
		_, err := s.Execute(context.Background(), "richard", "piedpiper@gmail.com", "password123#-!")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection failed")
	})
}
