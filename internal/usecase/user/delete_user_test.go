package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	"github.com/rober0xf/notifier/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser(t *testing.T) {
	t.Run("succesfully deleted user", func(t *testing.T) {
		s, mockRepo := setupDeleteUserTest(t)

		email := "richard@gmail.com"
		user := &entity.User{
			ID:       1,
			Username: "richard",
			Email:    email,
		}
		mockRepo.emails[email] = user
		mockRepo.users[1] = user

		err := s.Execute(context.Background(), 1, 1, user.Role)
		assert.NoError(t, err)

		_, err = mockRepo.GetUserByID(context.Background(), 1)
		assert.ErrorIs(t, err, repoErr.ErrNotFound)

		_, err = mockRepo.GetUserByEmail(context.Background(), email)
		assert.ErrorIs(t, err, repoErr.ErrNotFound)
	})

	t.Run("successfully deleted user as admin", func(t *testing.T) {
		s, mockRepo := setupDeleteUserTest(t)

		email := "richard@gmail.com"
		user := &entity.User{
			ID:       1,
			Username: "richard",
			Email:    email,
			Role:     entity.RoleAdmin,
		}
		mockRepo.emails[email] = user
		mockRepo.users[1] = user

		err := s.Execute(context.Background(), 1, 99, user.Role)
		assert.NoError(t, err)

		_, err = mockRepo.GetUserByID(context.Background(), 1)
		assert.ErrorIs(t, err, repoErr.ErrNotFound)

		_, err = mockRepo.GetUserByEmail(context.Background(), email)
		assert.ErrorIs(t, err, repoErr.ErrNotFound)
	})

	t.Run("returns forbidden when user tries to delete another user", func(t *testing.T) {
		s, mockRepo := setupDeleteUserTest(t)

		email := "richard@gmail.com"
		user := &entity.User{
			ID:       1,
			Username: "richard",
			Email:    email,
		}
		mockRepo.emails[email] = user
		mockRepo.users[1] = user

		err := s.Execute(context.Background(), 1, 2, user.Role)
		assert.ErrorIs(t, err, middleware.ErrForbidden)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		s, _ := setupDeleteUserTest(t)

		err := s.Execute(context.Background(), 99999, 99999, entity.RoleUser)
		assert.ErrorIs(t, err, domainErr.ErrUserNotFound)
	})

	t.Run("returns error for zero id", func(t *testing.T) {
		s, _ := setupDeleteUserTest(t)

		err := s.Execute(context.Background(), 0, 0, entity.RoleUser)
		assert.ErrorIs(t, err, domainErr.ErrInvalidUserData)
	})

	t.Run("returns error for negative id", func(t *testing.T) {
		s, _ := setupDeleteUserTest(t)

		err := s.Execute(context.Background(), -1, -1, entity.RoleUser)
		assert.ErrorIs(t, err, domainErr.ErrInvalidUserData)
	})

	t.Run("handles repository error", func(t *testing.T) {
		s, mockRepo := setupDeleteUserTest(t)

		mockRepo.err = errors.New("database connection failed")

		err := s.Execute(context.Background(), 1, 1, entity.RoleUser)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection failed")
	})
}
