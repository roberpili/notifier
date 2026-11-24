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
)

type CreatePaymentService struct {
	paymentRepo repository.PaymentRepository
	logger      *slog.Logger
}

func NewCreatePaymentService(paymentRepo repository.PaymentRepository, logger *slog.Logger) *CreatePaymentService {
	if logger == nil {
		logger = slog.Default()
	}

	return &CreatePaymentService{
		paymentRepo: paymentRepo,
		logger:      logger,
	}
}

func (s *CreatePaymentService) Execute(ctx context.Context, userID int, payment *entity.Payment) (*entity.Payment, error) {
	payment.UserID = userID

	created, err := s.paymentRepo.CreatePayment(ctx, payment)
	if err != nil {
		if errors.Is(err, repoErr.ErrAlreadyExists) {
			return nil, domainErr.ErrPaymentAlreadyExists
		}

		return nil, fmt.Errorf("create payment: %w", err)
	}
	s.logger.InfoContext(ctx, "payment created",
		"payment_id", created.ID,
		"user_id", userID)

	return created, nil
}
