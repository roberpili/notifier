package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/rober0xf/notifier/pkg/auth"
	authErr "github.com/rober0xf/notifier/pkg/auth"
)

// Login godoc
// @Summary      Login
// @Description  Authenticates a user and sets a session cookie.
// @Description  Possible 401 errors: invalid credentials.
// @Description  Possible 403 errors: email not verified.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.LoginPayload  true  "Login credentials"
// @Success      200      {object}  dto.LoginResponse
// @Failure      400      {object}  dto.ErrorResponse  "Missing fields"
// @Failure      401      {object}  dto.ErrorResponse  "Invalid credentials"
// @Failure      403      {object}  dto.ErrorResponse  "Email not verified"
// @Failure      500      {object}  dto.ErrorResponse  "Internal server error"
// @Router       /v1/users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var payload dto.LoginPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "email and password are required"})
		return
	}

	out, err := h.loginUC.Execute(c.Request.Context(), payload.Email, payload.Password)
	if err != nil {
		switch {
		case errors.Is(err, authErr.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid credentials"})
		case errors.Is(err, authErr.ErrEmailNotVerified):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "email not verified"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to login user", "error", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}

		return
	}

	auth.SetAuthCookie(c, out.Token, auth.CookieConfig{
		Name:            auth.SessionCookieName,
		TokenExpiration: auth.TokenExpiration,
		Secure:          false,
		HttpOnly:        true,
	})

	c.JSON(http.StatusOK, dto.LoginResponse{
		ID:    out.User.ID,
		Email: out.User.Email,
		Role:  string(out.User.Role),
	})
}
