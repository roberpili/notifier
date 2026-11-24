package payment

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	"github.com/rober0xf/notifier/internal/middleware"
)

type DeletePaymentService struct {
	paymentRepo repository.PaymentRepository
	logger      *slog.Logger
}

func NewDeletePaymentService(paymentRepo repository.PaymentRepository, logger *slog.Logger) *DeletePaymentService {
	if logger == nil {
		logger = slog.Default()
	}

	return &DeletePaymentService{
		paymentRepo: paymentRepo,
		logger:      logger,
	}
}

func (s *DeletePaymentService) Execute(ctx context.Context, id, userID int, userRole entity.Role) error {
	payment, err := s.paymentRepo.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return domainErr.ErrPaymentNotFound
		}

		return fmt.Errorf("failed getting payment by id: %w", err)
	}

	if payment.UserID != userID && userRole != entity.RoleAdmin {
		return middleware.ErrForbidden
	}

	if err := s.paymentRepo.DeletePayment(ctx, id); err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return domainErr.ErrPaymentNotFound
		}

		return fmt.Errorf("delete payment: %w", err)
	}
	s.logger.InfoContext(ctx, "payment deleted",
		"payment_id", id,
		"user_id", userID)

	return nil
}
