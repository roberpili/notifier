package http

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/middleware"
)

// Create godoc
// @Summary      Create a new payment
// @Description  Creates a new payment for the authenticated user.
// @Description  Possible 400 errors: invalid category, invalid transaction type, invalid frequency, frequency required if recurrent is set, recurrent required if frequency is set, paid required if paid_at is set, due_date cannot be before payment date.
// @Description  Possible 401 errors: missing or invalid user ID in context.
// @Description  Possible 409 errors: payment already exists.
// @Tags         payments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload  body      dto.CreatePaymentRequest      true  "Payment creation payload"
// @Success      201      {object}  dto.PaymentResponse           "Payment created successfully"
// @Failure      400      {object}  dto.PaymentValidationErrorResponse   "Validation or domain error"
// @Failure      401      {object}  dto.ErrorResponse             "Unauthorized"
// @Failure      409      {object}  dto.ErrorResponse             "Payment already exists"
// @Failure      500      {object}  dto.ErrorResponse             "Internal server error"
// @Router       /v1/auth/payments [post]
func (h *PaymentHandler) Create(c *gin.Context) {
	var req dto.CreatePaymentRequest

	validationMeta := dto.ValidationMeta{
		AllowedTypes:       entity.AllTransactionTypes,
		AllowedCategories:  entity.AllCategoryTypes,
		AllowedFrequencies: entity.AllFrequencyTypes,
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.PaymentValidationErrorResponse{
			Error: formatValidationError(err),
			Meta:  validationMeta,
		})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, dto.PaymentValidationErrorResponse{
			Error: err.Error(),
			Meta:  validationMeta,
		})
		return
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	paidAt := strPtrOrNil(req.PaidAt)
	if req.Paid && paidAt == nil {
		now := time.Now().Format("2006-01-02")
		paidAt = &now
	}

	// some fields are pointers because they are not always needed
	payment := &entity.Payment{
		Amount:     req.Amount,
		Name:       req.Name,
		Type:       req.Type,
		Category:   req.Category,
		Date:       req.Date,
		DueDate:    strPtrOrNil(req.DueDate),
		Paid:       req.Paid,
		PaidAt:     paidAt,
		Recurrent:  req.Recurrent,
		Frequency:  freqPtrOrNil(req.Frequency),
		ReceiptURL: strPtrOrNil(req.ReceiptURL),
	}

	createdPayment, err := h.createPaymentService.Execute(c.Request.Context(), userID, payment)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrPaymentAlreadyExists):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "payment already exists"})
		default:
			slog.ErrorContext(c.Request.Context(), "failed to create payment", "error", err)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}

		return
	}

	c.JSON(http.StatusCreated, dto.ToPaymentResponse(*createdPayment))
}
