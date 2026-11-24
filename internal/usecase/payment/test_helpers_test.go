package payment_test

import (
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/usecase/payment"
)

func setupTestPaymentRepo(t *testing.T) *MockPaymentRepository {
	t.Helper()

	return &MockPaymentRepository{
		payments: make(map[string]*entity.Payment),
	}
}

func setupTestUserRepo(t *testing.T) *MockUserRepositoryForPayment {
	t.Helper()

	return &MockUserRepositoryForPayment{
		users: make(map[string]*entity.User),
	}
}

func setupCreatePaymentTest(t *testing.T) (*payment.CreatePaymentService, *MockPaymentRepository) {
	t.Helper()

	mockRepo := setupTestPaymentRepo(t)

	s := payment.NewCreatePaymentService(mockRepo, nil)

	return s, mockRepo
}

func setupGetPaymentByIDTest(t *testing.T) (*payment.GetPaymentByIDService, *MockPaymentRepository) {
	t.Helper()

	mockRepo := setupTestPaymentRepo(t)
	s := payment.NewGetPaymentByIDService(mockRepo, nil)

	return s, mockRepo
}

func setupGetAllPaymentsTest(t *testing.T) (*payment.GetAllPaymentsService, *MockPaymentRepository) {
	t.Helper()

	mockRepo := setupTestPaymentRepo(t)
	s := payment.NewGetAllPaymentsService(mockRepo, nil)

	return s, mockRepo
}

func setupGetMyPaymentsTest(t *testing.T) (*payment.GetMyPaymentsService, *MockPaymentRepository, *MockUserRepositoryForPayment) {
	t.Helper()

	mockPaymentRepo := setupTestPaymentRepo(t)
	mockUserRepo := setupTestUserRepo(t)
	s := payment.NewGetMyPaymentsService(mockPaymentRepo, mockUserRepo, nil)

	return s, mockPaymentRepo, mockUserRepo
}

func setupDeletePaymentTest(t *testing.T) (*payment.DeletePaymentService, *MockPaymentRepository) {
	t.Helper()

	mockRepo := setupTestPaymentRepo(t)
	s := payment.NewDeletePaymentService(mockRepo, nil)

	return s, mockRepo
}

func setupUpdatePaymentTest(t *testing.T) (*payment.UpdatePaymentService, *MockPaymentRepository) {
	t.Helper()

	mockRepo := setupTestPaymentRepo(t)
	s := payment.NewUpdatePaymentService(mockRepo, nil)

	return s, mockRepo
}
