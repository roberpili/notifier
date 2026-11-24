package http

import "github.com/rober0xf/notifier/internal/usecase/payment"

type PaymentHandler struct {
	createPaymentService  *payment.CreatePaymentService
	getPaymentByIDService *payment.GetPaymentByIDService
	getMyPaymentsService  *payment.GetMyPaymentsService
	getAllPaymentsService *payment.GetAllPaymentsService
	updatePaymentService  *payment.UpdatePaymentService
	deletePaymentService  *payment.DeletePaymentService
}

func NewPaymentHandler(
	createPaymentService *payment.CreatePaymentService,
	getPaymentByIDService *payment.GetPaymentByIDService,
	getMyPaymentsService *payment.GetMyPaymentsService,
	getAllPaymentsService *payment.GetAllPaymentsService,
	updatePaymentService *payment.UpdatePaymentService,
	deletePaymentService *payment.DeletePaymentService,
) *PaymentHandler {
	return &PaymentHandler{
		createPaymentService:  createPaymentService,
		getPaymentByIDService: getPaymentByIDService,
		getMyPaymentsService:  getMyPaymentsService,
		getAllPaymentsService: getAllPaymentsService,
		updatePaymentService:  updatePaymentService,
		deletePaymentService:  deletePaymentService,
	}
}
