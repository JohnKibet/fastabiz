package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"logistics-backend/internal/domain/order"
	usecase "logistics-backend/internal/usecase/order"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OrderHandler struct {
	UC *usecase.UseCase
}

func NewOrderHandler(uc *usecase.UseCase) *OrderHandler {
	return &OrderHandler{UC: uc}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Creates an order and returns the new object
// @Tags orders
// @Accept json
// @Produce json
// @Param order body order.CreateOrderRequest true "Order input"
// @Success 201 {object} order.Order
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal server error"
// @Router /orders/create [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req order.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	o := req.ToOrder()

	if err := h.UC.CreateOrder(r.Context(), o); err != nil {
		log.Printf("create order failed: %v", err)
		http.Error(w, "could not create order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":                o.ID,
		"customer_id":       o.CustomerID,
		"pickup_location":   o.PickupLocation,
		"delivery_location": o.DeliveryLocation,
		"order_status":      o.OrderStatus,
		"created_at":        o.CreatedAt,
		"updated_at":        o.UpdatedAt,
	})
}

// GetOrderByID godoc
// @Summary Get order by ID
// @Description Fetch a single order using UUID
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} order.Order
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Not found"
// @Router /orders/id/{id} [get]
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid order ID", http.StatusBadRequest)
		return
	}
	o, err := h.UC.GetOrder(r.Context(), id)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(o)
}

// GetOrderByCustomer godoc
// @Summary Get order by Customer ID
// @Description Fetch order(s) using Customer ID
// @Tags orders
// @Produce json
// @Param customer_id path string true "Customer ID"
// @Success 200 {object} []order.Order
// @Failure 400 {string} string "Invalid Customer ID"
// @Failure 404 {string} string "Not found"
// @Router /orders/customer_id/{customer_id} [get]
func (h *OrderHandler) GetOrderByCustomer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "customer_id")
	fmt.Println("id:", idStr)
	id, err := uuid.Parse(idStr)
	fmt.Println("parsed id:", id)
	if err != nil {
		http.Error(w, "invalid customer ID", http.StatusBadRequest)
		return
	}

	o, err := h.UC.GetOrderByCustomer(r.Context(), id)
	if err != nil {
		http.Error(w, "no orders found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(o)
}

// UpdateOrderStatus godoc
// @Summary Update Order Status
// @Description Update the status of an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Param status body order.UpdateOrderStatusRequest true "New Order Status"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/{order_id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "order_id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid order ID", http.StatusBadRequest)
		return
	}

	var req order.UpdateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	newStatus := order.OrderStatus(req.Status)

	if err := h.UC.UpdateOrderStatus(r.Context(), orderID, newStatus); err != nil {
		http.Error(w, "could not update order status", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "order status updated successfully",
	})
}

// ListOrders godoc
// @Summary List all orders
// @Description Get a list of all orders
// @Tags orders
// @Produce  json
// @Success 200 {array} order.Order
// @Router /orders/all_orders [get]
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.UC.ListOrders(r.Context())
	if err != nil {
		http.Error(w, "could not fetch orders", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(orders)
}
