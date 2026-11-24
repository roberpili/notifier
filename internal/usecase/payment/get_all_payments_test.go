package payment_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetAllPayments(t *testing.T) {
	t.Run("successfully found all payments", func(t *testing.T) {
		s, mockRepo := setupGetAllPaymentsTest(t)

		payments := []*entity.Payment{
			{
				ID:        1,
				Name:      "claude",
				Amount:    100,
				Type:      entity.TransactionTypeExpense,
				Category:  entity.CategoryTypeEducation,
				Date:      "2026-01-01",
				Paid:      true,
				Recurrent: false,
			},
			{
				ID:        2,
				Name:      "paramount",
				Amount:    200,
				Type:      entity.TransactionTypeSubscription,
				Category:  entity.CategoryTypeEntertainment,
				Date:      "2026-01-02",
				Paid:      false,
				Recurrent: false,
			},
			{
				ID:        3,
				Name:      "project",
				Amount:    1500,
				Type:      entity.TransactionTypeIncome,
				Category:  entity.CategoryTypeWork,
				Date:      "2026-01-10",
				Paid:      true,
				Recurrent: false,
			},
		}

		for _, u := range payments {
			idStr := strconv.Itoa(int(u.ID))
			mockRepo.payments[idStr] = u
		}

		result, err := s.Execute(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 3)

		IDs := make(map[int32]bool)
		for _, p := range result {
			IDs[p.ID] = true
		}

		assert.True(t, IDs[1])
		assert.True(t, IDs[2])
		assert.True(t, IDs[3])
	})

	t.Run("returns empty list when no payments exists", func(t *testing.T) {
		s, _ := setupGetAllPaymentsTest(t)

		payments, err := s.Execute(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, payments)
		assert.Empty(t, payments)
		assert.Len(t, payments, 0)
	})

	t.Run("returns payments from different users", func(t *testing.T) {
		s, mockRepo := setupGetAllPaymentsTest(t)

		user1Payment := &entity.Payment{
			ID:       1,
			UserID:   1,
			Name:     "user 1 payment",
			Amount:   100,
			Type:     entity.TransactionTypeExpense,
			Category: entity.CategoryTypeClothing,
			Date:     "2026-01-01",
		}

		user2Payment := &entity.Payment{
			ID:       2,
			UserID:   2,
			Name:     "user 2 payment",
			Amount:   200,
			Type:     entity.TransactionTypeIncome,
			Category: entity.CategoryTypeSports,
			Date:     "2026-01-02",
		}

		mockRepo.payments["1"] = user1Payment
		mockRepo.payments["2"] = user2Payment

		payments, err := s.Execute(context.Background())

		assert.NoError(t, err)
		assert.Len(t, payments, 2)
	})

	t.Run("handles repository errors", func(t *testing.T) {
		s, mockRepo := setupGetAllPaymentsTest(t)

		mockRepo.err = errors.New("database connection failed")

		payments, err := s.Execute(context.Background())

		assert.Error(t, err)
		assert.Nil(t, payments)
		assert.Contains(t, err.Error(), "database connection failed")
	})
}
