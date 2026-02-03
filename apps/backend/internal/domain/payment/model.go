package payment

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string
type PaymentStatus string

const (
	MethodStripe         PaymentMethod = "stripe"
	MethodPayPal         PaymentMethod = "paypal"
	MethodMobileMoney    PaymentMethod = "mobile_money"
	MethodCashOnDelivery PaymentMethod = "cash_on_delivery"

	StatusPending   PaymentStatus = "pending"
	StatusCompleted PaymentStatus = "completed"
	StatusFailed    PaymentStatus = "failed"
)

type Payment struct {
	ID         uuid.UUID `db:"id" json:"id"`
	OrderID    uuid.UUID `db:"order_id" json:"order_id"`
	CustomerID uuid.UUID `db:"customer_id" json:"customer_id"`

	Amount   int64  `db:"amount" json:"amount"` // in cents
	Currency string `db:"currency" json:"currency"`

	Method PaymentMethod `db:"method" json:"method"`
	Status PaymentStatus `db:"status" json:"status"`

	PhoneNumber string `db:"phone_number" json:"phone_number,omitempty"`

	// M-Pesa / Stripe refs
	ProviderRef string `db:"provider_ref" json:"provider_ref,omitempty"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	PaidAt    *time.Time `db:"paid_at" json:"paid_at,omitempty"`
}
