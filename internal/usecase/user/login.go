package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rober0xf/notifier/internal/auth"
	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/domain/repository"
	"github.com/rober0xf/notifier/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	userRepo repository.UserRepository
	tokenGen auth.TokenGenerator
	logger   *slog.Logger
}

func NewLoginService(
	userRepo repository.UserRepository,
	tokenGen auth.TokenGenerator,
	logger *slog.Logger,
) *LoginService {
	if logger == nil {
		logger = slog.Default()
	}

	return &LoginService{
		userRepo: userRepo,
		tokenGen: tokenGen,
		logger:   logger,
	}
}

type LoginOutput struct {
	Token string
	User  *entity.User
}

func (s *LoginService) Execute(ctx context.Context, email, password string) (*LoginOutput, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, middleware.ErrInvalidCredentials
	}

	if err := auth.VerifyPassword(user.PasswordHash, password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, middleware.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("verify password: %w", err)
	}

	if !user.IsActive {
		return nil, auth.ErrEmailNotVerified
	}

	token, err := s.tokenGen.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}
	s.logger.InfoContext(ctx, "logged user",
		slog.Int("user_id", user.ID),
		slog.String("email", user.Email),
	)

	return &LoginOutput{Token: token, User: user}, nil
}
