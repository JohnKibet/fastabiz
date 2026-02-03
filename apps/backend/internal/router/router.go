package router

import (
	"net/http"

	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"backend/handlers"
	authMiddleware "backend/internal/middleware"

	"github.com/jmoiron/sqlx"
)

func NewRouter(
	u *handlers.UserHandler,
	o *handlers.OrderHandler,
	d *handlers.DriverHandler,
	e *handlers.DeliveryHandler,
	p *handlers.PaymentHandler, f *handlers.FeedbackHandler,
	n *handlers.NotificationHandler,
	publicApiBaseUrl string,
	c *handlers.InviteHandler,
	s *handlers.StoreHandler,
	pr *handlers.ProductHandler,
	db *sqlx.DB,
) http.Handler {
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

		})

		// Protected Routes (auth required)
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.JWTAuthMiddleware)

			// Users
			r.Route("/users", func(r chi.Router) {
				r.Get("/all_users", u.ListUsers)
				r.Get("/by-id/{id}", u.GetUserByID)
				r.Get("/by-email/{email}", u.GetUserByEmail)
				r.Patch("/{id}/driver_profile", u.UpdateDriverProfile)
				r.Patch("/{id}/profile", u.UpdateUserProfile)
				r.Put("/{id}/update", u.UpdateUser)
				r.Put("/{id}/password", u.ChangePassword)
				r.Patch("/{id}/status", u.UpdateUserStatus)
				r.Delete("/{id}", u.DeleteUser)
			})

			// Invites
			r.Route("/invites", func(r chi.Router) {
				r.Post("/create", c.CreateMember)
				r.Get("/by-token", c.GetMemberByToken)
				r.Get("/all_invites", c.ListPendingMembers)
				r.Delete("/{id}", c.DeleteMember)
			})

			// Orders
			r.Route("/orders", func(r chi.Router) {
				r.Post("/create", o.CreateOrder)
				r.Get("/all_orders", o.ListOrders)
				r.Post("/assign", o.AutoAssignOrders)
				r.Get("/by-id/{id}", o.GetOrderByID)
				r.Get("/by-customer/{customer_id}", o.GetOrderByCustomer)
				r.Put("/{id}/update", o.UpdateOrder)
				r.Delete("/{id}", o.DeleteOrder)
			})

			// Drivers
			r.Route("/drivers", func(r chi.Router) {
				r.Get("/all_drivers", d.ListDrivers)
				r.Get("/by-id/{id}", d.GetDriverByID)
				r.Get("/by-email/{email}", d.GetDriverByEmail)
				r.Patch("/{id}/profile", d.UpdateDriverProfile)
				r.Put("/{id}/update", d.UpdateDriver)
				r.Delete("/{id}", d.DeleteDriver)
			})

			// Deliveries
			r.Route("/deliveries", func(r chi.Router) {
				r.Get("/all_deliveries", e.ListDeliveries)
				r.Get("/by-id/{id}", e.GetDeliveryByID)
				r.Put("/{id}/update", e.UpdateDelivery)
				r.Put("/{id}/accept", e.AcceptDelivery)
				r.Delete("/{id}", e.DeleteDelivery)
			})

			// Payments
			r.Route("/payments", func(r chi.Router) {
				r.Post("/create", p.CreatePayment)
				r.Get("/all_payments", p.ListPayments)
				r.Get("/{id}", p.GetPaymentByID)
				r.Get("/{order_id}", p.GetPaymentByOrderID)

				// MPesa STK Push
				r.Post("/mpesa-express", p.MpesaExpress)
				r.Post("/mpesa-callback", p.MpesaCallback)
			})

			// Feedbacks
			r.Route("/feedbacks", func(r chi.Router) {
				r.Post("/create", f.CreateFeedback)
				r.Get("/all_feedbacks", f.ListFeedback)
				r.Get("/{id}", f.GetFeedbackByID)
			})

			// Notifications
			r.Route("/notifications", func(r chi.Router) {
				r.Post("/create", n.CreateNotification)
				r.Get("/all_pending_notifications", n.ListNotifications)
				r.Get("/all_my_notifications/{id}", n.ListUserNotifications)
				r.Put("/{id}/status", n.UpdateNotificationStatus)
				r.Patch("/{id}/read", n.MarkAsRead)
				r.Patch("/mark_all_as_read/{id}", n.MarkAllAsRead)
			})

			// Stores
			r.Route("/stores", func(r chi.Router) {
				r.Post("/create", s.CreateStore)
				r.Get("/all_stores", s.ListStores)
				r.Get("/me", s.ListOwnerStores)
				r.Get("/me/paged", s.ListStoresPaged)
				r.Get("/by-id/{id}", s.GetStoreByID)
				r.Get("/{id}/summary", s.GetStoreSummary)
				r.Put("/{id}/update", s.UpdateStore)
				r.Delete("/{id}/delete", s.DeleteStore)
			})

			// Products
			r.Route("/products", func(r chi.Router) {
				r.Post("/cloudinary/signature", pr.CloudinarySignature)
				r.Post("/create", pr.CreateProduct)
				r.Get("/{store_id}/all_products", pr.ListProductsByStore)
				r.Post("/images/add", pr.AddImage)
				r.Post("/options/add", pr.AddOptionName)
				r.Post("/options/values/add", pr.AddOptionValue)
				r.Post("/variants/add", pr.CreateVariant)
				r.Patch("/variants/stock/update", pr.UpdateVariantStock)
				r.Patch("/variants/price/update", pr.UpdateVariantPrice)
				r.Patch("/images/reorder", pr.ReorderImages)
				r.Patch("/{id}/product_details", pr.UpdateProductDetails)
				r.Get("/{productId}/options", pr.ListOptions)
				r.Get("/by-id/{id}", pr.GetProductByID)
				r.Delete("/{id}/delete", pr.DeleteProduct)
				r.Delete("/images/{imageId}/delete", pr.DeleteImage)
				r.Delete("/options/{optionId}/delete", pr.DeleteOptionName)
				r.Delete("/options/values/{valueId}/delete", pr.DeleteOptionValue)
				r.Delete("/variants/{variantId}/delete", pr.DeleteVariant)
			})
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("db unreachable"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	return r
}
