package feedback

import (
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	ID          uuid.UUID `db:"id" json:"id"`
	OrderID     uuid.UUID `db:"order_id" json:"order_id"`
	CustomerID  uuid.UUID `db:"customer_id" json:"customer_id"`
	Rating      int       `db:"rating" json:"rating"`
	Comments    string    `db:"comments" json:"comments"`
	SubmittedAt time.Time `db:"submitted_at" json:"submitted_at"`
}
