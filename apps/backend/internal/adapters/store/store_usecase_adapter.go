package store

import storeusecase "backend/internal/usecase/store"

type UseCaseAdapter struct {
	UseCase *storeusecase.UseCase
}
