package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"logistics-backend/internal/domain/delivery"
	usecase "logistics-backend/internal/usecase/delivery"
	"net/http"
	"strings"

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
// @Security JWT
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
		writeJSONError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	d := req.ToDelivery()

	if err := dh.DH.CreateDelivery(r.Context(), d); err != nil {
		log.Printf("create delivery failed: %v", err)

		switch err {
		case delivery.ErrorNoPendingOrder:
			writeJSONError(w, http.StatusConflict, "No pending orders")
		default:
			writeJSONError(w, http.StatusInternalServerError, "could not create delivery")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
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
// @Security JWT
// @Description Retrieve a delivery by their ID
// @Tags deliveries
// @Produce  json
// @Param id path string true "Delivery ID"
// @Success 200 {object} delivery.Delivery
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Delivery not found"
// @Router /deliveries/by-id/{id} [get]
func (dh *DeliveryHandler) GetDeliveryByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	deliveryID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid delivery ID")
		return
	}

	d, err := dh.DH.GetDeliveryByID(r.Context(), deliveryID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "No delivery found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

// UpdateDelivery godoc
// @Summary Update Delivery
// @Security JWT
// @Description Update any delivery struct field of an existing delivery
// @Tags deliveries
// @Accept json
// @Produce json
// @Param delivery_id path string true "Delivery ID"
// @Param update body delivery.UpdateDeliveryRequest true "Field and value to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /deliveries/{id}/update [put]
func (dh *DeliveryHandler) UpdateDelivery(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	deliveryID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid delivery ID")
		return
	}

	var req delivery.UpdateDeliveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	column := strings.TrimSpace(strings.ToLower(req.Column))
	if column == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing or invalid column name")
		return
	}

	if err := dh.DH.UpdateDelivery(r.Context(), deliveryID, column, req.Value); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("delivery %s updated successfully", column),
	})

}

// ListDeliveries godoc
// @Summary List all deliveries
// @Security JWT
// @Description Get a list of all deliveries
// @Tags deliveries
// @Produce  json
// @Success 200 {array} delivery.Delivery
// @Router /deliveries/all_deliveries [get]
func (dh *DeliveryHandler) ListDeliveries(w http.ResponseWriter, r *http.Request) {
	deliveries, err := dh.DH.ListDeliveries(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch deliveries")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deliveries)
}

// DeleteDelivery godoc
// @Summary Delete a delivery
// @Description Permanently deletes a delivery by their ID
// @Tags deliveries
// @Accept json
// @Produce json
// @Security JWT
// @Param id path string true "Delivery ID"
// @Success 200 {object} map[string]string "Delivery deleted"
// @Failure 400 {object} string "Invalid Delivery ID"
// @Failure 500 {object} string "Internal server error"
// @Router /deliveries/{id} [delete]
func (dh *DeliveryHandler) DeleteDelivery(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	deliveryID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid delivery ID")
		return
	}

	if err := dh.DH.DeleteDelivery(r.Context(), deliveryID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("delivery %s deleted", deliveryID),
	})
}
