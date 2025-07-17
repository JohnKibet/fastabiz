package user

import (
	"context"
	"fmt"
	"logistics-backend/internal/domain/driver"
	domain "logistics-backend/internal/domain/user"
	"time"

	"github.com/google/uuid"
)

type UseCase struct {
	repo  domain.Repository
	dRepo domain.CreateDriver
}

func NewUseCase(repo domain.Repository, dRepo domain.CreateDriver) *UseCase {
	return &UseCase{repo: repo, dRepo: dRepo}
}

func (uc *UseCase) RegisterUser(ctx context.Context, u *domain.User) error {
	user, err := uc.repo.GetByID(u.ID)
	if err != nil {
		return err
	}
	if user.Role == "driver" {
		driver := &driver.Driver{
			ID:              user.ID,
			FullName:        user.FullName,
			Email:           user.Email,
			VehicleInfo:     "not set",
			CurrentLocation: "not set",
			Available:       true,
			CreatedAt:       time.Now().String(),
		}
		err := uc.dRepo.RegisterDriver(ctx, driver)
		if err != nil {
			return fmt.Errorf("could not register driver: %w", err)
		}
	}

	return uc.repo.Create(u)
}

// PATCH method for drivers users to update details
func (uc *UseCase) UpdateDriverProfile(ctx context.Context, id uuid.UUID, req *domain.UpdateDriverUserProfileRequest) error {
	return uc.repo.UpdateProfile(ctx, id, req.Phone)
}

// PATCH method for user details
func (uc *UseCase) UpdateUser(ctx context.Context, userID uuid.UUID, req *domain.UpdateUserRequest) error {
	return uc.repo.UpdateColum(ctx, userID, req.Column, req.Value)
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

func (uc *UseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Delete(ctx, id)
}
