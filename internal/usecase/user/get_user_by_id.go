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

type GetUserByIDService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

func NewGetUserByIDService(userRepo repository.UserRepository, logger *slog.Logger) *GetUserByIDService {
	if logger == nil {
		logger = slog.Default()
	}

	return &GetUserByIDService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *GetUserByIDService) Execute(ctx context.Context, id int) (*entity.User, error) {
	if id <= 0 {
		return nil, domainErr.ErrInvalidUserData
	}

	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, domainErr.ErrUserNotFound
		}

		return nil, fmt.Errorf("get user by id: %w", err)
	}

	slog.InfoContext(ctx, "user got by id",
		slog.Int("user_id", user.ID),
	)

	return user, nil
}
