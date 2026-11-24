package payment_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/usecase/payment"
	"github.com/stretchr/testify/assert"

	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
)

func TestUpdatePayment(t *testing.T) {
	t.Run("successful update payment", func(t *testing.T) {
		s, mockRepo := setupUpdatePaymentTest(t)

		paymnt := &entity.Payment{
			UserID:    1,
			ID:        1,
			Name:      "Nike",
			Amount:    100.0,
			Type:      entity.TransactionTypeExpense,
			Category:  entity.CategoryTypeClothing,
			Date:      "2022-10-10",
			Paid:      false,
			Recurrent: false,
		}
		mockRepo.payments["1"] = paymnt

		newName := "Adidas"
		newAmount := 110.0
		input := payment.UpdatePaymentInput{
			Name:   &newName,
			Amount: &newAmount,
		}

		err := s.Execute(context.Background(), 1, 1, input)

		updatedPayment, err := mockRepo.GetPaymentByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, newName, updatedPayment.Name)
		assert.Equal(t, newAmount, updatedPayment.Amount)
		assert.Equal(t, paymnt.Type, updatedPayment.Type)
		assert.Equal(t, paymnt.Category, updatedPayment.Category)
		assert.Equal(t, paymnt.Date, updatedPayment.Date)
		assert.Equal(t, paymnt.Paid, updatedPayment.Paid)

		stored, _ := mockRepo.GetPaymentByID(context.Background(), 1)
		assert.Equal(t, newName, stored.Name)
		assert.Equal(t, newAmount, stored.Amount)
	})

	t.Run("successfully updates only name", func(t *testing.T) {
		s, mockRepo := setupUpdatePaymentTest(t)

		paymnt := &entity.Payment{
			ID:       1,
			UserID:   1,
			Name:     "Original",
			Amount:   50.0,
			Type:     entity.TransactionTypeExpense,
			Category: entity.CategoryTypeElectronics,
			Date:     "2022-10-10",
		}
		mockRepo.payments["1"] = paymnt

		newName := "updated name"
		input := payment.UpdatePaymentInput{
			Name: &newName,
		}

		err := s.Execute(context.Background(), 1, 1, input)
		updatedPayment, err := mockRepo.GetPaymentByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, newName, updatedPayment.Name)
		assert.Equal(t, 50.0, updatedPayment.Amount)
	})

	t.Run("successfully updates recurrent payment frequency", func(t *testing.T) {
		s, mockRepo := setupUpdatePaymentTest(t)

		monthly := entity.FrequencyTypeMonthly
		paymnt := &entity.Payment{
			ID:        1,
			UserID:    1,
			Name:      "NBA",
			Amount:    79.99,
			Type:      entity.TransactionTypeSubscription,
			Category:  entity.CategoryTypeEntertainment,
			Date:      "2022-10-10",
			Recurrent: true,
			Frequency: &monthly,
		}
		mockRepo.payments["1"] = paymnt

		yearly := entity.FrequencyTypeYearly
		input := payment.UpdatePaymentInput{
			Frequency: &yearly,
		}

		err := s.Execute(context.Background(), 1, 1, input)
		updatedPayment, err := mockRepo.GetPaymentByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, yearly, *updatedPayment.Frequency)
	})

	t.Run("returns error when payment not found", func(t *testing.T) {
		s, _ := setupUpdatePaymentTest(t)

		name := "GPT"
		input := payment.UpdatePaymentInput{
			Name: &name,
		}

		err := s.Execute(context.Background(), 99999, 1, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domainErr.ErrPaymentNotFound)
	})

	t.Run("handles repository get errors", func(t *testing.T) {
		s, mockRepo := setupUpdatePaymentTest(t)

		mockRepo.err = errors.New("database connection failed")

		name := "codex"
		input := payment.UpdatePaymentInput{
			Name: &name,
		}

		err := s.Execute(context.Background(), 1, 1, input)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection failed")
	})

	t.Run("handles repository update errors", func(t *testing.T) {
		s, mockRepo := setupUpdatePaymentTest(t)

		mockRepo.err = errors.New("database connection failed")

		paymnt := &entity.Payment{
			ID:       1,
			Name:     "leetcode",
			Amount:   10,
			Type:     entity.TransactionTypeSubscription,
			Category: entity.CategoryTypeEducation,
			Date:     "2022-11-11",
			Paid:     true,
		}
		mockRepo.payments["1"] = paymnt

		newName := "HBO"
		input := payment.UpdatePaymentInput{
			Name: &newName,
		}

		mockRepo.err = errors.New("update failed")

		err := s.Execute(context.Background(), 1, 1, input)

		assert.Error(t, err)
	})
}
