package store

import (
	"backend/internal/domain/store"
	"backend/internal/usecase/common"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type UseCase struct {
	repo      store.Repository
	txManager common.TxManager
}

func NewUseCase(repo store.Repository, txm common.TxManager) *UseCase {
	return &UseCase{repo: repo, txManager: txm}
}

func (uc *UseCase) CreateStore(ctx context.Context, s *store.Store) error {
	s.NameNormalized = NormalizeName(s.Name)

	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		return uc.repo.Create(txCtx, s)
	})
}

func (uc *UseCase) GetStoreByID(ctx context.Context, id uuid.UUID) (*store.Store, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetBasicStoreByID(ctx context.Context, storeID uuid.UUID) (*store.StoreBasic, error) {
	return uc.repo.GetBasicByID(ctx, storeID)
}

func (uc *UseCase) GetStoreSummary(ctx context.Context, storeID uuid.UUID) (*store.StoreSummary, error) {
	return uc.repo.GetStoreSummary(ctx, storeID)
}

func (uc *UseCase) UpdateStore(ctx context.Context, storeID uuid.UUID, ownerID uuid.UUID, req *store.UpdateStoreRequest) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {

		owned, err := uc.repo.IsOwnedBy(txCtx, storeID, ownerID)
		if err != nil {
			return fmt.Errorf("ownership check failed: %w", err)
		}
		if !owned {
			return fmt.Errorf("store not owned by user")
		}

		if err := uc.repo.UpdateStoreDetails(txCtx, storeID, req.Name, req.Logo, req.Location); err != nil {
			return fmt.Errorf("update store failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) ListAllStores(ctx context.Context) ([]*store.Store, error) {
	return uc.repo.ListStores(ctx)
}

func (uc *UseCase) ListOwnerStores(ctx context.Context, ownerID uuid.UUID) ([]*store.Store, error) {
	return uc.repo.ListStoresByOwner(ctx, ownerID)
}

func (uc *UseCase) ListStoresPaged(ctx context.Context, filter store.StoreFilter) ([]*store.StoreSummary, error) {
	return uc.repo.ListStoresPaged(ctx, filter)
}

func (uc *UseCase) DeleteStore(ctx context.Context, storeID uuid.UUID, ownerID uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {

		owned, err := uc.repo.IsOwnedBy(txCtx, storeID, ownerID)
		if err != nil {
			return fmt.Errorf("ownership check failed: %w", err)
		}
		if !owned {
			return fmt.Errorf("store not owned by user")
		}

		if err := uc.repo.Delete(txCtx, storeID); err != nil {
			return fmt.Errorf("delete store failed: %w", err)
		}

		return nil
	})
}

func (uc *UseCase) StoreExists(ctx context.Context, storeID uuid.UUID) (bool, error) {
	return uc.repo.Exists(ctx, storeID)
}

func (uc *UseCase) IsStoreOwnedBy(ctx context.Context, storeID uuid.UUID, ownerID uuid.UUID) (bool, error) {
	return uc.repo.IsOwnedBy(ctx, storeID, ownerID)
}

func NormalizeName(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
