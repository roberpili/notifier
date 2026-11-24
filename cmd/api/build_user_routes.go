package api

import (
	"log/slog"

	"github.com/rober0xf/notifier/internal/auth"
	routes "github.com/rober0xf/notifier/internal/delivery/http"
	"github.com/rober0xf/notifier/internal/domain/repository"
	email "github.com/rober0xf/notifier/internal/infraestructure/email"
	emailusecase "github.com/rober0xf/notifier/internal/usecase/email"
	"github.com/rober0xf/notifier/internal/usecase/user"
)

type Config struct {
	BaseURL        string
	GoogleClientID string
}

func buildUserRoutes(
	userRepo repository.UserRepository,
	tokenGen auth.TokenGenerator,
	emailSender email.EmailSender,
	config Config,
	logger *slog.Logger) *routes.UserHandler {
	return routes.NewUserHandler(
		user.NewCreateUserService(userRepo, emailSender, emailusecase.MustDisposableEmail(), config.BaseURL, logger),
		user.NewLoginService(userRepo, tokenGen, logger),
		user.NewGetUserByIDService(userRepo, logger),
		user.NewGetUserByEmailService(userRepo, logger),
		user.NewGetAllUsersService(userRepo, logger),
		user.NewUpdateUserService(userRepo, logger),
		user.NewDeleteUserService(userRepo, logger),
		user.NewVerifyEmailService(userRepo, logger),
		user.NewOAuthService(userRepo, tokenGen, logger),
		auth.NewGoogleVerifier(config.GoogleClientID),
	)
}
