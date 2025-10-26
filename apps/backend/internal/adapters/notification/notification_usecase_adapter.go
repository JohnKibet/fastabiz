package notificationadapter

import (
	notificationusecase "backend/internal/usecase/notification"
)

type UseCaseAdapter struct {
	UseCase *notificationusecase.UseCase
}
