package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"backend/internal/application"
	"backend/internal/domain/order"
	"backend/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type OrderHandler struct {
	UC *application.OrderService
}

func NewOrderHandler(uc *application.OrderService) *OrderHandler {
	return &OrderHandler{UC: uc}
}

// CreateOrder godoc
// @Summary Create a new order
// @Security BearerAuth
// @Description Creates an order and returns IDs of created orders
// @Tags orders
// @Accept json
// @Produce json
// @Param order body order.CreateOrderRequest true "Order payload"
// @Success 201 {object} map[string][]uuid.UUID "Created order IDs"
// @Failure 400 {object} handlers.ErrorResponse "Bad request"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 409 {object} handlers.ErrorResponse "Conflict"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /orders/create [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req order.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Get customerID from middleware/context
	customerID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	// Create orders
	orders, err := h.UC.Orders.UseCase.CreateOrder(ctx, customerID, &req)
	if err != nil {
		switch {
		case errors.Is(err, order.ErrorOutOfStock):
			writeJSONError(w, http.StatusConflict, "Product is out of stock", err)
		case errors.Is(err, order.ErrorInvalidQuantity):
			writeJSONError(w, http.StatusConflict, "Invalid product quantity", err)
		case errors.Is(err, order.ErrorVariantRequired):
			writeJSONError(w, http.StatusConflict, "Variant required", err)
		case errors.Is(err, order.ErrorVariantNotAllowed):
			writeJSONError(w, http.StatusConflict, "Variant not allowed for this product", err)
		default:
			writeJSONError(w, http.StatusInternalServerError, "Failed to create order", err)
		}
		return
	}

	// Respond with all created order IDs
	orderIDs := make([]uuid.UUID, len(orders))
	for i, o := range orders {
		orderIDs[i] = o.ID
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"order_ids": orderIDs,
	})
}

// CreatePending godoc
// @Summary Create pending orders (cart)
// @Security BearerAuth
// @Description Creates pending orders and returns IDs
// @Tags orders
// @Accept json
// @Produce json
// @Param order body order.CreateOrderRequest true "Order payload"
// @Success 200 {object} map[string][]uuid.UUID "Pending order IDs"
// @Failure 400 {object} handlers.ErrorResponse "Bad request"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /orders/pending [post]
func (h *OrderHandler) CreatePending(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req order.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	customerID, err := middleware.GetUserIDFromContext(ctx)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	orders, err := h.UC.Orders.UseCase.CreatePendingOrders(ctx, customerID, &req)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to create pending orders", err)
		return
	}

	orderIDs := make([]uuid.UUID, len(orders))
	for i, o := range orders {
		orderIDs[i] = o.ID
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"order_ids": orderIDs,
	})
}

// GetOrderByID godoc
// @Summary Get order by ID
// @Security BearerAuth
// @Description Fetch a single order using UUID
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} order.OrderDoc
// @Failure 400 {string} handlers.ErrorResponse "Invalid ID"
// @Failure 404 {string} handlers.ErrorResponse "Not found"
// @Router /orders/by-id/{id} [get]
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid order ID", nil)
		return
	}
	o, err := h.UC.Orders.UseCase.GetOrder(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Order not found", err)
		return
	}

	writeJSON(w, http.StatusOK, o)
}

// GetOrderByCustomer godoc
// @Summary Get order by Customer ID
// @Security BearerAuth
// @Description Fetch order(s) using Customer ID
// @Tags orders
// @Produce json
// @Param customer_id path string true "Customer ID"
// @Success 200 {object} []order.OrderDoc
// @Failure 400 {string} handlers.ErrorResponse "Invalid Customer ID"
// @Failure 404 {string} handlers.ErrorResponse "Not found"
// @Router /orders/by-customer/{customer_id} [get]
func (h *OrderHandler) GetOrderByCustomer(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "customer_id")
	customerID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid customer ID", nil)
		return
	}

	o, err := h.UC.Orders.UseCase.GetOrderByCustomer(r.Context(), customerID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "No orders found", err)
		return
	}

	writeJSON(w, http.StatusOK, o)
}

// UpdateOrder godoc
// @Summary Update Order
// @Security BearerAuth
// @Description Update any order struct field of an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Param update body order.UpdateOrderRequest true "Field and value to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handlers.ErrorResponse "Invalid order ID or request body"
// @Failure 404 {object} handlers.ErrorResponse "Not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /orders/{id}/update [put]
func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid order ID", nil)
		return
	}

	var req order.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	column := strings.TrimSpace(strings.ToLower(req.Column))
	if column == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing or invalid column name", err)
		return
	}

	if err := h.UC.Orders.UseCase.UpdateOrder(r.Context(), orderID, column, req.Value); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update order", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("order %s updated successfully", column),
	})
}

// ListOrders godoc
// @Summary List all orders
// @Security BearerAuth
// @Description Get a list of all orders
// @Tags orders
// @Produce  json
// @Success 200 {array} order.OrderDoc
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /orders/all_orders [get]
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.UC.Orders.UseCase.ListOrders(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch orders", err)
		return
	}

	writeJSON(w, http.StatusOK, orders)
}

// DeleteOrder godoc
// @Summary Delete a order
// @Description Permanently deletes an order by their ID
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} map[string]string "Order deleted"
// @Failure 400 {object} handlers.ErrorResponse "Invalid order ID"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /orders/{id} [delete]
func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid order ID", nil)
		return
	}

	if err := h.UC.Orders.UseCase.DeleteOrder(r.Context(), orderID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete order", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("order %s deleted", orderID),
	})
}

// @Summary Run auto-assignment for pending orders
// @Security BearerAuth
// @Description Assign nearest available drivers to pending orders
// @Tags orders
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /orders/assign [post]
func (h *OrderHandler) AutoAssignOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	assignments, err := h.UC.OrderAssignment(ctx, 5000) // 5km radius
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Assignment failed", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message":     "Auto-assignment complete",
		"assignments": assignments,
		"count":       len(assignments),
	})
}
