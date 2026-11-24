package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/rober0xf/notifier/internal/middleware"
)

// VerifyEmail godoc
// @Summary      Verify user email
// @Description  Verifies a user's email address using the token sent to their mail.
// @Description  Possible 422 errors: invalid or expired verification token.
// @Tags         users
// @Produce      json
// @Param        token  path      string  true  "Verification token"
// @Success      200    {object}  dto.VerifyEmailResponse
// @Failure      422    {object}  dto.ErrorResponse  "Invalid or expired token"
// @Failure      500    {object}  dto.ErrorResponse  "Internal server error"
// @Router       /v1/users/email_verification/{token} [get]
func (h *UserHandler) VerifyEmail(c *gin.Context) {
	token := c.Param("token")

	_, err := h.verifyEmailService.Execute(c.Request.Context(), token)
	if err != nil {
		switch {
		case errors.Is(err, middleware.ErrInvalidToken):
			c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: "invalid or expired verification link"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to verify user", "error", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, dto.VerifyEmailResponse{
		Message: "email verified successfully",
	})
}
