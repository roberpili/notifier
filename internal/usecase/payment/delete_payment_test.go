package payment_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/middleware"
	"github.com/stretchr/testify/assert"

	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
)

func TestDeletePayment(t *testing.T) {
	t.Run("successfully deleted payment as owner", func(t *testing.T) {
		s, mockRepo := setupDeletePaymentTest(t)

		monthly := entity.FrequencyTypeMonthly
		payment := &entity.Payment{
			ID:        1,
			UserID:    1,
			Name:      "fight pass",
			Amount:    110,
			Type:      entity.TransactionTypeSubscription,
			Category:  entity.CategoryTypeEntertainment,
			Date:      "2022-11-01",
			Paid:      true,
			Recurrent: true,
			Frequency: &monthly,
		}
		mockRepo.payments["1"] = payment

		err := s.Execute(context.Background(), 1, 1, entity.RoleUser)
		assert.NoError(t, err)

		_, exists := mockRepo.payments["1"]
		assert.False(t, exists)
	})

	t.Run("successfully deleted payment as admin", func(t *testing.T) {
		s, mockRepo := setupDeletePaymentTest(t)

		payment := &entity.Payment{
			ID:     1,
			UserID: 1,
			Name:   "fight pass",
			Amount: 110,
		}
		mockRepo.payments["1"] = payment

		err := s.Execute(context.Background(), 1, 99, entity.RoleAdmin)
		assert.NoError(t, err)

		_, exists := mockRepo.payments["1"]
		assert.False(t, exists)
	})

	t.Run("returns forbidden when user tries to delete another user payment", func(t *testing.T) {
		s, mockRepo := setupDeletePaymentTest(t)

		payment := &entity.Payment{
			ID:     1,
			UserID: 1,
			Name:   "fight pass",
			Amount: 110,
		}
		mockRepo.payments["1"] = payment

		err := s.Execute(context.Background(), 1, 2, entity.RoleUser)
		assert.ErrorIs(t, err, middleware.ErrForbidden)
	})

	t.Run("returns error when payment not found", func(t *testing.T) {
		s, _ := setupDeletePaymentTest(t)

		err := s.Execute(context.Background(), 99999, 1, entity.RoleUser)
		assert.ErrorIs(t, err, domainErr.ErrPaymentNotFound)
	})

	t.Run("handles repository errors", func(t *testing.T) {
		s, mockRepo := setupDeletePaymentTest(t)

		mockRepo.err = errors.New("database connection failed")
		payment := &entity.Payment{
			ID:        1,
			UserID:    1,
			Name:      "neetcode",
			Amount:    199,
			Type:      entity.TransactionTypeExpense,
			Category:  entity.CategoryTypeEducation,
			Date:      "2024-07-01",
			Paid:      false,
			Recurrent: false,
		}
		mockRepo.payments["1"] = payment

		err := s.Execute(context.Background(), 1, 1, entity.RoleUser)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection failed")

		_, exists := mockRepo.payments["1"]
		assert.True(t, exists)
	})
}
