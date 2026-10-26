package payment_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	authErr "github.com/rober0xf/notifier/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func TestGetPaymentByID(t *testing.T) {
	t.Run("succesfully found the payment by id", func(t *testing.T) {
		uc, mockRepo := setupGetPaymentByIDTest(t)

		expectedPayment := &entity.Payment{
			ID:        1,
			UserID:    1,
			Name:      "copilot",
			Amount:    100,
			Type:      entity.TransactionTypeExpense,
			Category:  entity.CategoryTypeEducation,
			Date:      "2026-02-10",
			Paid:      true,
			Recurrent: false,
		}
		mockRepo.payments["1"] = expectedPayment

		payment, err := uc.Execute(context.Background(), int(expectedPayment.ID), 1, entity.RoleUser)

		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.Equal(t, expectedPayment.ID, payment.ID)
		assert.Equal(t, expectedPayment.Name, payment.Name)
		assert.Equal(t, expectedPayment.Amount, payment.Amount)
		assert.Equal(t, entity.TransactionTypeExpense, payment.Type)
		assert.Equal(t, entity.CategoryTypeEducation, payment.Category)
	})

	t.Run("returns error when payment not found", func(t *testing.T) {
		uc, _ := setupGetPaymentByIDTest(t)

		nonExistingID := 99999
		payment, err := uc.Execute(context.Background(), nonExistingID, 1, entity.RoleUser)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.ErrorIs(t, err, domainErr.ErrPaymentNotFound)
	})

	t.Run("handles repository error", func(t *testing.T) {
		uc, mockRepo := setupGetPaymentByIDTest(t)

		mockRepo.err = errors.New("database connection failed")
		payment, err := uc.Execute(context.Background(), 1, 1, entity.RoleUser)

		assert.Error(t, err)
		assert.Nil(t, payment)
		assert.Contains(t, err.Error(), "database connection failed")
	})

	t.Run("admins can access payments from other users", func(t *testing.T) {
		uc, mockRepo := setupGetPaymentByIDTest(t)

		expectedPayment := &entity.Payment{
			ID:        1,
			UserID:    99,
			Name:      "netflix",
			Amount:    15,
			Type:      entity.TransactionTypeSubscription,
			Category:  entity.CategoryTypeEntertainment,
			Date:      "2026-02-10",
			Paid:      false,
			Recurrent: true,
		}
		mockRepo.payments["1"] = expectedPayment

		payment, err := uc.Execute(
			context.Background(),
			int(expectedPayment.ID),
			1,
			entity.RoleAdmin,
		)

		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.Equal(t, expectedPayment.ID, payment.ID)
	})

	t.Run("regular users cannot access payments from other users", func(t *testing.T) {
		uc, mockRepo := setupGetPaymentByIDTest(t)

		mockRepo.payments["1"] = &entity.Payment{
			ID:     1,
			UserID: 99,
			Name:   "netflix",
		}

		payment, err := uc.Execute(
			context.Background(),
			1,
			2,
			entity.RoleUser,
		)

		assert.ErrorIs(t, err, authErr.ErrForbidden)
		assert.Nil(t, payment)
	})
}
