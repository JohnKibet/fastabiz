package order

import (
	"context"
	"fmt"
	"logistics-backend/internal/domain/order"

	"github.com/google/uuid"
)

type UseCase struct {
	repo     order.Repository
	invLogic order.InventoryReader
}

func NewUseCase(repo order.Repository, inv order.InventoryReader) *UseCase {
	return &UseCase{repo: repo, invLogic: inv}
}

func (uc *UseCase) CreateOrder(ctx context.Context, o *order.Order) error {
	if o.Quantity <= 0 {
		return order.ErrorInvalidQuantity
	}

	inv, err := uc.invLogic.GetByID(ctx, o.InventoryID)
	if err != nil {
		return fmt.Errorf("could not fetch inventory: %w", err)
	}

	if inv.Stock < o.Quantity {
		return order.ErrorOutOfStock
	}

	// updating inventory table
	newStock := inv.Stock - o.Quantity
	if err := uc.invLogic.UpdateInventory(ctx, inv.ID, "stock", newStock); err != nil {
		return fmt.Errorf("could not update inventory stock: %w", err)
	}

	o.OrderStatus = order.Pending
	return uc.repo.Create(o)

	// TODO: goroutine
	// fire notification for order created
	// fire notification for restocking reminder
}

func (uc *UseCase) GetOrder(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	return uc.repo.GetByID(id)
}

func (uc *UseCase) GetOrderByCustomer(ctx context.Context, customerID uuid.UUID) ([]*order.Order, error) {
	return uc.repo.ListByCustomer(customerID)
}

func (uc *UseCase) UpdateOrder(ctx context.Context, orderID uuid.UUID, column string, value any) error {
	return uc.repo.UpdateColumn(ctx, orderID, column, value)
}

func (uc *UseCase) ListOrders(ctx context.Context) ([]*order.Order, error) {
	return uc.repo.List()
}
