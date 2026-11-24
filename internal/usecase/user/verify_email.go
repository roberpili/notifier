package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	"github.com/rober0xf/notifier/internal/middleware"
)

type VerifyEmailService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

func NewVerifyEmailService(userRepo repository.UserRepository, logger *slog.Logger) *VerifyEmailService {
	if logger == nil {
		logger = slog.Default()
	}

	return &VerifyEmailService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *VerifyEmailService) Execute(ctx context.Context, plainToken string) (*entity.User, error) {
	hash := sha256.Sum256([]byte(plainToken))
	tokenHash := hex.EncodeToString(hash[:])

	// verify and mark token as used
	token, err := s.userRepo.VerifyAndConsumeToken(ctx, tokenHash, entity.TokenPurposeEmailVerification)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, middleware.ErrInvalidToken
		}

		return nil, fmt.Errorf("verify and consume token: %w", err)
	}

	user, err := s.userRepo.UpdateUserIsActiveReturning(ctx, token.UserID, true)
	if err != nil {
		return nil, fmt.Errorf("update is_active: %w", err)
	}

	s.logger.InfoContext(ctx, "email verified",
		slog.Int("user_id", user.ID),
		slog.String("email", user.Email))

	return user, nil
}
