package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rober0xf/notifier/internal/auth"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
)

type GoogleOAuthService struct {
	userRepo repository.UserRepository
	tokenGen auth.TokenGenerator
	logger   *slog.Logger
}

func NewOAuthService(userRepo repository.UserRepository, tokenGen auth.TokenGenerator, logger *slog.Logger) *GoogleOAuthService {
	if logger == nil {
		logger = slog.Default()
	}

	return &GoogleOAuthService{
		userRepo: userRepo,
		tokenGen: tokenGen,
		logger:   logger,
	}
}

func (s *GoogleOAuthService) Execute(ctx context.Context, googleID, email, name string) (*LoginOutput, error) {
	if googleID == "" {
		return nil, domainErr.ErrInvalidGoogleID
	}

	if email == "" {
		return nil, domainErr.ErrInvalidEmailFormat
	}

	// search for google_id
	user, err := s.userRepo.GetUserByGoogleID(ctx, googleID)
	if err == nil {
		s.logger.DebugContext(ctx, "Google: existing account found",
			slog.Int("user_id", user.ID),
		)
		return s.buildLoginOutput(user)
	}

	if !errors.Is(err, repoErr.ErrNotFound) {
		s.logger.ErrorContext(ctx, "Google: failed to get user by Google ID",
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("get user by google id: %w", err)
	}

	// search for email
	user, err = s.userRepo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, repoErr.ErrNotFound) {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	// if exists then link account
	if user != nil {
		err = s.userRepo.UpdateUserGoogleID(ctx, user.ID, googleID)
		if err != nil {
			s.logger.ErrorContext(ctx, "Google: failed to link account",
				slog.Int("user_id", user.ID),
				slog.Any("error", err),
			)
			switch {
			case errors.Is(err, repoErr.ErrNotFound):
				return nil, domainErr.ErrUserNotFound
			case errors.Is(err, repoErr.ErrGoogleExists):
				return nil, auth.ErrGoogleAccountAlreadyLinked
			default:
				return nil, fmt.Errorf("link google account: %w", err)
			}
		}
		s.logger.InfoContext(ctx, "Google: account linked",
			slog.Int("user_id", user.ID),
		)

		user.GoogleID = googleID
		return s.buildLoginOutput(user)
	}

	// if it doesnt exists then create user
	user, err = s.userRepo.CreateOAuthUser(ctx, email, name, googleID)
	if err != nil {
		s.logger.ErrorContext(ctx, "Google: failed to create user",
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("create oauth user: %w", err)
	}
	s.logger.InfoContext(ctx, "Google: user created",
		slog.Int("user_id", user.ID),
	)

	return s.buildLoginOutput(user)
}

func (s *GoogleOAuthService) buildLoginOutput(user *entity.User) (*LoginOutput, error) {
	token, err := s.tokenGen.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &LoginOutput{Token: token, User: user}, nil
}
