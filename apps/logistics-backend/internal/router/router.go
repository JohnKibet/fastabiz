package router

import (
	"net/http"

	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"logistics-backend/handlers"
	authMiddleware "logistics-backend/internal/middleware"
)

func NewRouter(u *handlers.UserHandler, o *handlers.OrderHandler, d *handlers.DriverHandler, e *handlers.DeliveryHandler, p *handlers.PaymentHandler, f *handlers.FeedbackHandler, n *handlers.NotificationHandler, i *handlers.InventoryHandler, publicApiBaseUrl string) http.Handler {
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

	r.Route("/api", func(r chi.Router) {
		// Swagger docs (this will now be served under /api/swagger)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(publicApiBaseUrl+"/swagger/doc.json"),
		))

		// Public routes
		r.Route("/public", func(r chi.Router) {
			// Public auth
			r.Post("/create", u.CreateUser)
			r.Post("/login", u.LoginUser)

			// Public store pages
			r.Route("/store", func(r chi.Router) {
				r.Get("/{adminSlug}/product/{productSlug}", i.GetPublicProductPage)
				r.Get("/{adminSlug}", i.GetAdminStorePage)
			})
		})

		// Protected Routes (auth required)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.JWTAuthMiddleware)

			// Users
			r.Route("/users", func(r chi.Router) {
				r.Get("/all_users", u.ListUsers)
				r.Get("/id/{id}", u.GetUserByID)
				r.Get("/email/{email}", u.GetUserByEmail)
			})

			// Orders
			r.Route("/orders", func(r chi.Router) {
				r.Post("/create", o.CreateOrder)
				r.Get("/all_orders", o.ListOrders)
				r.Get("/id/{id}", o.GetOrderByID)
				r.Get("/customer_id/{customer_id}", o.GetOrderByCustomer)
				r.Put("/{order_id}/status", o.UpdateOrderStatus)
			})

			// Inventories
			r.Route("/inventories", func(r chi.Router) {
				r.Post("/create", i.CreateInventory)
				r.Get("/inventory_id/{inventory_id}", i.GetByInventoryID)
				r.Get("/inventory_name/{inventory_name}", i.GetByInventoryName)
				r.Get("/all_inventories", i.ListInventories)
			})

			// Drivers
			r.Route("/drivers", func(r chi.Router) {
				r.Post("/create", d.CreateDriver)
				r.Get("/all_drivers", d.ListDrivers)
				r.Get("/id/{id}", d.GetDriverByID)
				r.Get("/email/{email}", d.GetDriverByEmail)
			})

			// Deliveries
			r.Route("/deliveries", func(r chi.Router) {
				r.Post("/create", e.CreateDelivery)
				r.Get("/all_deliveries", e.ListDeliveries)
				r.Get("/id/{id}", e.GetDeliveryByID)
			})

			// Payments
			r.Route("/payments", func(r chi.Router) {
				r.Post("/create", p.CreatePayment)
				r.Get("/all_payments", p.ListPayments)
				r.Get("/id/{id}", p.GetPaymentByID)
				r.Get("/order_id/{order_id}", p.GetPaymentByOrderID)
			})

			// Feedbacks
			r.Route("/feedbacks", func(r chi.Router) {
				r.Post("/create", f.CreateFeedback)
				r.Get("/all_feedbacks", f.ListFeedback)
				r.Get("/id/{id}", f.GetFeedbackByID)
			})

			// Notifications
			r.Route("/notifications", func(r chi.Router) {
				r.Post("/create", n.CreateNotification)
				r.Get("/all_notifications", n.ListNotification)
				r.Get("/id/{id}", n.GetNotificationByID)
			})
		})
	})

	return r
}
