package user

import (
	"backend/internal/domain/driver"
	"backend/internal/domain/notification"
	domain "backend/internal/domain/user"
	"backend/internal/usecase/common"
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/cridenour/go-postgis"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo      domain.Repository
	drvRepo   domain.DriverReader
	txManager common.TxManager
	notfRepo  domain.NotificationReader
}

func NewUseCase(repo domain.Repository, drvRepo domain.DriverReader, txm common.TxManager, notf domain.NotificationReader) *UseCase {
	return &UseCase{repo: repo, drvRepo: drvRepo, txManager: txm, notfRepo: notf}
}

func (uc *UseCase) RegisterUser(ctx context.Context, u *domain.User) error {

	return uc.txManager.Do(ctx, func(txCtx context.Context) error {

		// 1. hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		u.PasswordHash = string(hashedPassword)

		// 2. insert user to DB
		if err := uc.repo.Create(txCtx, u); err != nil {
			return fmt.Errorf("could not create user: %w", err)
		}

		// 3. if role is driver, insert into drivers table
		if u.Role == "driver" {
			driver := &driver.Driver{
				ID:          u.ID,
				FullName:    u.FullName,
				Email:       u.Email,
				VehicleInfo: "not set",
				CurrentLocation: postgis.PointS{
					SRID: 4326,
					X:    36.8219, // longitude
					Y:    -1.2921, // latitude
				},
				Available: true,
				CreatedAt: time.Now(),
			}
			if err := uc.drvRepo.RegisterDriver(txCtx, driver); err != nil {
				return fmt.Errorf("could not register driver: %w", err)
			}
		}

		// async notification after commit
		go func() {
			msg := fmt.Sprintf("✅ New user '%s' has been registered with role '%s'.", u.FullName, u.Role)
			_ = uc.notify(ctx, u.ID, msg)
		}()

		return nil
	})
}

// PATCH method for user(driver) to update phone no. after registration
func (uc *UseCase) UpdateDriverProfile(ctx context.Context, id uuid.UUID, req *domain.UpdateDriverUserProfileRequest) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		// fetch user for notification
		user, err := uc.repo.GetByID(txCtx, id)
		if err != nil {
			return fmt.Errorf("could not fetch user: %w", err)
		}

		if err := uc.repo.UpdateDriverProfile(txCtx, id, req.Phone); err != nil {
			return fmt.Errorf("update user profile failed: %w", err)
		}

		go func() {
			msg := "ℹ️ Your profile was updated successfully."
			_ = uc.notify(ctx, user.ID, msg)
		}()

		return nil
	})
}

// PATCH method for users to update profile details - phone, email & name
func (uc *UseCase) UpdateUserProfile(ctx context.Context, id uuid.UUID, req *domain.UpdateUserProfileRequest) error {
	if err := validateUserUpdate(req); err != nil {
		return err // return domain error directly
	}

	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		// fetch user for notification
		user, err := uc.repo.GetByID(txCtx, id)
		if err != nil {
			return fmt.Errorf("could not fetch user: %w: %w", domain.ErrDB, err)
		}

		if err := uc.repo.UpdateUserProfile(txCtx, id, req.Phone, req.Email, req.FullName); err != nil {
			return fmt.Errorf("update user profile failed: %w", err)
		}

		go func() {
			msg := "ℹ️ Your profile was updated successfully."
			_ = uc.notify(ctx, user.ID, msg)
		}()

		return nil
	})
}

// PATCH method for user details
func (uc *UseCase) UpdateUser(ctx context.Context, userID uuid.UUID, req *domain.UpdateUserRequest) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		user, err := uc.repo.GetByID(txCtx, userID)
		if err != nil {
			return fmt.Errorf("could not fetch user: %w", err)
		}

		if err := uc.repo.UpdateColum(txCtx, userID, req.Column, req.Value); err != nil {
			return fmt.Errorf("update user failed: %w", err)
		}

		go func() {
			msg := fmt.Sprintf("ℹ️ Your account column '%s' was updated.", req.Column)
			_ = uc.notify(ctx, user.ID, msg)
		}()

		return nil
	})
}

func (uc *UseCase) ChangePassword(ctx context.Context, userID uuid.UUID, req *domain.ChangePasswordRequest) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		user, err := uc.repo.GetByID(txCtx, userID)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		// 1. verify old password
		if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)) != nil {
			return fmt.Errorf("current password incorrect")
		}

		// 2. hash new password
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hashing failed: %w", err)
		}

		// 3. save
		if err := uc.repo.UpdateColum(txCtx, userID, "password_hash", string(hashed)); err != nil {
			return fmt.Errorf("failed to update password: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return uc.repo.GetByEmail(ctx, email)
}

func (uc *UseCase) ListUsers(ctx context.Context) ([]*domain.User, error) {
	return uc.repo.List(ctx)
}

func (uc *UseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		user, err := uc.repo.GetByID(txCtx, id)
		if err != nil {
			return fmt.Errorf("could not fetch user: %w", err)
		}

		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete user failed: %w", err)
		}

		go func() {
			msg := fmt.Sprintf("🗑️ User '%s' has been deleted.", user.FullName)
			_ = uc.notify(ctx, user.ID, msg)
		}()

		return nil
	})
}

func (uc *UseCase) GetAllCustomers(ctx context.Context) ([]domain.AllCustomers, error) {
	return uc.repo.GetAllCustomers(ctx)
}

func (uc *UseCase) notify(ctx context.Context, userID uuid.UUID, message string) error {
	n := &notification.Notification{
		UserID:  userID,
		Message: message,
		Type:    notification.System,
		Status:  notification.Pending,
	}
	return uc.notfRepo.Create(ctx, n)
}

func validateUserUpdate(req *domain.UpdateUserProfileRequest) error {
	if req.Phone != "" {
		// Very simple Kenya phone format example
		matched, _ := regexp.MatchString(`^\+254\s?\d{3}\s?\d{3}\s?\d{3}$`, req.Phone)
		if !matched {
			return domain.ErrInvalidPhone
		}
	}

	if req.Email != "" {
		if !strings.Contains(req.Email, "@") {
			return domain.ErrInvalidEmail
		}
	}

	if req.FullName != "" && len(req.FullName) < 2 {
		return domain.ErrInvalidName
	}

	return nil
}
