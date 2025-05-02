package user

import (
	"context"
	domain "logistics-backend/internal/domain/user"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

type UseCase struct {
	repo domain.Repository
}

func NewUseCase(repo domain.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) RegisterUser(ctx context.Context, u *domain.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashed)
	return uc.repo.Create(u)
}

func (uc *UseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return uc.repo.GetByID(id)
}

func (uc *UseCase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return uc.repo.GetByEmail(email)
}

func (uc *UseCase) ListUsers(ctx context.Context) ([]*domain.User, error) {
	return uc.repo.List()
}
