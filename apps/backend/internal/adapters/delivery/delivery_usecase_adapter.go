package delivery

import (
	"context"
	"backend/internal/domain/delivery"
	deliveryusecase "backend/internal/usecase/delivery"

	"github.com/google/uuid"
)

type UseCaseAdapter struct {
	UseCase *deliveryusecase.UseCase
}

func (a *UseCaseAdapter) GetDeliveryByID(ctx context.Context, id uuid.UUID) (*delivery.Delivery, error) {
	return a.UseCase.GetDeliveryByID(ctx, id)
}
