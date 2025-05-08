package payment

import "github.com/google/uuid"

type CreatePaymentRequest struct {
	OrderID uuid.UUID     `json:"order_id"`
	Amount  float64       `json:"amount"`
	Method  PaymentMethod `json:"method"`
	Status  PaymentStatus `json:"status"`
}

func (r *CreatePaymentRequest) ToPayment() *Payment {
	return &Payment{
		OrderID: r.OrderID,
		Amount:  r.Amount,
		Method:  r.Method,
		Status:  StatusPending,
	}
}
