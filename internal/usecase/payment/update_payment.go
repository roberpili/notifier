package payment

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	"github.com/rober0xf/notifier/internal/middleware"
)

type UpdatePaymentService struct {
	paymentRepo repository.PaymentRepository
	logger      *slog.Logger
}

func NewUpdatePaymentService(paymentRepo repository.PaymentRepository, logger *slog.Logger) *UpdatePaymentService {
	if logger == nil {
		logger = slog.Default()
	}

	return &UpdatePaymentService{
		paymentRepo: paymentRepo,
		logger:      logger,
	}
}

type UpdatePaymentInput struct {
	Name       *string
	Amount     *float64
	Type       *entity.TransactionType
	Category   *entity.CategoryType
	Date       *string
	DueDate    *string
	Paid       *bool
	PaidAt     *string
	Recurrent  *bool
	Frequency  *entity.FrequencyType
	ReceiptURL *string
}

func (s *UpdatePaymentService) Execute(ctx context.Context, id, userID int, input UpdatePaymentInput) error {
	existingPayment, err := s.paymentRepo.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return domainErr.ErrPaymentNotFound
		}

		return fmt.Errorf("failed to get payment by id: %w", err)
	}

	if userID != existingPayment.UserID {
		return middleware.ErrForbidden
	}

	// partial updates
	applyPaymentUpdates(&input, existingPayment)

	if err := s.paymentRepo.UpdatePayment(ctx, existingPayment); err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return domainErr.ErrPaymentNotFound
		}

		return fmt.Errorf("update payment: %w", err)
	}
	s.logger.InfoContext(ctx, "updated payment",
		"payment_id", existingPayment.ID,
		"user_id", existingPayment.UserID,
	)

	return nil
}

func applyPaymentUpdates(input *UpdatePaymentInput, existing *entity.Payment) {
	if input.Name != nil {
		existing.Name = *input.Name
	}

	if input.Amount != nil {
		existing.Amount = *input.Amount
	}

	if input.Type != nil {
		existing.Type = *input.Type
	}

	if input.Category != nil {
		existing.Category = *input.Category
	}

	if input.Date != nil {
		existing.Date = *input.Date
	}

	if input.DueDate != nil {
		existing.DueDate = input.DueDate
	}

	if input.Recurrent != nil {
		existing.Recurrent = *input.Recurrent
	}

	if input.Frequency != nil {
		existing.Frequency = input.Frequency
	}

	if input.ReceiptURL != nil {
		existing.ReceiptURL = input.ReceiptURL
	}

	if input.Paid != nil {
		existing.Paid = *input.Paid
		if *input.Paid {
			if existing.PaidAt == nil || *existing.PaidAt == "" {
				now := time.Now().Format("2006-01-02")
				existing.PaidAt = &now
			}
		} else {
			existing.PaidAt = nil
		}
	}
}
