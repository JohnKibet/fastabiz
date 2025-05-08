package payment

import (
	"context"
	domain "logistics-backend/internal/domain/payment"

	"github.com/google/uuid"
)

type UseCase struct {
	repo domain.Repository
}

func NewUseCase(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) CreatePayment(ctx context.Context, p *domain.Payment) error {
	return uc.repo.Create(p)
}

func (uc *UseCase) GetPaymentByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	return uc.repo.GetByID(id)
}

func (uc *UseCase) GetPaymentByOrderID(ctx context.Context, order_id uuid.UUID) (*domain.Payment, error) {
	return uc.repo.GetByOrder(order_id)
}

func (uc *UseCase) ListPayments(ctx context.Context) ([]*domain.Payment, error) {
	return uc.repo.List()
}
