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

type GetPaymentByIDService struct {
	paymentRepo repository.PaymentRepository
	logger      *slog.Logger
}

func NewGetPaymentByIDService(paymentRepo repository.PaymentRepository, logger *slog.Logger) *GetPaymentByIDService {
	if logger == nil {
		logger = slog.Default()
	}

	return &GetPaymentByIDService{
		paymentRepo: paymentRepo,
		logger:      logger,
	}
}

func (s *GetPaymentByIDService) Execute(ctx context.Context, id, userID int, userRole entity.Role) (*entity.Payment, error) {
	payment, err := s.paymentRepo.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, domainErr.ErrPaymentNotFound
		}

		return nil, fmt.Errorf("get payment by id: %w", err)
	}

	if payment.UserID != userID && userRole != entity.RoleAdmin {
		return nil, middleware.ErrForbidden
	}

	s.logger.InfoContext(ctx, "fetched payment by id",
		"payment_id", payment.ID,
		"user_id", userID)

	return payment, nil
}
