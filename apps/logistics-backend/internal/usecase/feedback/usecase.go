package feedback

import (
	"context"
	domain "logistics-backend/internal/domain/feedback"

	"github.com/google/uuid"
)

type UseCase struct {
	repo domain.Repository
}

func NewUseCase(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) CreateFeedback(ctx context.Context, f *domain.Feedback) error {
	return uc.repo.Create(f)
}

func (uc *UseCase) GetFeedbackByID(ctx context.Context, id uuid.UUID) (*domain.Feedback, error) {
	return uc.repo.GetByID(id)
}

func (uc *UseCase) ListFeedback(ctx context.Context) ([]*domain.Feedback, error) {
	return uc.repo.List()
}
