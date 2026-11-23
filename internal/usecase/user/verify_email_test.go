package user_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func TestVerifyEmail(t *testing.T) {
	t.Run("successfully activates when a verification is valid", func(t *testing.T) {
		s, mockRepo := setupVerifyEmailTest(t)

		plainToken := "verification_123"
		tokenHash := sha256.Sum256([]byte(plainToken))
		tokenHashString := hex.EncodeToString(tokenHash[:])

		email := "richard@gmail.com"
		user := &entity.User{
			ID:           1,
			Username:     "richard",
			Email:        email,
			PasswordHash: "hashedPassword",
			IsActive:     false,
		}
		mockRepo.tokens[tokenHashString] = &entity.UserToken{
			ID:        1,
			UserID:    1,
			TokenHash: tokenHashString,
			Purpose:   entity.TokenPurposeEmailVerification,
			Used:      false,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		mockRepo.users[1] = user

		verifiedUser, err := s.Execute(context.Background(), plainToken)
		assert.NoError(t, err)
		assert.True(t, verifiedUser.IsActive)

		storedUser, err := mockRepo.GetUserByID(context.Background(), 1)
		assert.True(t, storedUser.IsActive)
	})

	t.Run("returns error when token is invalid", func(t *testing.T) {
		s, _ := setupVerifyEmailTest(t)

		_, err := s.Execute(context.Background(), "wrong_token")
		assert.ErrorIs(t, err, auth.ErrInvalidToken)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		s, mockRepo := setupVerifyEmailTest(t)
		plainToken := "plain_token"
		tokenHash := sha256.Sum256([]byte(plainToken))
		tokenHashString := hex.EncodeToString(tokenHash[:])

		mockRepo.tokens[tokenHashString] = &entity.UserToken{
			ID:        1,
			UserID:    1,
			TokenHash: tokenHashString,
			Purpose:   entity.TokenPurposeEmailVerification,
			Used:      false,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}

		_, err := s.Execute(context.Background(), plainToken)
		assert.Contains(t, err.Error(), "resource not found")
	})
}
