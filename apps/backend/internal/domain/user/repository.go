package user

import (
	"context"

	"github.com/google/uuid"
)

// User = actual onboarded account in the system.
type Repository interface {
	Create(ctx context.Context, user *User) error                                         // POST
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)                             // GET
	GetByEmail(ctx context.Context, email string) (*User, error)                          // GET
	List(ctx context.Context) ([]*User, error)                                            // GET
	GetAllCustomers(ctx context.Context) ([]AllCustomers, error)                          // GET
	UpdateColum(ctx context.Context, userID uuid.UUID, column string, value any) error    // PATCH update one specified column at a time
	UpdateUserStatus(ctx context.Context, userID uuid.UUID, status UserStatus) error      // PATCH method to update user status
	UpdateDriverProfile(ctx context.Context, id uuid.UUID, phone string) error            // PUT update user(driver) phone number
	UpdateUserProfile(ctx context.Context, id uuid.UUID, phone, email, name string) error // PATCH method to update user profile - name, email & phone number
	Delete(ctx context.Context, id uuid.UUID) error                                       // DELETE
}
