package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	t.Run("succesfully found the user by id", func(t *testing.T) {
		s, mockRepo := setupGetUserByIDTest(t)

		expectedUser := &entity.User{
			ID:           1,
			Email:        "richard@gmail.com",
			Username:     "piper",
			PasswordHash: "password123#!-",
			IsActive:     true,
		}
		mockRepo.emails["richard@gmail.com"] = expectedUser
		mockRepo.users[1] = expectedUser

		user, err := s.Execute(context.Background(), expectedUser.ID)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Email, user.Email)
		assert.Equal(t, "piper", user.Username)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		s, _ := setupGetUserByIDTest(t)

		nonExistingID := 99999
		user, err := s.Execute(context.Background(), nonExistingID)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, domainErr.ErrUserNotFound)
	})

	t.Run("returns error for zero id", func(t *testing.T) {
		s, _ := setupGetUserByIDTest(t)

		user, err := s.Execute(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, domainErr.ErrInvalidUserData)
	})

	t.Run("returns error for negative id", func(t *testing.T) {
		s, _ := setupGetUserByIDTest(t)

		user, err := s.Execute(context.Background(), -1)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, domainErr.ErrInvalidUserData)
	})

	t.Run("handles repository error", func(t *testing.T) {
		s, mockRepo := setupGetUserByIDTest(t)

		// repo error
		mockRepo.err = errors.New("database connection failed")
		user, err := s.Execute(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "database connection failed")
	})
}
