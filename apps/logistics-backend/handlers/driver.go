package handlers

import (
	"encoding/json"
	"fmt"
	"logistics-backend/internal/domain/driver"
	usecase "logistics-backend/internal/usecase/driver"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DriverHandler struct {
	DH *usecase.UseCase
}

func NewDriverHandler(dh *usecase.UseCase) *DriverHandler {
	return &DriverHandler{DH: dh}
}

// CreateDriver godoc
// @Summary Create a new driver
// @Security JWT
// @Description Register a new driver with name, email, etc.
// @Tags drivers
// @Accept  json
// @Produce  json
// @Param user body driver.CreateDriverRequest true "User Input"
// @Success 201 {object} driver.Driver
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Failed to create driver"
// @Router /drivers/create [post]
func (dh *DriverHandler) CreateDriver(w http.ResponseWriter, r *http.Request) {
	var req driver.CreateDriverRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Basic validation
	if req.FullName == "" || req.VehicleInfo == "" || req.CurrentLocation == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	d := req.ToDriver()

	if err := dh.DH.RegisterDriver(r.Context(), d); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not create driver")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Driver profile created"})
	json.NewEncoder(w).Encode(map[string]any{
		"id":               d.ID,
		"full_name":        d.FullName,
		"email":            d.Email,
		"vehicle_info":     d.VehicleInfo,
		"current_location": d.CurrentLocation,
		"available":        d.Available,
		"created_at":       d.CreatedAt,
	})
}

// UpdateDriverProfile godoc
// @Summary Update driver profile
// @Description Updates the vehicle information and current location of a driver
// @Tags drivers
// @Security JWT
// @Accept json
// @Produce json
// @Param driver_id path string true "Driver ID"
// @Param body body driver.UpdateDriverRequest true "Driver profile fields to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /drivers/{id}/profile [put]
func (dh *DriverHandler) UpdateDriverProfile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	driverID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid driver ID")
		return
	}

	var req driver.UpdateDriverProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := dh.DH.UpdateDriverProfile(r.Context(), driverID, &req); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Driver profile updated"})

}

// UpdateDriver godoc
// @Summary Update a specific driver field
// @Description Updates a driver's specific field (e.g., VehicleInfo, CurrentLocation) based on driver ID
// @Tags drivers
// @Accept json
// @Produce json
// @Security JWT
// @Param id path string true "Driver ID"
// @Param data body driver.UpdateDriverRequest true "Field and value to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} string "Invalid driver ID or request body"
// @Failure 500 {object} string "Internal server error"
// @Router /drivers/{id}/update [put]
func (dh *DriverHandler) UpdateDriver(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	driverID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid driver ID")
		return
	}

	var req driver.UpdateDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := dh.DH.UpdateDriver(r.Context(), driverID, &req); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("driver %s updated successfully", req.Column),
	})
}

// GetDriverByID godoc
// @Summary Get driver by ID
// @Security JWT
// @Description Retrieve a driver by their ID
// @Tags drivers
// @Produce  json
// @Param id path string true "Driver ID"
// @Success 200 {object} driver.Driver
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Driver not found"
// @Router /drivers/{id} [get]
func (dh *DriverHandler) GetDriverByID(w http.ResponseWriter, r *http.Request) {
	driverID := chi.URLParam(r, "id")
	id, err := uuid.Parse(driverID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid driver ID")
		return
	}

	d, err := dh.DH.GetDriverByID(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Driver not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

// GetUserByDriver godoc
// @Summary Get driver by Email
// @Security JWT
// @Description Retrieve a driver by their Email
// @Tags drivers
// @Produce  json
// @Param email path string true "Driver Email"
// @Success 200 {object} driver.Driver
// @Failure 400 {string} string "Invalid Email"
// @Failure 404 {string} string "Driver not found"
// @Router /drivers/{email} [get]
func (dh *DriverHandler) GetDriverByEmail(w http.ResponseWriter, r *http.Request) {
	emailParam := chi.URLParam(r, "email")
	email, err := url.PathUnescape(emailParam)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	d, err := dh.DH.GetDriverByEmail(r.Context(), email)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Driver not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

// ListDrivers godoc
// @Summary List all drivers
// @Security JWT
// @Description Get a list of all registered drivers
// @Tags drivers
// @Produce  json
// @Success 200 {array} driver.Driver
// @Router /drivers/all_drivers [get]
func (dh *DriverHandler) ListDrivers(w http.ResponseWriter, r *http.Request) {
	drivers, err := dh.DH.ListDrivers(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch drivers")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drivers)
}

// DeleteDriver godoc
// @Summary Delete a driver
// @Description Permanently deletes a driver by their ID
// @Tags drivers
// @Accept json
// @Produce json
// @Security JWT
// @Param id path string true "Driver ID"
// @Success 200 {object} map[string]string "Driver deleted"
// @Failure 400 {object} string "Invalid driver ID"
// @Failure 500 {object} string "Internal server error"
// @Router /drivers/{id} [delete]
func (h *DriverHandler) DeleteDriver(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	driverID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid driver ID")
		return
	}

	if err := h.DH.DeleteDriver(r.Context(), driverID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("driver %s deleted", driverID),
	})
}
