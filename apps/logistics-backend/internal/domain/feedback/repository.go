package feedback

import (
	"github.com/google/uuid"
)

type Repository interface {
	Create(feedback *Feedback) error
	GetByID(id uuid.UUID) (*Feedback, error)
	List() ([]*Feedback, error)
}
