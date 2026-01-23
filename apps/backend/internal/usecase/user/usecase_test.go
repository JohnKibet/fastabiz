package user

import (
	"backend/internal/domain/driver"
	"backend/internal/domain/notification"
	"backend/internal/domain/user"
	"backend/internal/usecase/common"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type fakeTxManager struct {
	called bool
	usedTx bool
}

func (f *fakeTxManager) Do(ctx context.Context, fn func(context.Context) error) error {
	f.called = true
	txCtx := common.MarkTx(ctx)
	return fn(txCtx)
}

type fakeUserRepo struct {
	createdUser *user.User
	createErr   error
	usedTx bool
}

func (f *fakeUserRepo) Create(ctx context.Context, u *user.User) error {
	f.createdUser = u
	f.usedTx = common.IsTransactional(ctx)
	return f.createErr
}

type fakeDriverRepo struct {
	registered bool
	err        error
}

func (f *fakeDriverRepo) RegisterDriver(ctx context.Context, d *driver.Driver) error {
	f.registered = true
	return f.err
}

type fakeNotificationRepo struct{}

func (f *fakeNotificationRepo) Create(ctx context.Context, n *notification.Notification) error {
	return nil
}

func TestRegisterUser_UsesTransaction(t *testing.T) {
	txm := &fakeTxManager{}
	repo := &fakeUserRepo{}
	drv := &fakeDriverRepo{}
	notf := &fakeNotificationRepo{}

	uc := NewUseCase(repo, drv, txm, notf)

	u := &user.User{
		FullName:     "John Doe",
		Email:        "john@example.com",
		PasswordHash: "plain-password",
		Role:         user.Customer,
	}

	err := uc.RegisterUser(context.Background(), u)

	require.NoError(t, err)
	require.True(t, txm.called, "expected transaction manager to be used")
	require.True(t, repo.usedTx, "Create must run inside transaction")
}

func TestRegisterUser_HashesPassword(t *testing.T) {
	repo := &fakeUserRepo{}
	uc := NewUseCase(repo, &fakeDriverRepo{}, &fakeTxManager{}, &fakeNotificationRepo{})

	u := &user.User{
		PasswordHash: "plain-password",
		Role:         user.Customer,
	}

	err := uc.RegisterUser(context.Background(), u)

	require.NoError(t, err)
	require.NotEqual(t, "plain-password", repo.createdUser.PasswordHash)

	err = bcrypt.CompareHashAndPassword(
		[]byte(repo.createdUser.PasswordHash),
		[]byte("plain-password"),
	)
	require.NoError(t, err)
}

func TestRegisterUser_DriverCreatesDriverRecord(t *testing.T) {
	drv := &fakeDriverRepo{}
	uc := NewUseCase(&fakeUserRepo{}, drv, &fakeTxManager{}, &fakeNotificationRepo{})

	u := &user.User{
		Role: user.Driver,
	}

	err := uc.RegisterUser(context.Background(), u)

	require.NoError(t, err)
	require.True(t, drv.registered)
}

func TestRegisterUser_NonDriverDoesNotCreateDriver(t *testing.T) {
	drv := &fakeDriverRepo{}
	uc := NewUseCase(&fakeUserRepo{}, drv, &fakeTxManager{}, &fakeNotificationRepo{})

	u := &user.User{Role: user.Customer}

	_ = uc.RegisterUser(context.Background(), u)

	require.False(t, drv.registered)
}

func TestRegisterUser_RepoErrorReturned(t *testing.T) {
	repo := &fakeUserRepo{createErr: errors.New("db down")}
	uc := NewUseCase(repo, &fakeDriverRepo{}, &fakeTxManager{}, &fakeNotificationRepo{})

	err := uc.RegisterUser(context.Background(), &user.User{})

	require.Error(t, err)
	require.Contains(t, err.Error(), "could not create user")
}

// --- unused methods (minimal stubs) ---
func (f *fakeUserRepo) GetByID(context.Context, uuid.UUID) (*user.User, error) {
	return nil, nil
}
func (f *fakeUserRepo) GetByEmail(context.Context, string) (*user.User, error) {
	return nil, nil
}
func (f *fakeUserRepo) List(context.Context) ([]*user.User, error) {
	return nil, nil
}
func (f *fakeUserRepo) GetAllCustomers(context.Context) ([]user.AllCustomers, error) {
	return nil, nil
}
func (f *fakeUserRepo) UpdateColum(context.Context, uuid.UUID, string, any) error {
	return nil
}
func (f *fakeUserRepo) UpdateUserStatus(context.Context, uuid.UUID, user.UserStatus) error {
	return nil
}
func (f *fakeUserRepo) UpdateDriverProfile(context.Context, uuid.UUID, string) error {
	return nil
}
func (f *fakeUserRepo) UpdateUserProfile(context.Context, uuid.UUID, string, string, string) error {
	return nil
}
func (f *fakeUserRepo) Delete(context.Context, uuid.UUID) error {
	return nil
}