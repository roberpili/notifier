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
	authErr "github.com/rober0xf/notifier/pkg/auth"
)

type DeleteUserService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

func NewDeleteUserService(userRepo repository.UserRepository, logger *slog.Logger) *DeleteUserService {
	if logger == nil {
		logger = slog.Default()
	}

	return &DeleteUserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *DeleteUserService) Execute(ctx context.Context, id, userID int, role entity.Role) error {
	if id <= 0 {
		return domainErr.ErrInvalidUserData
	}

	if id != userID && role != entity.RoleAdmin {
		return authErr.ErrForbidden
	}

	err := s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return domainErr.ErrUserNotFound
		}

		return fmt.Errorf("delete user: %w", err)
	}

	s.logger.InfoContext(ctx, "user deleted",
		slog.Int("target_user_id", id),
		slog.Int("actor_id", userID),
		slog.String("actor_role", string(role)),
	)

	return nil
}
