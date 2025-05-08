package feedback

import "github.com/google/uuid"

type CreateFeedbackRequest struct {
	OrderID    uuid.UUID `json:"order_id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Rating     int       `json:"rating"`
	Comments   string    `json:"comments"`
}

func (r *CreateFeedbackRequest) ToFeedback() *Feedback {
	return &Feedback{
		OrderID:    r.OrderID,
		CustomerID: r.CustomerID,
		Rating:     r.Rating,
		Comments:   r.Comments,
	}
}
