package payment_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	repoErr "github.com/rober0xf/notifier/internal/infraestructure/errors"
	"github.com/stretchr/testify/assert"
)

func TestCreatePayment(t *testing.T) {
	t.Run("successfully creates a new payment", func(t *testing.T) {
		s, mockRepo := setupCreatePaymentTest(t)

		paidAt := "2026-03-12"
		input := &entity.Payment{
			UserID:    1,
			Name:      "copilot",
			Amount:    100,
			Type:      entity.TransactionTypeExpense,
			Category:  entity.CategoryTypeEducation,
			Date:      "2026-02-10",
			Paid:      true,
			Recurrent: false,
			PaidAt:    &paidAt,
		}

		payment, err := s.Execute(context.Background(), input.UserID, input)

		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.Equal(t, "copilot", payment.Name)
		assert.Equal(t, 100.0, payment.Amount)
		assert.Equal(t, entity.TransactionTypeExpense, payment.Type)
		assert.Equal(t, entity.CategoryTypeEducation, payment.Category)
		assert.NotEqual(t, 0, payment.ID)

		storedPayment, err := mockRepo.GetPaymentByID(context.Background(), int(payment.ID))
		assert.NoError(t, err)
		assert.Equal(t, payment.Name, storedPayment.Name)
	})

	t.Run("returns error when payment already exists", func(t *testing.T) {
		s, mockRepo := setupCreatePaymentTest(t)

		mockRepo.err = repoErr.ErrAlreadyExists
		payment := &entity.Payment{
			UserID:    1,
			Name:      "infra",
			Amount:    3500,
			Type:      entity.TransactionTypeIncome,
			Category:  entity.CategoryTypeWork,
			Date:      "2025-11-07",
			Paid:      false,
			Recurrent: false,
		}

		_, err := s.Execute(context.Background(), payment.UserID, payment)
		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrPaymentAlreadyExists)
	})

	t.Run("successfully creates recurrent payment with frequency", func(t *testing.T) {
		s, _ := setupCreatePaymentTest(t)

		frequency := entity.FrequencyTypeWeekly
		paidAt := "2026-06-12"
		input := &entity.Payment{
			UserID:    1,
			Name:      "copilot",
			Amount:    100,
			Type:      entity.TransactionTypeExpense,
			Category:  entity.CategoryTypeEducation,
			Date:      "2026-02-10",
			Paid:      true,
			Recurrent: true,
			Frequency: &frequency,
			PaidAt:    &paidAt,
		}

		payment, err := s.Execute(context.Background(), input.UserID, input)

		assert.NoError(t, err)
		assert.NotNil(t, payment)
		assert.True(t, payment.Recurrent)
		assert.NotNil(t, payment.Frequency)
		assert.Equal(t, entity.FrequencyTypeWeekly, *payment.Frequency)
	})

	t.Run("handles repository errors", func(t *testing.T) {
		s, mockRepo := setupCreatePaymentTest(t)

		mockRepo.err = errors.New("database connection failed")

		input := &entity.Payment{
			UserID:   1,
			Name:     "Test",
			Amount:   100,
			Type:     entity.TransactionTypeExpense,
			Category: entity.CategoryTypeClothing,
			Date:     "2026-02-10",
		}

		payment, err := s.Execute(context.Background(), 1, input)

		assert.Error(t, err)
		assert.Nil(t, payment)
	})
}
