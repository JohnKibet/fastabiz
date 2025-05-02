package feedback

type Feedback struct {
	ID          int64  `json:"id"`
	OrderID     int64  `json:"order_id"`
	CustomerID  int64  `json:"customer_id"`
	Rating      int    `json:"rating"`
	Comments    string `json:"comments"`
	SubmittedAt string `json:"submitted_at"`
}
