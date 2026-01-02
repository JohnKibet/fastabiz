package product

import productusecase "backend/internal/usecase/product"

type UseCaseAdapter struct {
	UseCase *productusecase.UseCase
}
