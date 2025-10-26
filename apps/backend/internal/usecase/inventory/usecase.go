package inventory

import (
	"context"
	"fmt"
	domain "backend/internal/domain/inventory"
	"backend/internal/domain/notification"
	"backend/internal/usecase/common"

	"github.com/google/uuid"
)

type UseCase struct {
	repo      domain.Repository
	txManager common.TxManager
	notfRepo  domain.NotificationReader
	storeRepo domain.StoreReader
}

func NewUseCase(repo domain.Repository, txm common.TxManager, notf domain.NotificationReader, str domain.StoreReader) *UseCase {
	return &UseCase{repo: repo, txManager: txm, notfRepo: notf, storeRepo: str}
}

func (uc *UseCase) CreateInventory(ctx context.Context, i *domain.Inventory) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Create(txCtx, i); err != nil {
			return fmt.Errorf("could not create inventory: %w", err)
		}

		// Fetch store to get AdminID/ OwnerID
		store, err := uc.storeRepo.GetByID(txCtx, i.StoreID)
		if err != nil {
			return fmt.Errorf("could not fetch store: %w", err)
		}

		// After successful creation, fire notification (async)
		go func() {
			msg := fmt.Sprintf("✅ New inventory '%s' has been added with stock %d.", i.Category, i.Stock)
			_ = uc.notify(ctx, store.OwnerID, msg)

			// Optional: immediately alert if created with low stock
			if i.Stock <= 5 {
				lowMsg := fmt.Sprintf("⚠️ Inventory '%s' was created with low stock (%d).", i.Category, i.Stock)
				_ = uc.notify(ctx, store.OwnerID, lowMsg)
			}
		}()

		return nil
	})
}

func (uc *UseCase) GetInventory(ctx context.Context, id uuid.UUID) (*domain.Inventory, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetInventoryByName(ctx context.Context, name string) (*domain.Inventory, error) {
	return uc.repo.GetByName(ctx, name)
}

func (uc *UseCase) UpdateInventory(ctx context.Context, inventoryId uuid.UUID, column string, value any) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		// 1. Fetch inventory to get Category
		inv, err := uc.repo.GetByID(txCtx, inventoryId)
		if err != nil {
			return fmt.Errorf("could not fetch inventory: %w", err)
		}

		// 2. Fetch store to get AdminID/ OwnerID
		store, err := uc.storeRepo.GetByID(txCtx, inv.StoreID)
		if err != nil {
			return fmt.Errorf("could not fetch store: %w", err)
		}

		// 3. Update column
		if err := uc.repo.UpdateColumn(txCtx, inventoryId, column, value); err != nil {
			return fmt.Errorf("update inventory failed: %w", err)
		}

		// 4. Fire notification async (after commit)
		go func() {
			msg := fmt.Sprintf("ℹ️ Inventory %s updated: column '%s' changed.", inv.Category, column)
			_ = uc.notify(ctx, store.OwnerID, msg) // you can use AdminID if available
		}()

		return nil
	})
}

func (uc *UseCase) List(ctx context.Context, limit, offset int) ([]*domain.Inventory, error) {
	return uc.repo.List(ctx, limit, offset)
}

func (uc *UseCase) GetByCategory(ctx context.Context, category string) ([]*domain.Inventory, error) {
	return uc.repo.GetByCategory(ctx, category)
}

func (uc *UseCase) GetByStore(ctx context.Context, storeID uuid.UUID) ([]*domain.Inventory, error) {
	return uc.repo.GetByStoreID(ctx, storeID)
}

func (uc *UseCase) ListCategories(ctx context.Context) ([]string, error) {
	return uc.repo.ListCategories(ctx)
}

func (uc *UseCase) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		// 1. Fetch inventory to get Category
		inv, err := uc.repo.GetByID(txCtx, id)
		if err != nil {
			return fmt.Errorf("could not fetch inventory: %w", err)
		}

		// 2. Fetch store to get AdminID/ OwnerID
		store, err := uc.storeRepo.GetByID(txCtx, inv.StoreID)
		if err != nil {
			return fmt.Errorf("could not fetch store: %w", err)
		}

		// 3. Delete
		if err := uc.repo.Delete(txCtx, id); err != nil {
			return fmt.Errorf("delete inventory failed: %w", err)
		}

		// 4. Fire notification async
		go func() {
			msg := fmt.Sprintf("🗑️ Inventory '%s' has been deleted.", inv.Category)
			_ = uc.notify(ctx, store.OwnerID, msg)
		}()

		return nil
	})
}

func (uc *UseCase) GetAllInventories(ctx context.Context) ([]domain.AllInventory, error) {
	return uc.repo.GetAllInventories(ctx)
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
