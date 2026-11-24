package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/domain/repository"
)

type GetAllUsersService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

func NewGetAllUsersService(userRepo repository.UserRepository, logger *slog.Logger) *GetAllUsersService {
	if logger == nil {
		logger = slog.Default()
	}

	return &GetAllUsersService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *GetAllUsersService) Execute(ctx context.Context) ([]entity.User, error) {
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	s.logger.InfoContext(ctx, "fetched all users",
		slog.Int("count", len(users)),
	)

	return users, nil
}
