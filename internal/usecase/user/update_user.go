package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rober0xf/notifier/internal/auth"
	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	"github.com/rober0xf/notifier/internal/middleware"
)

type UpdateUserService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

func NewUpdateUserService(userRepo repository.UserRepository, logger *slog.Logger) *UpdateUserService {
	if logger == nil {
		logger = slog.Default()
	}

	return &UpdateUserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

type UpdateUserInput struct {
	ID       int
	Username *string
	Email    *string
	Password *string
}

func (s *UpdateUserService) Execute(ctx context.Context, input UpdateUserInput, userID int, role entity.Role) (*entity.User, error) {
	existingUser, err := s.fetchUser(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	if err := s.authorizeUpdate(userID, existingUser.ID, role); err != nil {
		return nil, err
	}

	profileChanged, passwordChanged, err := s.applyChanges(existingUser, input)
	if err != nil {
		return nil, err
	}

	if !profileChanged && !passwordChanged {
		return existingUser, nil
	}

	if err := s.persistChanges(ctx, existingUser, profileChanged, passwordChanged); err != nil {
		return nil, err
	}

	s.logger.InfoContext(ctx, "user updated",
		slog.Int("target_user_id", existingUser.ID),
		slog.Int("actor_id", userID),
		slog.String("actor_role", string(role)),
		slog.Bool("profile_changed", profileChanged),
		slog.Bool("password_changed", passwordChanged),
	)

	return existingUser, nil
}

func (s *UpdateUserService) fetchUser(ctx context.Context, id int) (*entity.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repoErr.ErrNotFound):
			return nil, domainErr.ErrUserNotFound
		default:
			return nil, fmt.Errorf("get user by id: %w", err)
		}
	}

	return user, nil
}

func (s *UpdateUserService) authorizeUpdate(givenID, expectedID int, role entity.Role) error {
	if givenID != expectedID && role != entity.RoleAdmin {
		return middleware.ErrForbidden
	}

	return nil
}

func (s *UpdateUserService) applyChanges(user *entity.User, input UpdateUserInput) (profileChanged, passwordChanged bool, err error) {
	if input.Username != nil && *input.Username != "" && *input.Username != user.Username {
		user.Username = *input.Username
		profileChanged = true
	}

	if input.Email != nil && *input.Email != "" && *input.Email != user.Email {
		if err := ValidateEmail(*input.Email, nil); err != nil {
			return false, false, domainErr.ErrInvalidEmailFormat
		}

		user.Email = *input.Email
		profileChanged = true
	}

	if input.Password != nil && *input.Password != "" {
		if err := ValidatePassword(*input.Password); err != nil {
			return false, false, domainErr.ErrInvalidPassword
		}

		hashedPassword, err := auth.HashPassword(*input.Password)
		if err != nil {
			return false, false, fmt.Errorf("hash password: %w", err)
		}

		user.PasswordHash = hashedPassword
		passwordChanged = true
	}

	return profileChanged, passwordChanged, nil
}

func (s *UpdateUserService) persistChanges(ctx context.Context, user *entity.User, profileChanged, passwordChanged bool) error {
	if profileChanged {
		if err := s.userRepo.UpdateUserProfile(ctx, user.ID, user.Username, user.Email); err != nil {
			switch {
			case errors.Is(err, repoErr.ErrNotFound):
				return domainErr.ErrUserNotFound
			case errors.Is(err, repoErr.ErrEmailAlreadyExists):
				return domainErr.ErrEmailAlreadyExists
			case errors.Is(err, repoErr.ErrUsernameAlreadyExists):
				return domainErr.ErrUsernameAlreadyExists
			default:
				return fmt.Errorf("update user profile: %w", err)
			}
		}
	}

	if passwordChanged {
		if err := s.userRepo.UpdateUserPassword(ctx, user.ID, user.PasswordHash); err != nil {
			switch {
			case errors.Is(err, repoErr.ErrNotFound):
				return domainErr.ErrUserNotFound
			default:
				return fmt.Errorf("update user password: %w", err)
			}
		}
	}

	return nil
}
