package http

import (
	"github.com/rober0xf/notifier/internal/auth"
	"github.com/rober0xf/notifier/internal/usecase/user"
)

type UserHandler struct {
	createUserService     *user.CreateUserService
	loginService          *user.LoginService
	getUserByIDService    *user.GetUserByIDService
	getUserByEmailService *user.GetUserByEmailService
	getAllUsersService    *user.GetAllUsersService
	updateUserService     *user.UpdateUserService
	deleteUserService     *user.DeleteUserService
	verifyEmailService    *user.VerifyEmailService
	oauthService          *user.GoogleOAuthService
	googleVerifier        auth.GoogleVerifier
}

func NewUserHandler(
	createUserService *user.CreateUserService,
	loginService *user.LoginService,
	getUserByIDService *user.GetUserByIDService,
	getUserByEmailService *user.GetUserByEmailService,
	getAllUsersService *user.GetAllUsersService,
	updateUserService *user.UpdateUserService,
	deleteUserService *user.DeleteUserService,
	verifyEmailService *user.VerifyEmailService,
	oauthService *user.GoogleOAuthService,
	googleVerifier auth.GoogleVerifier,
) *UserHandler {
	return &UserHandler{
		createUserService:     createUserService,
		loginService:          loginService,
		getUserByIDService:    getUserByIDService,
		getUserByEmailService: getUserByEmailService,
		getAllUsersService:    getAllUsersService,
		updateUserService:     updateUserService,
		deleteUserService:     deleteUserService,
		verifyEmailService:    verifyEmailService,
		oauthService:          oauthService,
		googleVerifier:        googleVerifier,
	}
}
