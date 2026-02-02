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
		return uc.repo.Create(txCtx, p)
	})
}

func (uc *UseCase) GetProductByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	p, err := uc.repo.GetFullProductByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("usecase GetProductByID: %w", err)
	}

	if p.HasVariants && len(p.Variants) == 0 {
		return nil, fmt.Errorf("product %s has variants flag but no variants found", p.ID)
	}

	return p, nil
}

func (uc *UseCase) UpdateProductDetails(ctx context.Context, req *product.UpdateProductDetailsRequest) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateDetails(txCtx, req.ProductID, req.Category, req.Name, req.Description); err != nil {
			return fmt.Errorf("%w", err)
		}
		return nil
	})
}

func (uc *UseCase) GetAllProducts(ctx context.Context, storeID uuid.UUID) ([]product.ProductListItem, error) {
	return uc.repo.ListProductsByStore(ctx, storeID)
}

func (uc *UseCase) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
	return uc.repo.Delete(ctx, productID)
}

func (uc *UseCase) AddImage(ctx context.Context, productID uuid.UUID, images []product.Image) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		for _, image := range images {
			if err := uc.repo.AddImage(txCtx, productID, image.URL, image.IsPrimary); err != nil {
				return fmt.Errorf("%w", err)
			}
		}
		return nil
	})
}

func (uc *UseCase) DeleteImage(ctx context.Context, imageID uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.RemoveImage(txCtx, imageID); err != nil {
			return fmt.Errorf("%w", err)
		}
		return nil
	})
}

func (uc *UseCase) ReorderImages(ctx context.Context, productID uuid.UUID, imageIDs []uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.ReorderImages(txCtx, productID, imageIDs); err != nil {
			return fmt.Errorf("%w", err)
		}
		return nil
	})
}

func (uc *UseCase) AddOptionName(ctx context.Context, productID uuid.UUID, name string) (uuid.UUID, error) {
	var optionID uuid.UUID
	err := uc.txManager.Do(ctx, func(txCtx context.Context) error {
		id, err := uc.repo.AddOption(txCtx, productID, name)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		optionID = id
		return nil
	})
	return optionID, err
}

func (uc *UseCase) DeleteOptionName(ctx context.Context, optionID uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		used, err := uc.repo.IsOptionUsed(txCtx, optionID)
		if err != nil {
			return err
		}

		if used {
			return product.ErrOptionInUse
		}

		return uc.repo.RemoveOption(txCtx, optionID)
	})
}

func (uc *UseCase) AddOptionValue(ctx context.Context, productID uuid.UUID, optionID uuid.UUID, values []string) error {
	if len(values) == 0 {
		return product.ErrMissingOptValues
	}

	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		for _, value := range values {
			if err := uc.repo.AddOptionValue(txCtx, productID, optionID, value); err != nil {
				return fmt.Errorf("%w", err)
			}
		}
		return nil
	})
}

func (uc *UseCase) ListProductOptions(ctx context.Context, productID uuid.UUID) ([]product.Option, error) {
	options, err := uc.repo.ListOptionsByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	for i := range options {
		values, err := uc.repo.ListOptionValuesByOptionID(ctx, options[i].ID)
		if err != nil {
			return nil, err
		}
		options[i].Values = values
	}

	return options, nil
}

func (uc *UseCase) DeleteOptionValue(ctx context.Context, optionValueID uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		used, err := uc.repo.IsOptionValueUsed(txCtx, optionValueID)
		if err != nil {
			return err
		}

		if used {
			return product.ErrOptionValueInUse
		}

		return uc.repo.RemoveOptionValue(txCtx, optionValueID)
	})
}

func (uc *UseCase) CreateVariant(ctx context.Context, req product.CreateVariantRequest) (*product.VariantWithOptions, error) {
	var variantWithOptions *product.VariantWithOptions

	err := uc.txManager.Do(ctx, func(txCtx context.Context) error {

		// 1. Resolve option name → value ID
		var optionvalueIDs []uuid.UUID

		for optName, optValue := range req.Options {
			optionID, err := uc.repo.GetOptionIDByName(txCtx, req.ProductID, optName)
			if err != nil {
				return product.ErrOptionNotFound
			}

			valueID, err := uc.repo.GetOptionValueID(txCtx, optionID, optValue)
			if err != nil {
				return product.ErrOptionValueNotFound
			}

			optionvalueIDs = append(optionvalueIDs, valueID)
		}

		// 2. Create variant
		variant := &product.Variant{
			ProductID: req.ProductID,
			SKU:       req.SKU,
			Price:     req.Price,
			Stock:     req.Stock,
			ImageURL:  req.ImageURL,
		}

		if err := uc.repo.CreateVariant(txCtx, variant); err != nil {
			return fmt.Errorf("%w", err)
		}

		// 3. Associate option values with variant
		for _, valueID := range optionvalueIDs {
			if err := uc.repo.AddVariantOptionValue(txCtx, variant.ID, valueID); err != nil {
				return fmt.Errorf("%w", err)
			}
		}

		// 4. Prepare response DTO
		variantWithOptions = &product.VariantWithOptions{
			ID:        variant.ID,
			ProductID: variant.ProductID,
			SKU:       variant.SKU,
			Price:     variant.Price,
			Stock:     variant.Stock,
			ImageURL:  variant.ImageURL,
			Options:   req.Options, // keep names for API
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return variantWithOptions, nil
}

func (uc *UseCase) UpdateVariantStock(ctx context.Context, variantID uuid.UUID, stock int) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateVariantStock(txCtx, variantID, stock); err != nil {
			return fmt.Errorf("%w", err)
		}
		return nil
	})
}

func (uc *UseCase) UpdateVariantPrice(ctx context.Context, variantID uuid.UUID, price float64) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateVariantPrice(txCtx, variantID, price); err != nil {
			return fmt.Errorf("%w", err)
		}
		return nil
	})
}

func (uc *UseCase) DeleteVariant(ctx context.Context, variantID uuid.UUID) error {
	return uc.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := uc.repo.RemoveVariant(txCtx, variantID); err != nil {
			return fmt.Errorf("%w", err)
		}
		return nil
	})
}
