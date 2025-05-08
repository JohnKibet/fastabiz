package notification

import (
	"context"
	domain "logistics-backend/internal/domain/notification"

	"github.com/google/uuid"
)

type UseCase struct {
	repo domain.Repository
}

func NewUseCase(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) CreateNotification(ctx context.Context, n *domain.Notification) error {
	return uc.repo.Create(n)
}

func (uc *UseCase) GetNotification(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	return uc.repo.GetByID(id)
}

func (uc *UseCase) ListNotification(ctx context.Context) ([]*domain.Notification, error) {
	return uc.repo.List()
}
