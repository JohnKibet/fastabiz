package main

import (
	"log"
	"net/http"
	"os"

	"logistics-backend/handlers"
	"logistics-backend/internal/repository/postgres"
	"logistics-backend/internal/router"
	orderUsecase "logistics-backend/internal/usecase/order"
	userUsercase "logistics-backend/internal/usecase/user"

	_ "logistics-backend/docs"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// @title Logistics API
// @version 1.0
// @description This is the API for logistics operations.
// @host localhost:8080
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

	// Set up usecase
	uUsecase := userUsercase.NewUseCase(userRepo)
	oUsecase := orderUsecase.NewUseCase(orderRepo)

	// Set up Handlers
	userHandler := handlers.NewUserHandler(uUsecase)
	orderHandler := handlers.NewOrderHandler(oUsecase)

	// Start server
	r := router.NewRouter(userHandler, orderHandler)

	log.Println("Server starting at :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal("could not start server at: %v", err)
	}
}
