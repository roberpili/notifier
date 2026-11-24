package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/middleware"
)

// Delete godoc
// @Summary      Delete a payment
// @Description  Deletes a payment by ID. Admins can delete any payment, regular users can only delete their own.
// @Tags         payments
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Payment ID"
// @Success      204  "No content"
// @Failure      400  {object}  dto.ErrorResponse  "Invalid payment ID"
// @Failure      401  {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      403  {object}  dto.ErrorResponse  "Forbidden"
// @Failure      404  {object}  dto.ErrorResponse  "Payment not found"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Router       /v1/auth/payments/{id} [delete]
func (h *PaymentHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid payment id"})
		return
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	role, err := middleware.GetRoleFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	err = h.deletePaymentService.Execute(c.Request.Context(), id, userID, role)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrPaymentNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "payment not found"})
		case errors.Is(err, middleware.ErrForbidden):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "cannot delete payment from other user"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to delete payment", "error", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}

		return
	}

	c.Status(http.StatusNoContent)
}
