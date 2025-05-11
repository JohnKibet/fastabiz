package main

import (
	"log"
	"net/http"
	"os"

	"logistics-backend/handlers"
	"logistics-backend/internal/repository/postgres"
	"logistics-backend/internal/router"
	deliveryUsecase "logistics-backend/internal/usecase/delivery"
	driverUsecase "logistics-backend/internal/usecase/driver"
	feedbackUsecase "logistics-backend/internal/usecase/feedback"
	notificationUsecase "logistics-backend/internal/usecase/notification"
	orderUsecase "logistics-backend/internal/usecase/order"
	paymentUsecase "logistics-backend/internal/usecase/payment"
	userUsecase "logistics-backend/internal/usecase/user"

	_ "logistics-backend/docs"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// @title Logistics API
// @version 1.0
// @description This is the API for logistics operations.
// @host 192.168.1.18:8080
// @BasePath /
// @schemes http
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// Set up repositories
	userRepo := postgres.NewUserRepository(db)
	orderRepo := postgres.NewOrderRepository(db)
	driverRepo := postgres.NewDriverRepository(db)
	deliveryRepo := postgres.NewDeliveryRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	feedbackRepo := postgres.NewFeedbackRepository(db)
	notificationRepo := postgres.NewNotificationRepository(db)

	// Set up usecase
	uUsecase := userUsecase.NewUseCase(userRepo)
	oUsecase := orderUsecase.NewUseCase(orderRepo)
	dUsecase := driverUsecase.NewUseCase(driverRepo)
	eUsecase := deliveryUsecase.NewUseCase(deliveryRepo)
	pUsecase := paymentUsecase.NewUseCase(paymentRepo)
	fUsecase := feedbackUsecase.NewUseCase(feedbackRepo)
	nUsecase := notificationUsecase.NewUseCase(notificationRepo)

	// Set up Handlers
	userHandler := handlers.NewUserHandler(uUsecase)
	orderHandler := handlers.NewOrderHandler(oUsecase)
	driverHandler := handlers.NewDriverHandler(dUsecase)
	deliveryHandler := handlers.NewDeliveryHandler(eUsecase)
	paymentHandler := handlers.NewPaymentHandler(pUsecase)
	feedbackHandler := handlers.NewFeedbackHandler(fUsecase)
	notificationHandler := handlers.NewNotificationHandler(nUsecase)

	// Start server
	r := router.NewRouter(userHandler, orderHandler, driverHandler, deliveryHandler, paymentHandler, feedbackHandler, notificationHandler)

	log.Println("Server starting at :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal("could not start server at: %v", err)
	}
}
