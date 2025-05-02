package payment

type PaymentMethod string
type PaymentStatus string

const (
	MethodStripe         PaymentMethod = "stripe"
	MethodPayPal         PaymentMethod = "paypal"
	MethodMobileMoney    PaymentMethod = "mobile_money"
	MethodCashOnDelivery PaymentMethod = "cash_on_delivery"

	StatusPending   PaymentStatus = " pending"
	StatusCompleted PaymentStatus = "completed"
	StatusFailed    PaymentStatus = "failed"
)

type Payment struct {
	ID      int64         `json:"id"`
	OrderID int64         `json:"order_id"`
	Amount  float64       `json:"amount"`
	Method  PaymentMethod `json:"method"`
	Status  PaymentStatus `json:"status"`
	PaidAt  string        `json:"paid_at"`
}
