package main

import (
	"log"
	"net/http"
	"os"

	"backend/handlers"
	deliveryadapter "backend/internal/adapters/delivery"
	driveradapter "backend/internal/adapters/driver"
	inventoryadapter "backend/internal/adapters/inventory"
	notificationadapter "backend/internal/adapters/notification"
	orderadapter "backend/internal/adapters/order"
	useradapter "backend/internal/adapters/user"
	"backend/internal/repository/postgres"
	"backend/internal/router"
	deliveryUsecase "backend/internal/usecase/delivery"
	driverUsecase "backend/internal/usecase/driver"
	feedbackUsecase "backend/internal/usecase/feedback"
	inventoryUsecase "backend/internal/usecase/inventory"
	inviteUsecase "backend/internal/usecase/invite"
	notificationUsecase "backend/internal/usecase/notification"
	orderUsecase "backend/internal/usecase/order"
	paymentUsecase "backend/internal/usecase/payment"
	storeUsecase "backend/internal/usecase/store"
	userUsecase "backend/internal/usecase/user"

	"backend/internal/application"

	_ "backend/docs"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"time"
)

// @title FastaBiz API
// @version 1.0
// @description This is the API for fastabiz operations.
// @host localhost:8000
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey JWT
// @in header
// @name Authorization
func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL not set")
	}

	publicApiBaseUrl := os.Getenv("PUBLIC_API_BASE_URL")
	if publicApiBaseUrl == "" {
		log.Fatal("PUBLIC_API_BASE_URL not set")
	}

	db := waitForPostgres(dbUrl, 10, 5*time.Second)

	txm := application.NewTxManager(db)

	// Set up repositories
	userRepo := postgres.NewUserRepository(db)
	orderRepo := postgres.NewOrderRepository(db)
	driverRepo := postgres.NewDriverRepository(db)
	deliveryRepo := postgres.NewDeliveryRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	feedbackRepo := postgres.NewFeedbackRepository(db)
	notificationRepo := postgres.NewNotificationRepository(db)
	inventoryRepo := postgres.NewInventoryRespository(db)
	inviteRepo := postgres.NewInviteRepository(db)
	storeRepo := postgres.NewStoreRepository(db)

	// Set up usecase
	// Individual
	inviteUC := inviteUsecase.NewUseCase(inviteRepo, txm)
	driverUC := driverUsecase.NewUseCase(driverRepo, txm, notificationRepo)
	userUC := userUsecase.NewUseCase(userRepo, driverUC, txm, notificationRepo)
	inventoryUC := inventoryUsecase.NewUseCase(inventoryRepo, txm, notificationRepo, storeRepo)
	orderUC := orderUsecase.NewUseCase(orderRepo, &inventoryadapter.UseCaseAdapter{UseCase: inventoryUC}, &useradapter.UseCaseAdapter{UseCase: userUC}, txm, notificationRepo, storeRepo)
	deliveryUC := deliveryUsecase.NewUseCase(deliveryRepo, &orderadapter.UseCaseAdapter{UseCase: orderUC}, &driveradapter.UseCaseAdapter{UseCase: driverUC}, txm, notificationRepo)
	notificationUC := notificationUsecase.NewUseCase(notificationRepo, txm)
	storeUC := storeUsecase.NewUseCase(storeRepo, txm)

	// Combined cross-domain service
	orderService := application.NewOrderService(
		&useradapter.UseCaseAdapter{UseCase: userUC},
		&orderadapter.UseCaseAdapter{UseCase: orderUC},
		&driveradapter.UseCaseAdapter{UseCase: driverUC},
		&deliveryadapter.UseCaseAdapter{UseCase: deliveryUC},
		&inventoryadapter.UseCaseAdapter{UseCase: inventoryUC},
		&notificationadapter.UseCaseAdapter{UseCase: notificationUC},
	)

	// Other usecases
	paymentUC := paymentUsecase.NewUseCase(paymentRepo, txm)
	feedbackUC := feedbackUsecase.NewUseCase(feedbackRepo, txm)

	// Set up Handlers
	userHandler := handlers.NewUserHandler(orderService)
	orderHandler := handlers.NewOrderHandler(orderService)
	driverHandler := handlers.NewDriverHandler(orderService)
	deliveryHandler := handlers.NewDeliveryHandler(orderService)
	paymentHandler := handlers.NewPaymentHandler(paymentUC)
	feedbackHandler := handlers.NewFeedbackHandler(feedbackUC)
	notificationHandler := handlers.NewNotificationHandler(orderService)
	inviteHandler := handlers.NewInviteHandler(inviteUC)
	inventoryHandler := handlers.NewInventoryHandler(orderService)
	storeHandler := handlers.NewStoreHandler(storeUC)

	// Start server
	r := router.NewRouter(
		userHandler,
		orderHandler,
		driverHandler,
		deliveryHandler,
		paymentHandler,
		feedbackHandler,
		notificationHandler,
		inventoryHandler,
		publicApiBaseUrl,
		inviteHandler,
		storeHandler,
		db,
	)

	log.Println("Server starting at :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Printf("could not start server at: %v", err)
	}
}

// helper to wait for postgres
func waitForPostgres(dsn string, maxRetries int, delay time.Duration) *sqlx.DB {
	var db *sqlx.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			// ping to confirm
			if pingErr := db.Ping(); pingErr == nil {
				log.Println("✅ Connected to Postgres successfully")
				return db
			}
		}

		log.Printf("⏳ Waiting for Postgres... (%d/%d)", i+1, maxRetries)
		time.Sleep(delay)
	}

	log.Fatalf("❌ Could not connect to Postgres after %d attempts: %v", maxRetries, err)
	return nil
}
