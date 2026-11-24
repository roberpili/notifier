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

// Update godoc
// @Summary      Update a payment
// @Description  Updates a payment by ID. Users can only update their own payments.
// @Tags         payments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                        true  "Payment ID"
// @Param        payload  body      dto.UpdatePaymentRequest   true  "Update payload"
// @Success      204  "No content"
// @Failure      400  {object}  dto.PaymentValidationErrorResponse  "Validation error or invalid id"
// @Failure      401  {object}  dto.ErrorResponse            "Unauthorized"
// @Failure      403  {object}  dto.ErrorResponse            "Forbidden"
// @Failure      404  {object}  dto.ErrorResponse            "Payment not found"
// @Failure      500  {object}  dto.ErrorResponse            "Internal server error"
// @Router       /v1/auth/payments/{id} [put]
func (h *PaymentHandler) Update(c *gin.Context) {
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

	var req dto.UpdatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.PaymentValidationErrorResponse{
			Error: formatValidationError(err),
		})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.PaymentValidationErrorResponse{Error: err.Error()})
		return
	}

	err = h.updatePaymentService.Execute(c.Request.Context(), id, userID, req.ToInput())
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrPaymentNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "payment not found"})
		case errors.Is(err, middleware.ErrForbidden):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "cannot update another user's payment"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to update payment by id", "error", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}

		return
	}

	c.Status(http.StatusNoContent)
}
