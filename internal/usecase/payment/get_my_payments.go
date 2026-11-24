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

type UserIDGetter interface {
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)
}

type GetMyPaymentsService struct {
	paymentRepo repository.PaymentRepository
	userRepo    UserIDGetter
	logger      *slog.Logger
}

func NewGetMyPaymentsService(paymentRepo repository.PaymentRepository, userRepo UserIDGetter, logger *slog.Logger) *GetMyPaymentsService {
	if logger == nil {
		logger = slog.Default()
	}

	return &GetMyPaymentsService{
		paymentRepo: paymentRepo,
		userRepo:    userRepo,
		logger:      logger,
	}
}

func (s *GetMyPaymentsService) Execute(ctx context.Context, userID int) ([]entity.Payment, error) {
	foundUser, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, domainErr.ErrUserNotFound
		}

		return nil, fmt.Errorf("get user %d: %w", userID, err)
	}

	if !foundUser.IsActive {
		return nil, domainErr.ErrUserNotVerified
	}

	payments, err := s.paymentRepo.GetMyPayments(ctx, foundUser.ID)
	if err != nil {
		return nil, fmt.Errorf("user %d: %w", userID, err)
	}
	s.logger.InfoContext(ctx, "fetched my payments",
		"user_id", userID,
		"count", len(payments))

	return payments, nil
}
