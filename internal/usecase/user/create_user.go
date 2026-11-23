package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	infraemail "github.com/rober0xf/notifier/internal/infraestructure/email"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	emailusecase "github.com/rober0xf/notifier/internal/usecase/email"
	"github.com/rober0xf/notifier/pkg/auth"
	"github.com/rober0xf/notifier/pkg/token"
)

type CreateUserService struct {
	userRepo          repository.UserRepository
	emailService      *emailusecase.EmailService
	disposableDomains []string
	baseURL           string
	logger            *slog.Logger
}

func NewCreateUserService(userRepo repository.UserRepository,
	emailSender infraemail.EmailSender,
	disposableDomains []string,
	baseURL string,
	logger *slog.Logger) *CreateUserService {
	if logger == nil {
		logger = slog.Default()
	}
	return &CreateUserService{
		userRepo:          userRepo,
		emailService:      emailusecase.NewEmailService(emailSender, logger),
		disposableDomains: disposableDomains,
		baseURL:           baseURL,
		logger:            logger,
	}
}

func (s *CreateUserService) Execute(ctx context.Context, username string, email string, password string) (*entity.User, error) {
	if err := s.validateInput(email, password); err != nil {
		return nil, err
	}

	user, err := s.buildUser(username, email, password)
	if err != nil {
		return nil, err
	}

	createdUser, err := s.persistUser(ctx, user)
	if err != nil {
		return nil, err
	}

	go s.dispatchVerificationEmail(createdUser)

	return createdUser, nil
}

func (s *CreateUserService) validateInput(email, password string) error {
	if err := ValidateEmail(email, s.disposableDomains); err != nil {
		return err
	}

	if err := ValidatePassword(password); err != nil {
		return domainErr.ErrInvalidPassword
	}

	return nil
}

func (s *CreateUserService) buildUser(username, email, password string) (*entity.User, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	return &entity.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}, nil
}

func (s *CreateUserService) persistUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, repoErr.ErrEmailAlreadyExists):
			return nil, domainErr.ErrEmailAlreadyExists
		case errors.Is(err, repoErr.ErrUsernameAlreadyExists):
			return nil, domainErr.ErrUsernameAlreadyExists
		default:
			return nil, fmt.Errorf("create user: %w", err)
		}
	}

	return createdUser, nil
}

func (s *CreateUserService) sendVerificationEmail(ctx context.Context, user *entity.User) error {
	tokenData, err := token.GenerateVerificationToken(12)
	if err != nil {
		return fmt.Errorf("generate verification token: %w", err)
	}

	_, err = s.userRepo.CreateUserToken(ctx, &entity.UserToken{
		UserID:    user.ID,
		TokenHash: tokenData.Hash,
		Purpose:   entity.TokenPurposeEmailVerification,
		ExpiresAt: tokenData.ExpiresAt,
	})
	if err != nil {
		return fmt.Errorf("create verification token: %w", err)
	}

	body := emailusecase.VerificationEmailHTML(tokenData.Token, s.baseURL)

	if err := s.emailService.SendWithRetry(ctx, []string{user.Email}, "verify account", body); err != nil {
		return fmt.Errorf("send verification email: %w", err)
	}

	return nil
}

func (s *CreateUserService) dispatchVerificationEmail(user *entity.User) {
	// to not block req
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("panic in verification email goroutine",
				"user_id", user.ID,
				"panic", r,
				"stack", string(debug.Stack()),
			)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.sendVerificationEmail(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "failed to send verification email",
			"user_id", user.ID,
			"error", err,
		)
	}
}
