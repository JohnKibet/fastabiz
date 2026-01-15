package handlers

import (
	"backend/internal/application"
	"backend/internal/domain/store"
	"backend/internal/middleware"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type StoreHandler struct {
	UC *application.OrderService
}

func NewStoreHandler(uc *application.OrderService) *StoreHandler {
	return &StoreHandler{UC: uc}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write json response: %v", err)
	}
}

// CreateStore godoc
// @Summary Create a new store
// @Security BearerAuth
// @Description Creates a new store for the authenticated owner
// @Tags stores
// @Accept json
// @Produce json
// @Param store body store.CreateStoreRequest true "Store input"
// @Success 201 {object} map[string]any "Created store"
// @Failure 400 {object} handlers.ErrorResponse "Invalid request"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /stores/create [post]
func (h *StoreHandler) CreateStore(w http.ResponseWriter, r *http.Request) {
	var req store.CreateStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	ownerID, err := middleware.GetOwnerIDFromContext(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	s := req.ToStore()
	s.OwnerID = ownerID

	if err := h.UC.Stores.UseCase.CreateStore(r.Context(), s); err != nil {
		if errors.Is(err, store.ErrStoreNameAlreadyExists) {
			writeJSONError(w, http.StatusConflict, "A store with this name already exists", nil)
			return
		}

		writeJSONError(w, http.StatusInternalServerError, "Could not create store", err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":         s.ID,
		"owner_id":   s.OwnerID,
		"name":       s.Name,
		"logo_url":   s.LogoURL,
		"location":   s.Location,
		"created_at": s.CreatedAt,
		"updated_at": s.UpdatedAt,
	})
}

// GetStoreByID godoc
// @Summary Get a store by ID
// @Description Fetch a single store by its UUID
// @Tags stores
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {object} store.Store
// @Failure 400 {object} handlers.ErrorResponse "Invalid ID"
// @Failure 404 {object} handlers.ErrorResponse "Store not found"
// @Router /stores/by-id/{id} [get]
func (h *StoreHandler) GetStoreByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid store ID", nil)
		return
	}

	s, err := h.UC.Stores.UseCase.GetStoreByID(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Store not found", err)
		return
	}

	writeJSON(w, http.StatusOK, s)
}

// GetStoreSummary godoc
// @Summary Get store summary
// @Description Fetch summarized info about a store
// @Tags stores
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {object} store.StoreSummary
// @Failure 400 {object} handlers.ErrorResponse "Invalid ID"
// @Failure 404 {object} handlers.ErrorResponse "Store not found"
// @Router /stores/{id}/summary [get]
func (h *StoreHandler) GetStoreSummary(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	storeID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid store ID", nil)
		return
	}

	summary, err := h.UC.Stores.UseCase.GetStoreSummary(r.Context(), storeID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Store not found", err)
		return
	}

	writeJSON(w, http.StatusOK, summary)
}

// UpdateStore godoc
// @Summary Update a store
// @Security BearerAuth
// @Description Update details of a store owned by the authenticated user
// @Tags stores
// @Accept json
// @Produce json
// @Param id path string true "Store ID"
// @Param update body store.UpdateStoreRequest true "Update request"
// @Success 200 {object} map[string]string "Update message"
// @Failure 400 {object} handlers.ErrorResponse "Invalid request"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /stores/{id}/update [put]
func (h *StoreHandler) UpdateStore(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	storeID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid store ID", nil)
		return
	}

	var req store.UpdateStoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	ownerID, err := middleware.GetOwnerIDFromContext(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	if err := h.UC.Stores.UseCase.UpdateStore(r.Context(), storeID, ownerID, &req); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update store", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "store updated successfully",
	})
}

// ListStores godoc
// @Summary List all stores
// @Description Returns all stores (admin view or public)
// @Tags stores
// @Produce json
// @Success 200 {array} store.Store
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /stores/all_stores [get]
func (h *StoreHandler) ListStores(w http.ResponseWriter, r *http.Request) {
	stores, err := h.UC.Stores.UseCase.ListAllStores(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to list stores", err)
		return
	}

	writeJSON(w, http.StatusOK, stores)
}

// ListOwnerStores godoc
// @Summary List authenticated owner's stores
// @Security BearerAuth
// @Description Returns all stores owned by the current user
// @Tags stores
// @Produce json
// @Success 200 {array} store.Store
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /stores/me [get]
func (h *StoreHandler) ListOwnerStores(w http.ResponseWriter, r *http.Request) {
	ownerID, err := middleware.GetOwnerIDFromContext(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	stores, err := h.UC.Stores.UseCase.ListOwnerStores(r.Context(), ownerID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to list stores", err)
		return
	}

	writeJSON(w, http.StatusOK, stores)
}

// ListStoresPaged godoc
// @Summary List stores with pagination
// @Security BearerAuth
// @Description Returns stores owned by authenticated user with limit and offset
// @Tags stores
// @Produce json
// @Param limit query int false "Maximum number of items to return"
// @Param offset query int false "Number of items to skip"
// @Success 200 {array} store.StoreSummary
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /stores/me/paged [get]
func (h *StoreHandler) ListStoresPaged(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // default
	}
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	ownerID, err := middleware.GetOwnerIDFromContext(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	var filter = &store.StoreFilter{
		OwnerID: &ownerID,
		Limit:   limit,
		Offset:  offset,
	}

	stores, err := h.UC.Stores.UseCase.ListStoresPaged(r.Context(), *filter)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to list stores", err)
		return
	}

	writeJSON(w, http.StatusOK, stores)
}

// DeleteStore godoc
// @Summary Delete a store
// @Security BearerAuth
// @Description Deletes a store owned by the authenticated user
// @Tags stores
// @Produce json
// @Param id path string true "Store ID"
// @Success 200 {object} map[string]string "Deletion message"
// @Failure 400 {object} handlers.ErrorResponse "Invalid store ID"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /stores/{id}/delete [delete]
func (h *StoreHandler) DeleteStore(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	storeID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid store ID", nil)
		return
	}

	ownerID, err := middleware.GetOwnerIDFromContext(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	name, err := h.UC.Stores.UseCase.DeleteStore(r.Context(), storeID, ownerID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete store", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": fmt.Sprintf("Store '%q' deleted successfully", name),
	})
}
