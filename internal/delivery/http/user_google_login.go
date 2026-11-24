package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/auth"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/session"
)

// GoogleLogin godoc
// @Summary      Google OAuth login
// @Description  Authenticates a user using a Google ID token.
// @Description  Possible 400 errors: invalid google ID, invalid email format.
// @Description  Possible 409 errors: google account already linked to another account.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.GoogleLoginRequest  true  "Google ID token"
// @Success      200      {object}  dto.LoginResponse
// @Failure      400      {object}  dto.ErrorResponse  "Invalid user data"
// @Failure      401      {object}  dto.ErrorResponse  "Invalid google token"
// @Failure      409      {object}  dto.ErrorResponse  "Google account already linked"
// @Failure      500      {object}  dto.ErrorResponse  "Internal server error"
// @Router       /v1/users/login/google [post]
func (h *UserHandler) GoogleLogin(c *gin.Context) {
	var req dto.GoogleLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id_token is required"})
		return
	}

	// validate google token
	googleData, err := h.googleVerifier.Verify(req.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid google token"})
		return
	}

	out, err := h.oauthService.Execute(c.Request.Context(), googleData.Sub, googleData.Email, googleData.Name)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrInvalidGoogleID), errors.Is(err, domainErr.ErrInvalidEmailFormat):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user data"})
		case errors.Is(err, auth.ErrGoogleAccountAlreadyLinked):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "google account already linked"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to login with oauth", "error", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}

		return
	}

	session.SetAuthCookie(c, out.Token, session.CookieConfig{
		Name:            session.SessionCookieName,
		TokenExpiration: session.TokenExpiration,
		Secure:          false,
		HttpOnly:        true,
	})

	c.JSON(http.StatusOK, dto.LoginResponse{
		ID:    out.User.ID,
		Email: out.User.Email,
	})
}
