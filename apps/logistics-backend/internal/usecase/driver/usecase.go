package driver

import (
	"context"
	domain "logistics-backend/internal/domain/driver"

	"github.com/google/uuid"
)

type UseCase struct {
	repo domain.Repository
}

func NewUseCase(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) RegisterDriver(ctx context.Context, d *domain.Driver) error {
	return uc.repo.Create(d)
}

func (uc *UseCase) GetDriverByID(ctx context.Context, id uuid.UUID) (*domain.Driver, error) {
	return uc.repo.GetByID(id)
}

func (uc *UseCase) GetDriverByEmail(ctx context.Context, email string) (*domain.Driver, error) {
	return uc.repo.GetByEmail(email)
}

func (uc *UseCase) ListDrivers(ctx context.Context) ([]*domain.Driver, error) {
	return uc.repo.List()
}
