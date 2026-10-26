package http

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/pkg/auth"
)

// GetByID godoc
// @Summary      Get payment by ID
// @Description  Returns a single payment by ID. Users can only access their own payments.
// @Tags         payments
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Payment ID"
// @Success      200  {object}  dto.PaymentResponse
// @Failure      400  {object}  dto.ErrorResponse  "Invalid payment ID"
// @Failure      401  {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      403  {object}  dto.ErrorResponse  "Forbidden"
// @Failure      404  {object}  dto.ErrorResponse  "Payment not found"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Router       /v1/admin/payments/{id} [get]
func (h *PaymentHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid payment id"})
		return
	}

	userID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	role, err := auth.GetRoleFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	payment, err := h.getPaymentByIDUC.Execute(c.Request.Context(), id, userID, role)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrPaymentNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "payment not found"})
		case errors.Is(err, auth.ErrForbidden):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "cannot access payments from other users"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to get payment by id", "error", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, dto.ToPaymentResponse(*payment))
}

// GetMyPayments godoc
// @Summary      Get my payments
// @Description  Returns all payments for the authenticated user.
// @Tags         payments
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   dto.PaymentResponse
// @Failure      401  {object}  dto.ErrorResponse  "Unauthorized"
// @Failure      403  {object}  dto.ErrorResponse  "Forbidden"
// @Failure      404  {object}  dto.ErrorResponse  "User not found"
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Router       /v1/auth/payments/me [get]
func (h *PaymentHandler) GetMyPayments(c *gin.Context) {
	userID, err := auth.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	payments, err := h.getMyPaymentsUC.Execute(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrUserNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "user not found"})
		case errors.Is(err, domainErr.ErrUserNotVerified):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "user not verified"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to get all payments from user", "error", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}

		return
	}

	response := make([]dto.PaymentResponse, 0, len(payments))
	for _, p := range payments {
		response = append(response, dto.ToPaymentResponse(p))
	}

	c.JSON(http.StatusOK, response)
}

// GetAll godoc
// @Summary      Get all payments
// @Description  Returns all payments. Admin only.
// @Tags         payments
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   dto.PaymentResponse
// @Failure      500  {object}  dto.ErrorResponse  "Internal server error"
// @Router       /v1/admin/payments [get]
func (h *PaymentHandler) GetAll(c *gin.Context) {
	payments, err := h.getAllPaymentsUC.Execute(c.Request.Context())
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "failed to get all payments", "error", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	response := make([]dto.PaymentResponse, 0, len(payments))
	for _, p := range payments {
		response = append(response, dto.ToPaymentResponse(p))
	}

	c.JSON(http.StatusOK, response)
}
