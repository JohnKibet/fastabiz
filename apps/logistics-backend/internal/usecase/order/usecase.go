package order

import (
	"context"

	domain "logistics-backend/internal/domain/order"

	"github.com/google/uuid"
)

type UseCase struct {
	repo domain.Repository
}

func NewUseCase(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) CreateOrder(ctx context.Context, o *domain.Order) error {

	o.OrderStatus = domain.Pending
	return uc.repo.Create(o)
}

func (uc *UseCase) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	return uc.repo.GetByID(id)
}

func (uc *UseCase) GetOrderByCustomer(ctx context.Context, customerID uuid.UUID) ([]*domain.Order, error) {
	return uc.repo.ListByCustomer(customerID)
}

func (uc *UseCase) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error {
	return uc.repo.UpdateStatus(orderID, status)
}

func (uc *UseCase) ListOrders(ctx context.Context) ([]*domain.Order, error) {
	return uc.repo.List()
}
