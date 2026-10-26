package payment

import (
	"context"
	"errors"
	"fmt"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/rober0xf/notifier/internal/domain/repository"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	authErr "github.com/rober0xf/notifier/pkg/auth"
)

type GetPaymentByIDUseCase struct {
	paymentRepo repository.PaymentRepository
}

func NewGetPaymentByIDUseCase(paymentRepo repository.PaymentRepository) *GetPaymentByIDUseCase {
	return &GetPaymentByIDUseCase{
		paymentRepo: paymentRepo,
	}
}

func (uc *GetPaymentByIDUseCase) Execute(ctx context.Context, id, userID int, userRole entity.Role) (*entity.Payment, error) {
	payment, err := uc.paymentRepo.GetPaymentByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return nil, domainErr.ErrPaymentNotFound
		}

		return nil, fmt.Errorf("GetPaymentByIDUC.Execute failed to get payment by id: %w", err)
	}

	if payment.UserID != userID && userRole != entity.RoleAdmin {
		return nil, authErr.ErrForbidden
	}

	return payment, nil
}
