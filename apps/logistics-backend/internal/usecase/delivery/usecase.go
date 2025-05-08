package delivery

import (
	"context"
	domain "logistics-backend/internal/domain/delivery"

	"github.com/google/uuid"
)

type UseCase struct {
	repo domain.Repository
}

func NewUseCase(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) CreateDelivery(ctx context.Context, d *domain.Delivery) error {
	return uc.repo.Create(d)
}

func (uc *UseCase) GetDeliveryByID(ctx context.Context, id uuid.UUID) (*domain.Delivery, error) {
	return uc.repo.GetByID(id)
}

func (uc *UseCase) ListDeliveries(ctx context.Context) ([]*domain.Delivery, error) {
	return uc.repo.List()
}
