package api

import (
	"log/slog"

	routes "github.com/rober0xf/notifier/internal/delivery/http"
	"github.com/rober0xf/notifier/internal/domain/repository"
	"github.com/rober0xf/notifier/internal/usecase/payment"
)

func buildPaymentRoutes(
	paymentRepo repository.PaymentRepository,
	userRepo payment.UserIDGetter,
	logger *slog.Logger) *routes.PaymentHandler {
	return routes.NewPaymentHandler(
		payment.NewCreatePaymentService(paymentRepo, logger),
		payment.NewGetPaymentByIDService(paymentRepo, logger),
		payment.NewGetMyPaymentsService(paymentRepo, userRepo, logger),
		payment.NewGetAllPaymentsService(paymentRepo, logger),
		payment.NewUpdatePaymentService(paymentRepo, logger),
		payment.NewDeletePaymentService(paymentRepo, logger),
	)
}
