package handlers

import (
	"encoding/json"
	"logistics-backend/internal/domain/delivery"
	usecase "logistics-backend/internal/usecase/delivery"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DeliveryHandler struct {
	DH *usecase.UseCase
}

func NewDeliveryHandler(dh *usecase.UseCase) *DeliveryHandler {
	return &DeliveryHandler{DH: dh}
}

// CreateDelivery godoc
// @Summary Create a new delivery
// @Description Create a new delivery with order_id, driver_id, etc.
// @Tags deliveries
// @Accept  json
// @Produce  json
// @Param user body delivery.CreateDeliveryRequest true "User Input"
// @Success 201 {object} delivery.Delivery
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Failed to create delivery"
// @Router /deliveries/create [post]
func (dh *DeliveryHandler) CreateDelivery(w http.ResponseWriter, r *http.Request) {
	var req *delivery.CreateDeliveryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	d := req.ToDelivery()

	err := dh.DH.CreateDelivery(r.Context(), d)
	if err != nil {
		http.Error(w, "could not create delivery", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"id":           d.ID,
		"order_id":     d.OrderID,
		"driver_id":    d.DriverID,
		"assigned_at":  d.AssignedAt,
		"picked_up_at": d.PickedUpAt,
		"delivered_at": d.DeliveredAt,
		"status":       d.Status,
	})
}

// GetDeliveryByID godoc
// @Summary Get delivery by ID
// @Description Retrieve a delivery by their ID
// @Tags deliveries
// @Produce  json
// @Param id path string true "Delivery ID"
// @Success 200 {object} delivery.Delivery
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Delivery not found"
// @Router /deliveries/id/{id} [get]
func (dh *DeliveryHandler) GetDeliveryByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	deliveryID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid delivery ID", http.StatusBadRequest)
		return
	}

	d, err := dh.DH.GetDeliveryByID(r.Context(), deliveryID)
	if err != nil {
		http.Error(w, "no delivery found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(d)
}

// ListDeliveries godoc
// @Summary List all deliveries
// @Description Get a list of all deliveries
// @Tags deliveries
// @Produce  json
// @Success 200 {array} delivery.Delivery
// @Router /deliveries/all_deliveries [get]
func (dh *DeliveryHandler) ListDeliveries(w http.ResponseWriter, r *http.Request) {
	deliveries, err := dh.DH.ListDeliveries(r.Context())
	if err != nil {
		http.Error(w, "could not fetch deliveries", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(deliveries)
}
