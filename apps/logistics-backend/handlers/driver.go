package handlers

import (
	"encoding/json"
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
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.FullName == "" || req.VehicleInfo == "" || req.CurrentLocation == "" || req.Available == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	d := req.ToDriver()

	if err := dh.DH.RegisterDriver(r.Context(), d); err != nil {
		http.Error(w, "could not create driver", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
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

// GetDriverByID godoc
// @Summary Get driver by ID
// @Description Retrieve a driver by their ID
// @Tags drivers
// @Produce  json
// @Param id path string true "Driver ID"
// @Success 200 {object} driver.Driver
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Driver not found"
// @Router /drivers/id/{id} [get]
func (dh *DriverHandler) GetDriverByID(w http.ResponseWriter, r *http.Request) {
	driverID := chi.URLParam(r, "id")
	id, err := uuid.Parse(driverID)
	if err != nil {
		http.Error(w, "invalid driver ID", http.StatusBadRequest)
		return
	}

	d, err := dh.DH.GetDriverByID(r.Context(), id)
	if err != nil {
		http.Error(w, "driver not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(d)
}

// GetUserByDriver godoc
// @Summary Get driver by Email
// @Description Retrieve a driver by their Email
// @Tags drivers
// @Produce  json
// @Param email path string true "Driver Email"
// @Success 200 {object} driver.Driver
// @Failure 400 {string} string "Invalid Email"
// @Failure 404 {string} string "Driver not found"
// @Router /drivers/email/{email} [get]
func (dh *DriverHandler) GetDriverByEmail(w http.ResponseWriter, r *http.Request) {
	emailParam := chi.URLParam(r, "email")
	email, err := url.PathUnescape(emailParam)
	if err != nil {
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}

	d, err := dh.DH.GetDriverByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "driver not found", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(d)
}

// ListDrivers godoc
// @Summary List all drivers
// @Description Get a list of all registered drivers
// @Tags drivers
// @Produce  json
// @Success 200 {array} driver.Driver
// @Router /drivers/all_drivers [get]
func (dh *DriverHandler) ListDrivers(w http.ResponseWriter, r *http.Request) {
	drivers, err := dh.DH.ListDrivers(r.Context())
	if err != nil {
		http.Error(w, "could not fetch drivers", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(drivers)
}
