package payment_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	domainErr "github.com/rober0xf/notifier/internal/domain/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetMyPayments(t *testing.T) {
	t.Run("returns payments for a user", func(t *testing.T) {
		s, mockPaymentRepo, mockUserRepo := setupGetMyPaymentsTest(t)

		email := "richard@gmail.com"
		user := &entity.User{
			ID:       1,
			Email:    email,
			IsActive: true,
		}
		mockUserRepo.users[email] = user
		mockUserRepo.users["1"] = user

		payment1 := &entity.Payment{
			ID:       1,
			UserID:   1,
			Name:     "zed",
			Amount:   30,
			Type:     entity.TransactionTypeSubscription,
			Category: entity.CategoryTypeEducation,
			Date:     "2026-02-13",
		}
		payment2 := &entity.Payment{
			ID:       2,
			UserID:   1,
			Name:     "cursor",
			Amount:   50,
			Type:     entity.TransactionTypeSubscription,
			Category: entity.CategoryTypeEducation,
			Date:     "2026-02-13",
		}
		mockPaymentRepo.payments["1"] = payment1
		mockPaymentRepo.payments["2"] = payment2

		payments, err := s.Execute(context.Background(), user.ID)

		assert.NoError(t, err)
		assert.Len(t, payments, 2)

		for _, p := range payments {
			assert.Equal(t, 1, p.UserID)
		}
	})

	t.Run("returns empty list when user has no payments", func(t *testing.T) {
		s, _, mockUserRepo := setupGetMyPaymentsTest(t)

		user := &entity.User{
			ID:       1,
			Email:    "richard@gmail.com",
			IsActive: true,
		}
		mockUserRepo.users["richard@gmail.com"] = user
		mockUserRepo.users["1"] = user

		payments, err := s.Execute(context.Background(), user.ID)

		assert.NoError(t, err)
		assert.NotNil(t, payments)
		assert.Empty(t, payments)
	})

	t.Run("returns error when user not found", func(t *testing.T) {
		s, _, _ := setupGetMyPaymentsTest(t)

		nonExistentID := 99999

		payments, err := s.Execute(context.Background(), nonExistentID)

		assert.Error(t, err)
		assert.Nil(t, payments)
		assert.ErrorIs(t, err, domainErr.ErrUserNotFound)
	})

	t.Run("returns only payments from specified user", func(t *testing.T) {
		s, mockPaymentRepo, mockUserRepo := setupGetMyPaymentsTest(t)

		user1Email := "richard@gmail.com"
		user1 := &entity.User{
			ID:       1,
			Email:    user1Email,
			IsActive: true,
		}
		mockUserRepo.users[user1Email] = user1
		mockUserRepo.users["1"] = user1

		user2Email := "gilfoyle@gmail.com"
		user2 := &entity.User{
			ID:       2,
			Email:    user1Email,
			IsActive: true,
		}
		mockUserRepo.users[user2Email] = user2
		mockUserRepo.users["2"] = user2

		mockPaymentRepo.payments["1"] = &entity.Payment{
			ID: 1, UserID: 1, Name: "richard's payment",
			Amount: 100, Type: entity.TransactionTypeExpense,
			Category: entity.CategoryTypeElectronics, Date: "2026-01-01",
		}
		mockPaymentRepo.payments["2"] = &entity.Payment{
			ID: 2, UserID: 1, Name: "richard's payment 2",
			Amount: 200, Type: entity.TransactionTypeIncome,
			Category: entity.CategoryTypeSports, Date: "2026-01-02",
		}
		mockPaymentRepo.payments["3"] = &entity.Payment{
			ID: 3, UserID: 2, Name: "gilfoyle's payment",
			Amount: 300, Type: entity.TransactionTypeExpense,
			Category: entity.CategoryTypeWork, Date: "2026-01-03",
		}

		payments, err := s.Execute(context.Background(), user1.ID)

		assert.NoError(t, err)
		assert.Len(t, payments, 2) // just first user

		for _, p := range payments {
			assert.Equal(t, int32(1), int32(p.UserID))
		}
	})

	t.Run("returns error when user is not verified", func(t *testing.T) {
		s, _, mockUserRepo := setupGetMyPaymentsTest(t)

		user := &entity.User{
			ID:       1,
			Email:    "richard@gmail.com",
			IsActive: false,
		}
		mockUserRepo.users["richard@gmail.com"] = user
		mockUserRepo.users["1"] = user

		payments, err := s.Execute(context.Background(), user.ID)
		assert.ErrorIs(t, err, domainErr.ErrUserNotVerified)
		assert.Nil(t, payments)
	})

	t.Run("handles user repository errors", func(t *testing.T) {
		s, _, mockUserRepo := setupGetMyPaymentsTest(t)

		mockUserRepo.err = errors.New("user database connection failed")

		payments, err := s.Execute(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, payments)
	})

	t.Run("handles payment, repository errors", func(t *testing.T) {
		s, mockPaymentRepo, mockUserRepo := setupGetMyPaymentsTest(t)

		user := &entity.User{ID: 1, Email: "richard@gmail.com"}
		mockUserRepo.users["richard@gmail.com"] = user
		mockUserRepo.users["1"] = user

		mockPaymentRepo.err = errors.New("payment database connection error")
		payments, err := s.Execute(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, payments)
	})
}
