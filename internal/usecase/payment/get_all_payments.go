package payment

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/domain/repository"
)

type GetAllPaymentsService struct {
	paymentRepo repository.PaymentRepository
	logger      *slog.Logger
}

func NewGetAllPaymentsService(paymentRepo repository.PaymentRepository, logger *slog.Logger) *GetAllPaymentsService {
	if logger == nil {
		logger = slog.Default()
	}

	return &GetAllPaymentsService{
		paymentRepo: paymentRepo,
		logger:      logger,
	}
}

func (s *GetAllPaymentsService) Execute(ctx context.Context) ([]entity.Payment, error) {
	payments, err := s.paymentRepo.GetAllPayments(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all payments: %w", err)
	}
	s.logger.InfoContext(ctx, "fetched all payments",
		"count", len(payments))

	return payments, nil
}
