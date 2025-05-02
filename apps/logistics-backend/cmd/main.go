package cmd

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln("Database connection failed:", err)
	}
	defer db.Close()

	// Use db to init repos:
	// userRepo := postgres.NewUserRepository(db)
	// orderRepo := postgres.NewOrderRepository(db)
}
