package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
	t.Run("successfully found all users", func(t *testing.T) {
		s, mockRepo := setupGetAllUsersTest(t)

		users := []*entity.User{
			{ID: 1, Email: "user1@test.com", Username: "user1", IsActive: true},
			{ID: 2, Email: "user2@test.com", Username: "user2", IsActive: true},
			{ID: 3, Email: "user3@test.com", Username: "user3", IsActive: false},
		}

		for _, u := range users {
			mockRepo.emails[u.Email] = u
			mockRepo.users[u.ID] = u
		}

		result, err := s.Execute(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)

		userIDs := make(map[int]bool)
		for _, u := range result {
			userIDs[u.ID] = true
		}

		assert.True(t, userIDs[1])
		assert.True(t, userIDs[2])
		assert.True(t, userIDs[3])
	})

	t.Run("returns empty list when no users exist", func(t *testing.T) {
		s, _ := setupGetAllUsersTest(t)

		users, err := s.Execute(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Empty(t, users) // empty, not nil
	})

	t.Run("returns users with different status", func(t *testing.T) {
		s, mockRepo := setupGetAllUsersTest(t)

		activeUser := &entity.User{ID: 1, Email: "active@test.com", Username: "active", IsActive: true}
		inactiveUser := &entity.User{ID: 2, Email: "inactive@test.com", Username: "inactive", IsActive: false}

		mockRepo.emails["active@test.com"] = activeUser
		mockRepo.users[1] = activeUser
		mockRepo.emails["inactive@test.com"] = inactiveUser
		mockRepo.users[2] = inactiveUser

		users, err := s.Execute(context.Background())

		assert.NoError(t, err)
		assert.Len(t, users, 2)

		var hasActive, hasInactive bool
		for _, u := range users {
			if u.IsActive {
				hasActive = true
			} else {
				hasInactive = true
			}
		}
		assert.True(t, hasActive)
		assert.True(t, hasInactive)
	})

	t.Run("handles repository errors", func(t *testing.T) {
		s, mockRepo := setupGetAllUsersTest(t)

		mockRepo.err = errors.New("database connection failed")

		users, err := s.Execute(context.Background())

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Contains(t, err.Error(), "database connection failed")
	})
}
