package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
)

type GetUserByEmailService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

func NewGetUserByEmailService(userRepo repository.UserRepository, logger *slog.Logger) *GetUserByEmailService {
	if logger == nil {
		logger = slog.Default()
	}

	return &GetUserByEmailService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *GetUserByEmailService) Execute(ctx context.Context, email string) (*entity.User, error) {
	if err := ValidateEmail(email, nil); err != nil {
		return nil, domainErr.ErrInvalidEmailFormat
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, domainErr.ErrUserNotFound
		}

		return nil, fmt.Errorf("get user by email: %w", err)
	}

	s.logger.InfoContext(ctx, "user got by email",
		slog.String("email", email),
		slog.Int("user_id", user.ID),
	)

	return user, nil
}
