package user_test

import (
	"testing"

	"github.com/rober0xf/notifier/internal/domain/entity"
	"github.com/rober0xf/notifier/internal/usecase/email"
	"github.com/rober0xf/notifier/internal/usecase/user"
)

func setupTestRepo(t *testing.T) *MockUserRepository {
	t.Helper()

	return &MockUserRepository{
		users:     make(map[int]*entity.User),
		emails:    make(map[string]*entity.User),
		usernames: make(map[string]*entity.User),
		tokens:    make(map[string]*entity.UserToken),
	}
}

func setupCreateUserTest(t *testing.T) (*user.CreateUserService, *MockUserRepository, *email.MockSender) {
	t.Helper()

	mockRepo := setupTestRepo(t)

	mockEmailSender := email.NewMockSender()

	disposableDomains := []string{}
	baseURL := "http://localhost:3000"

	s := user.NewCreateUserService(mockRepo, mockEmailSender, disposableDomains, baseURL, nil)

	return s, mockRepo, mockEmailSender
}

func setupGetUserByEmailTest(t *testing.T) (*user.GetUserByEmailService, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	s := user.NewGetUserByEmailService(mockRepo, nil)

	return s, mockRepo
}

func setupGetUserByIDTest(t *testing.T) (*user.GetUserByIDService, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	s := user.NewGetUserByIDService(mockRepo, nil)

	return s, mockRepo
}

func setupGetAllUsersTest(t *testing.T) (*user.GetAllUsersService, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	s := user.NewGetAllUsersService(mockRepo, nil)

	return s, mockRepo
}

func setupDeleteUserTest(t *testing.T) (*user.DeleteUserService, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	s := user.NewDeleteUserService(mockRepo, nil)

	return s, mockRepo
}

func setupUpdateUserTest(t *testing.T) (*user.UpdateUserService, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	s := user.NewUpdateUserService(mockRepo, nil)

	return s, mockRepo
}

func setupVerifyEmailTest(t *testing.T) (*user.VerifyEmailService, *MockUserRepository) {
	t.Helper()

	mockRepo := setupTestRepo(t)
	s := user.NewVerifyEmailService(mockRepo, nil)

	return s, mockRepo
}
