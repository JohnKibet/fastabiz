package product

import (
	"backend/internal/domain/product"
	"backend/internal/usecase/common"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type UseCase struct {
	repo      product.Repository
	txManager common.TxManager
}

func NewUseCase(repo product.Repository, txm common.TxManager) *UseCase {
	return &UseCase{repo: repo, txManager: txm}
}

func (uc *UseCase) CreateProduct(ctx context.Context, p *product.Product) (err error) {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.Create(txCtx, p); err != nil {
			return fmt.Errorf("could not create product: %w", err)
		}
		return nil
	})
}

func (uc *UseCase) GetProductByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) UpdateProductDetails(ctx context.Context, req *product.UpdateProductDetailsRequest) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateDetails(txCtx, req.ProductID, req.Category, req.Name, req.Description); err != nil {
			return fmt.Errorf("update product details failed: %w", err)
		}
		return nil
	})
}

func (uc *UseCase) GetAllProducts(ctx context.Context) ([]product.Product, error) {
	return uc.repo.List(ctx)
}

func (uc *UseCase) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	return uc.repo.Delete(ctx, productID)
}

func (uc *UseCase) AddImage(ctx context.Context, productID uuid.UUID, url string, isPrimary bool) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.AddImage(txCtx, productID, url, isPrimary); err != nil {
			return fmt.Errorf("image upload failed: %w", err)
		}
		return nil
	})
}

func (uc *UseCase) DeleteImage(ctx context.Context, imageID uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.RemoveImage(txCtx, imageID); err != nil {
			return fmt.Errorf("delete image failed: %w", err)
		}
		return nil
	})
}

func (uc *UseCase) ReorderImages(ctx context.Context, productID uuid.UUID, imageIDs []uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.ReorderImages(txCtx, productID, imageIDs); err != nil {
			return fmt.Errorf("reorder images failed: %w", err)
		}
		return nil
	})
}

func (uc *UseCase) AddOptionName(ctx context.Context, productID uuid.UUID, name string) (uuid.UUID, error) {
	var optionID uuid.UUID
	err := uc.txManager.Do(ctx, func(txCtx context.Context) error {
		id, err := uc.repo.AddOption(txCtx, productID, name)
		if err != nil {
			return fmt.Errorf("add option name failed: %w", err)
		}
		optionID = id
		return nil
	})
	return optionID, err
}

func (uc *UseCase) DeleteOptionName(ctx context.Context, optionID uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.RemoveOption(txCtx, optionID); err != nil {
			return fmt.Errorf("delete option name failed: %w", err)
		}
		return nil
	})
}

func (uc *UseCase) AddOptionValue(ctx context.Context, optionID uuid.UUID, value string) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.AddOptionValue(txCtx, optionID, value); err != nil {
			return fmt.Errorf("add option value failed: %w", err)
		}
		return nil
	})
}

func (uc *UseCase) DeleteOptionValue(ctx context.Context, optionValueID uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.RemoveOptionValue(txCtx, optionValueID); err != nil {
			return fmt.Errorf("delete option value failed: %w", err)
		}
		return nil
	})
}

func (uc *UseCase) CreateVariant(ctx context.Context, variant *product.Variant) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.CreateVariant(txCtx, variant); err != nil {
			return fmt.Errorf("create new variant failed: %w", err)
		}
		return nil
	})
}

func (uc *UseCase) UpdateVariantStock(ctx context.Context, variantID uuid.UUID, stock int) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateVariantStock(txCtx, variantID, stock); err != nil {
			return fmt.Errorf("update variant stock failed: %w", err)
		}
		return nil
	})
}

func (uc *UseCase) UpdateVariantPrice(ctx context.Context, variantID uuid.UUID, price float64) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateVariantPrice(txCtx, variantID, price); err != nil {
			return fmt.Errorf("update variant price failed: %w", err)
		}
		return nil
	})
}

func (uc *UseCase) DeleteVariant(ctx context.Context, variantID uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.RemoveVariant(txCtx, variantID); err != nil {
			return fmt.Errorf("delete variant failed: %w", err)
		}
		return nil
	})
}
