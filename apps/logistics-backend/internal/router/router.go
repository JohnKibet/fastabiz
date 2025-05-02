package router

import (
	"net/http"

	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"logistics-backend/handlers"
)

func NewRouter(u *handlers.UserHandler, o *handlers.OrderHandler) http.Handler {
	r := chi.NewRouter()

	// Enable Cors
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Basic middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Swagger documentation endpoint
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// User routes
	r.Route("/users", func(r chi.Router) {
		r.Post("/create", u.CreateUser)
		r.Get("/all_users", u.ListUsers)
		r.Get("/id/{id}", u.GetUserByID)
		r.Get("/email/{email}", u.GetUserByEmail)
	})

	// Order routes
	r.Route("/orders", func(r chi.Router) {
		r.Post("/create", o.CreateOrder)
		r.Get("/id/{id}", o.GetOrderByID)
		r.Get("/customer/{customer_id}", o.GetOrderByCustomer)
		r.Put("/{order_id}/status", o.UpdateOrderStatus)
	})

	return r

}
