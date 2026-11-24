package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/usecase/user"
)

// Create godoc
// @Summary      Register a new user
// @Description  Creates a new user account and sends a verification email.
// @Description  Possible 400 errors: invalid email format, invalid domain, disposable email, email not reachable, weak password, missing fields.
// @Description  Possible 409 errors: email already in use, username already in use.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.RegisterPayload          true  "Registration data"
// @Success      201      {object}  dto.UserCreatedResponse
// @Failure      400      {object}  dto.UserValidationErrorResponse  "Validation or domain error"
// @Failure      409      {object}  dto.ErrorResponse            "Conflict: email or username already exists"
// @Failure      500      {object}  dto.ErrorResponse            "Internal server error"
// @Router       /v1/users/register [post]
func (h *UserHandler) Create(c *gin.Context) {
	var payload dto.RegisterPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "username, email and password are required. password min length 8",
		})
		return
	}

	// here we use the service logic
	user, err := h.createUserService.Execute(c.Request.Context(), payload.Username, payload.Email, payload.Password)
	if err != nil {
		handleCreateUserError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.UserCreatedResponse{
		Message: "check your email to verify your account",
		User: dto.UserPayload{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Active:   user.IsActive,
		},
	})
}

func handleCreateUserError(c *gin.Context, err error) {
	var validationErr *user.EmailValidationError

	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, dto.UserValidationErrorResponse{
			Error:      validationErr.Message,
			Suggestion: validationErr.Suggestion,
		})
		return
	}

	switch {
	case errors.Is(err, domainErr.ErrEmailAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
	case errors.Is(err, domainErr.ErrUsernameAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": "username already in use"})
	case errors.Is(err, domainErr.ErrInvalidEmailFormat):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
	case errors.Is(err, domainErr.ErrInvalidDomain):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email domain"})
	case errors.Is(err, domainErr.ErrDisposableEmail):
		c.JSON(http.StatusBadRequest, gin.H{"error": "disposable emails are not valid"})
	case errors.Is(err, domainErr.ErrEmailNotReachable):
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is not reachable"})
	case errors.Is(err, domainErr.ErrInvalidPassword):
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be stronger"})
	default:
		slog.ErrorContext(c.Request.Context(), "failed to register user", "error", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
	}
}
