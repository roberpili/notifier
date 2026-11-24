package integration

import (
	"io"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rober0xf/notifier/internal/auth"
	"github.com/rober0xf/notifier/internal/delivery/http"
	"github.com/rober0xf/notifier/internal/domain/repository"
	"github.com/rober0xf/notifier/internal/infraestructure/persistance/postgres"
	"github.com/rober0xf/notifier/internal/middleware"
	"github.com/rober0xf/notifier/internal/session"
	email "github.com/rober0xf/notifier/internal/usecase/email"
	"github.com/rober0xf/notifier/internal/usecase/payment"
	"github.com/rober0xf/notifier/internal/usecase/user"

	"testing"
)

type TestUserDependencies struct {
	router   *gin.Engine
	db       *pgxpool.Pool
	userRepo repository.UserRepository
	jwtGen   *auth.JWTGenerator
}

type TestPaymentDependencies struct {
	router      *gin.Engine
	db          *pgxpool.Pool
	paymentRepo repository.PaymentRepository
	userRepo    payment.UserIDGetter
}

type MockGoogleVerifier struct{}

func setupTestDependencies(t *testing.T) (*TestUserDependencies, *TestPaymentDependencies) {
	t.Helper()
	db := setupTestDB(t)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	userRepo := postgres.NewUserRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)

	jwtKey := []byte("test-jwt")
	tokenGen := auth.NewJWTGenerator(jwtKey, 160)
	authMiddleware := middleware.AuthMiddleware(tokenGen, session.SessionCookieName)

	emailSender := email.NewMockSender()

	createUserService := user.NewCreateUserService(userRepo, emailSender, []string{}, "http://localhost:8080", logger)
	loginUserService := user.NewLoginService(userRepo, tokenGen, logger)
	getUserByIDService := user.NewGetUserByIDService(userRepo, logger)
	getUserByEmailService := user.NewGetUserByEmailService(userRepo, logger)
	getAllUsersService := user.NewGetAllUsersService(userRepo, logger)
	updateUserService := user.NewUpdateUserService(userRepo, logger)
	verifyEmailService := user.NewVerifyEmailService(userRepo, logger)
	deleteUserService := user.NewDeleteUserService(userRepo, logger)
	oauthUC := user.NewOAuthService(userRepo, tokenGen, logger)
	googleVerifier := &MockGoogleVerifier{}

	userHandler := http.NewUserHandler(
		createUserService,
		loginUserService,
		getUserByIDService,
		getUserByEmailService,
		getAllUsersService,
		updateUserService,
		deleteUserService,
		verifyEmailService,
		oauthUC,
		googleVerifier,
	)

	createPaymentService := payment.NewCreatePaymentService(paymentRepo, logger)
	getPaymentByIDService := payment.NewGetPaymentByIDService(paymentRepo, logger)
	getMyPaymentsService := payment.NewGetMyPaymentsService(paymentRepo, userRepo, logger)
	getAllPaymentsService := payment.NewGetAllPaymentsService(paymentRepo, logger)
	updatePaymentService := payment.NewUpdatePaymentService(paymentRepo, logger)
	deletePaymentService := payment.NewDeletePaymentService(paymentRepo, logger)

	paymentHandler := http.NewPaymentHandler(
		createPaymentService,
		getPaymentByIDService,
		getMyPaymentsService,
		getAllPaymentsService,
		updatePaymentService,
		deletePaymentService,
	)

	gin.SetMode(gin.TestMode)
	router := http.SetupRoutes(userHandler, paymentHandler, authMiddleware)

	return &TestUserDependencies{
		router:   router,
		db:       db,
		userRepo: userRepo,
		jwtGen:   tokenGen,
	}, &TestPaymentDependencies{
		router:      router,
		db:          db,
		paymentRepo: paymentRepo,
		userRepo:    userRepo,
	}
}

func (m *MockGoogleVerifier) Verify(idToken string) (*auth.GoogleUser, error) {
	return &auth.GoogleUser{
		Sub:   "test-google-id",
		Email: "test@example.com",
		Name:  "Test User",
	}, nil
}
