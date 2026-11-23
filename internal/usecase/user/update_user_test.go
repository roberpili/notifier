package user_test

import (
	"context"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/usecase/user"
	"github.com/rober0xf/notifier/pkg/auth"
	"github.com/stretchr/testify/assert"

	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
)

func TestUpdateUserProfile(t *testing.T) {
	t.Run("successfully updates user profile", func(t *testing.T) {
		s, mockRepo := setupUpdateUserTest(t)

		// original data
		originalEmail := "original@gmail.com"
		hashedPassword, _ := auth.HashPassword("password123#!%")

		originalUser := &entity.User{
			ID:           1,
			Username:     "original username",
			Email:        originalEmail,
			PasswordHash: hashedPassword,
			IsActive:     true,
		}
		mockRepo.emails[originalEmail] = originalUser
		mockRepo.users[1] = originalUser

		// updated data
		input := user.UpdateUserInput{
			ID:       1,
			Username: strPtr("username changed"),
			Email:    strPtr("changed@gmail.com"),
		}
		updatedUser, err := s.Execute(context.Background(), input, input.ID, originalUser.Role)
		assert.NoError(t, err)
		assert.Equal(t, "username changed", updatedUser.Username)
		assert.Equal(t, "changed@gmail.com", updatedUser.Email)
		assert.Equal(t, hashedPassword, updatedUser.PasswordHash)

		getUserByEmail := user.NewGetUserByEmailService(mockRepo, nil)
		_, err = getUserByEmail.Execute(context.Background(), originalEmail)
		assert.Error(t, err)
	})

	t.Run("successfully updates only username", func(t *testing.T) {
		s, mockRepo := setupUpdateUserTest(t)

		email := "richard@gmail.com"
		usr := &entity.User{
			ID:           1,
			Username:     "old username",
			Email:        email,
			PasswordHash: "hashedpassword",
			IsActive:     true,
		}
		mockRepo.emails[email] = usr
		mockRepo.users[1] = usr

		input := user.UpdateUserInput{
			ID:       1,
			Username: strPtr("new username"),
		}

		updatedUser, err := s.Execute(context.Background(), input, input.ID, usr.Role)

		assert.NoError(t, err)
		assert.Equal(t, "new username", updatedUser.Username)
		assert.Equal(t, email, updatedUser.Email)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		s, _ := setupUpdateUserTest(t)

		input := user.UpdateUserInput{
			ID:       99999,
			Username: strPtr("richard"),
		}

		_, err := s.Execute(context.Background(), input, input.ID, entity.RoleUser)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrUserNotFound)
	})
}

func TestUpdateUserPassword(t *testing.T) {
	t.Run("successfully updates password", func(t *testing.T) {
		s, mockRepo := setupUpdateUserTest(t)

		oldHashedPassword, _ := auth.HashPassword("password123#!-")
		email := "richard@gmail.com"

		usr := &entity.User{
			ID:           1,
			Username:     "richard",
			Email:        email,
			PasswordHash: oldHashedPassword,
			IsActive:     true,
		}
		mockRepo.emails[email] = usr
		mockRepo.users[1] = usr

		plainPassword := "new#2Password123#!-"
		input := user.UpdateUserInput{
			ID:       1,
			Password: strPtr(plainPassword),
		}

		updatedUser, err := s.Execute(context.Background(), input, input.ID, usr.Role)

		assert.NoError(t, err)
		assert.NotEqual(t, oldHashedPassword, updatedUser.PasswordHash)
		assert.NotEqual(t, plainPassword, updatedUser.PasswordHash)

		isValid := auth.VerifyPassword(updatedUser.PasswordHash, plainPassword)
		assert.NoError(t, isValid)
	})
}

// pointer to string
func strPtr(s string) *string {
	return &s
}
