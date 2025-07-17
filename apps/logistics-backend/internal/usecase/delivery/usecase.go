package delivery

import (
	"context"
	"fmt"
	"logistics-backend/internal/domain/delivery"

	"github.com/google/uuid"
)

type UseCase struct {
	repo       delivery.Repository
	orderLogic delivery.OrderReader
}

func NewUseCase(repo delivery.Repository, ord delivery.OrderReader) *UseCase {
	return &UseCase{repo: repo, orderLogic: ord}
}

func (uc *UseCase) CreateDelivery(ctx context.Context, d *delivery.Delivery) error {
	tx, err := uc.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. fetch order
	order, err := uc.orderLogic.GetByID(ctx, d.OrderID)
	if err != nil {
		return fmt.Errorf("could not fetch order: %w", err)
	}
	if order.OrderStatus != "pending" {
		return delivery.ErrorNoPendingOrder
	}

	// 2. update order status using tx
	if err := uc.orderLogic.UpdateOrderTx(ctx, tx, order.ID, "status", "assigned"); err != nil {
		return fmt.Errorf("could not update order status: %w", err)
	}

	// 3. create delivery using tx
	if err := uc.repo.CreateTx(ctx, tx, d); err != nil {
		return fmt.Errorf("could not create delivery: %w", err)
	}

	// 4. commit
	return tx.Commit()
}

func (uc *UseCase) GetDeliveryByID(ctx context.Context, deliveryId uuid.UUID) (*delivery.Delivery, error) {
	return uc.repo.GetByID(deliveryId)
}

func (uc *UseCase) UpdateDelivery(ctx context.Context, deliveryID uuid.UUID, column string, value any) error {
	return uc.repo.Update(ctx, deliveryID, column, value)
}

func (uc *UseCase) ListDeliveries(ctx context.Context) ([]*delivery.Delivery, error) {
	return uc.repo.List()
}

func (uc *UseCase) DeleteDelivery(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}
