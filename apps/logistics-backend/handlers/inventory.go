package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"logistics-backend/internal/domain/inventory"
	context "logistics-backend/internal/middleware"
	usecase "logistics-backend/internal/usecase/inventory"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	UC *usecase.UseCase
}

func NewInventoryHandler(uc *usecase.UseCase) *InventoryHandler {
	return &InventoryHandler{UC: uc}
}

// @Summary Create a new inventory
// @Security JWT
// @Description Creates an inventory and returns the created object
// @Tags inventories
// @Accept json
// @Produce json
// @Param inventory body inventory.CreateInventoryRequest true "Inventory input"
// @Success 201 {object} inventory.Inventory
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /inventories/create [post]
func (h *InventoryHandler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	var req inventory.CreateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	i := req.ToInventory()
	if err := h.UC.CreateInventory(r.Context(), i); err != nil {
		log.Printf("create inventory failed: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "could not create inventory")
		return
	}

	adminID, err := context.GetAdminIDFromContext(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized access")
		return
	}
	req.AdminID = adminID

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":          i.ID,
		"admin_id":    req.AdminID,
		"name":        i.Name,
		"category":    i.Category,
		"stock":       i.Stock,
		"price":       i.Price,
		"images":      i.Images,
		"unit":        i.Unit,
		"packaging":   i.Packaging,
		"description": i.Description,
		"location":    i.Location,
		"created_at":  i.CreatedAt,
		"updated_at":  i.UpdatedAt,
	})
}

// @Summary Get inventory by ID
// @Security JWT
// @Description Get a specific inventory item by UUID
// @Tags inventories
// @Produce json
// @Param id query string true "Inventory ID"
// @Success 200 {object} inventory.Inventory
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /inventories/{id} [get]
func (h *InventoryHandler) GetByInventoryID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid inventory ID")
		return
	}

	i, err := h.UC.GetByID(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "No inventory found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(i)
}

// @Summary Get inventory by name
// @Security JWT
// @Description Search inventory by item name (exact match)
// @Tags inventories
// @Produce json
// @Param name query string true "Inventory Name"
// @Success 200 {object} inventory.Inventory
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /inventories/by-name [get]
func (h *InventoryHandler) GetByInventoryName(w http.ResponseWriter, r *http.Request) {
	nameStr := r.URL.Query().Get("name")
	if nameStr == "" {
		writeJSONError(w, http.StatusBadRequest, "Name query parameter is required")
		return
	}

	i, err := h.UC.GetByName(r.Context(), nameStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, http.StatusNotFound, fmt.Sprintf("No inventory found with name '%s'", nameStr))
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(i)
}

// @Summary List all inventories
// @Security JWT
// @Description List all inventories with optional pagination
// @Tags inventories
// @Produce json
// @Param limit query int false "Limit number of items"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} inventory.Inventory
// @Failure 500 {object} map[string]string
// @Router /inventories/all_inventories [get]
func (h *InventoryHandler) ListInventories(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 0 // default
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	inventories, err := h.UC.List(r.Context(), limit, offset)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch orders")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventories)
}

// @Summary Get inventories by category
// @Security JWT
// @Description Get all inventory items in a specific category
// @Tags inventories
// @Produce json
// @Param category query string true "Category Name"
// @Success 200 {array} inventory.Inventory
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /inventories/by-category [get]
func (h *InventoryHandler) GetInventoryByCategory(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		writeJSONError(w, http.StatusBadRequest, "Category query parameter is required")
		return
	}

	inventories, err := h.UC.GetByCategory(r.Context(), category)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "No Inventories found for this category")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventories)
}

// @Summary List all categories
// @Security JWT
// @Description Get all unique categories that belong to existing inventories
// @Tags inventories
// @Produce json
// @Success 200 {array} string
// @Failure 500 {object} map[string]string
// @Router /inventories/categories [get]
func (h *InventoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.UC.ListCategories(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// @Summary View public product page
// @Description View a specific product by admin slug and product slug
// @Tags public
// @Produce json
// @Param adminSlug path string true "Admin Slug"
// @Param productSlug path string true "Product Slug"
// @Success 200 {object} inventory.Inventory
// @Failure 404 {object} map[string]string
// @Router /public/store/{adminSlug}/product/{productSlug} [get]
func (h *InventoryHandler) GetPublicProductPage(w http.ResponseWriter, r *http.Request) {
	adminSlug := chi.URLParam(r, "adminSlug")
	productSlug := chi.URLParam(r, "productSlug")

	product, err := h.UC.GetBySlugs(r.Context(), adminSlug, productSlug)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Product not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// @Summary View public store page
// @Description View all public products for an admin store by slug
// @Tags public
// @Produce json
// @Param adminSlug path string true "Admin Slug"
// @Success 200 {object} inventory.StorePublicView
// @Failure 404 {object} map[string]string
// @Router /public/store/{adminSlug} [get]
func (h *InventoryHandler) GetAdminStorePage(w http.ResponseWriter, r *http.Request) {
	adminSlug := chi.URLParam(r, "adminSlug")

	storeView, err := h.UC.GetStorePublicView(r.Context(), adminSlug)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Store not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(storeView)
}
