package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"logistics-backend/internal/domain/user"
	usecase "logistics-backend/internal/usecase/user"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	UC *usecase.UseCase
}

func NewUserHandler(uc *usecase.UseCase) *UserHandler {
	return &UserHandler{UC: uc}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Register a new user with name, email, etc.
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body user.CreateUserRequest true "User Input"
// @Success 201 {object} user.User
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Failed to create user"
// @Router /users/create [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Basic manual validation
	if req.FullName == "" || req.Email == "" || req.Password == "" || req.Role == "" || req.Phone == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	u := req.ToUser()

	if err := h.UC.RegisterUser(r.Context(), u); err != nil {
		log.Printf("failed to create user: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":         u.ID,
		"full_name":  u.FullName,
		"password":   u.PasswordHash,
		"email":      u.Email,
		"role":       u.Role,
		"phone":      u.Phone,
		"created_at": u.CreatedAt,
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve a user by their ID
// @Tags users
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} user.User
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "User not found"
// @Router /users/id/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	u, err := h.UC.GetUserByID(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// GetUserByEmail godoc
// @Summary Get user by Email
// @Description Retrieve a user by their Email
// @Tags users
// @Produce  json
// @Param email path string true "User Email"
// @Success 200 {object} user.User
// @Failure 400 {string} string "Invalid Email"
// @Failure 404 {string} string "User not found"
// @Router /users/email/{email} [get]
func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	emailParam := chi.URLParam(r, "email")
	email, err := url.PathUnescape(emailParam)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	u, err := h.UC.GetUserByEmail(r.Context(), email)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// ListUsers godoc
// @Summary List all users
// @Description Get a list of all registered users
// @Tags users
// @Produce  json
// @Success 200 {array} user.User
// @Router /users/all_users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UC.ListUsers(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch users")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
