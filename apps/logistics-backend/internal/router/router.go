package router

import (
	"net/http"

	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"logistics-backend/handlers"
)

func NewRouter(u *handlers.UserHandler, o *handlers.OrderHandler, d *handlers.DriverHandler, e *handlers.DeliveryHandler, p *handlers.PaymentHandler, f *handlers.FeedbackHandler, n *handlers.NotificationHandler) http.Handler {
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
		httpSwagger.URL("http://192.168.1.18:8080/swagger/doc.json"),
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
		r.Get("/all_orders", o.ListOrders)
		r.Get("/id/{id}", o.GetOrderByID)
		r.Get("/customer_id/{customer_id}", o.GetOrderByCustomer)
		r.Put("/{order_id}/status", o.UpdateOrderStatus)
	})

	// Driver routes
	r.Route("/drivers", func(r chi.Router) {
		r.Post("/create", d.CreateDriver)
		r.Get("/all_drivers", d.ListDrivers)
		r.Get("/id/{id}", d.GetDriverByID)
		r.Get("/email/{email}", d.GetDriverByEmail)
	})

	// Delivery routes
	r.Route("/deliveries", func(r chi.Router) {
		r.Post("/create", e.CreateDelivery)
		r.Get("/all_deliveries", e.ListDeliveries)
		r.Get("/id/{id}", e.GetDeliveryByID)
	})

	// Payment Routes
	r.Route("/payments", func(r chi.Router) {
		r.Post("/create", p.CreatePayment)
		r.Get("/all_payments", p.ListPayments)
		r.Get("/id/{id}", p.GetPaymentByID)
		r.Get("/order_id/{order_id}", p.GetPaymentByOrderID)
	})

	// Feedback Routes
	r.Route("/feedbacks", func(r chi.Router) {
		r.Post("/create", f.CreateFeedback)
		r.Get("/all_feedbacks", f.ListFeedback)
		r.Get("/id/{id}", f.GetFeedbackByID)
	})

	// Notification Routes
	r.Route("/notifications", func(r chi.Router) {
		r.Post("/create", n.CreateNotification)
		r.Get("/all_notifications", n.ListNotification)
		r.Get("/id/{id}", n.GetNotificationByID)
	})

	return r

}
