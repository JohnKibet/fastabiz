package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"backend/internal/application"
	"backend/internal/domain/user"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type UserHandler struct {
	UC *application.OrderService
}

// ErrorResponse is a generic error model for API responses.
// swagger:model
type ErrorResponse struct {
	Error  string `json:"error" example:"Invalid request"`                               // user-friendly message
	Detail string `json:"detail,omitempty" example:"validation failed on field 'email'"` // optional internal error
}

func NewUserHandler(uc *application.OrderService) *UserHandler {
	return &UserHandler{UC: uc}
}

func writeJSONError(w http.ResponseWriter, status int, message string, internalErr error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{
		Error: message,
	}

	// Optional: Only expose internal errors during development or server errors
	if status >= 500 && internalErr != nil {
		resp.Detail = internalErr.Error()
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to write error response: %v", err)
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Register a new user with name, email, etc.
// @Tags public
// @Accept  json
// @Produce  json
// @Param user body user.CreateUserRequest true "User Input"
// @Success 201 {object} user.User
// @Failure 400 {object} handlers.ErrorResponse "Bad request"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /public/create [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req user.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// Basic manual validation
	if req.FullName == "" || req.Email == "" || req.Password == "" || req.Role == "" || req.Phone == "" {
		writeJSONError(w, http.StatusBadRequest, "Missing required fields", nil)
		return
	}

	u := req.ToUser()

	if err := h.UC.Users.UseCase.RegisterUser(r.Context(), u); err != nil {
		log.Printf("failed to create user: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":         u.ID,
		"fullName":   u.FullName,
		"password":   u.PasswordHash,
		"email":      u.Email,
		"role":       u.Role,
		"status":     u.Status,
		"phone":      u.Phone,
		"slug":       u.Slug,
		"created_at": u.CreatedAt,
	})
}

// UpdateDriverProfile godoc
// @Summary Update user (driver) phone number
// @Description Updates only the phone number of a driver (commonly used after onboarding)
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param body body user.UpdateDriverUserProfileRequest true "Driver phone update payload"
// @Success 200 {object} map[string]string "Profile updated successfully"
// @Failure 400 {object} handlers.ErrorResponse "Bad request"
// @Failure 404 {object} handlers.ErrorResponse "User not found"
// @Failure 500 {object} handlers.ErrorResponse "Server error"
// @Router /users/{id}/driver_profile [patch]
func (h *UserHandler) UpdateDriverProfile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	var req user.UpdateDriverUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.UC.Users.UseCase.UpdateDriverProfile(r.Context(), userID, &req); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update user(driver) profile", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "User(driver) profile updated"})
}

// UpdateUserProfile godoc
// @Summary Update user profile
// @Description Updates the user's name, email, and/or phone number. Partial updates allowed.
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param body body user.UpdateUserProfileRequest true "User profile update payload"
// @Success 200 {object} map[string]string "Profile updated successfully"
// @Failure 400 {object} handlers.ErrorResponse "Bad request"
// @Failure 404 {object} handlers.ErrorResponse "User not found"
// @Failure 500 {object} handlers.ErrorResponse "Server error"
// @Router /users/{id}/profile [patch]
func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	var req user.UpdateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.UC.Users.UseCase.UpdateUserProfile(r.Context(), userID, &req); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update user profile", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "User profile updated"})
}

// UpdateUserStatus godoc
// @Summary      Update a user's status
// @Description  Updates the status (e.g., active, inactive, suspended) of a specific user.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string                        true  "User ID (UUID)"
// @Param        body body      user.UpdateUserStatusRequest  true  "User status update payload"
// @Success      200  {object}  map[string]string             "User status updated"
// @Failure      400  {object}  ErrorResponse                 "Invalid ID or request body"
// @Failure      500  {object}  ErrorResponse                 "Internal server error"
// @Router       /users/{id}/status [patch]
func (h *UserHandler) UpdateUserStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	var req user.UpdateUserStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.UC.Users.UseCase.UpdateStatus(r.Context(), userID, &req); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update user status", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "User status updated"})
}

// UpdateUser godoc
// @Summary Update a specific user field
// @Description Updates a user's specific field (e.g., FullName, Email) based on user ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param data body user.UpdateUserRequest true "Field and value to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} handlers.ErrorResponse "Invalid user ID or request body"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /users/{id}/update [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	var req user.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.UC.Users.UseCase.UpdateUser(r.Context(), userID, &req); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("user %s updated successfully", req.Column),
	})
}

// ChangePassword godoc
// @Summary Change a user's password
// @Description Verifies the user's current password and updates it to a new password
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param data body user.ChangePasswordRequest true "Current and new passwords"
// @Success 200 {object} map[string]string "Password updated successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid user ID or request body"
// @Failure 401 {object} handlers.ErrorResponse "Current password incorrect"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /users/{id}/password [put]
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid user ID", nil)
		return
	}

	var req user.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.UC.Users.UseCase.ChangePassword(r.Context(), userID, &req); err != nil {
		// you can optionally inspect errors here:
		// if errors.Is(err, domain.ErrInvalidPassword) ...
		writeJSONError(w, http.StatusInternalServerError, "Failed to update password", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "password updated successfully",
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Security BearerAuth
// @Description Retrieve a user by their ID
// @Tags users
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} user.User
// @Failure 400 {string} handlers.ErrorResponse "Invalid ID"
// @Failure 404 {string} handlers.ErrorResponse "User not found"
// @Router /users/by-id/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	u, err := h.UC.Users.UseCase.GetUserByID(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "User not found", err)
		return
	}

	writeJSON(w, http.StatusOK, u)
}

// GetUserByEmail godoc
// @Summary Get user by Email
// @Security BearerAuth
// @Description Retrieve a user by their Email
// @Tags users
// @Produce  json
// @Param email path string true "User Email"
// @Success 200 {object} user.User
// @Failure 400 {string} handlers.ErrorResponse "Invalid Email"
// @Failure 404 {string} handlers.ErrorResponse "User not found"
// @Router /users/by-email/{email} [get]
func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	emailParam := chi.URLParam(r, "email")
	email, err := url.PathUnescape(emailParam)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid email format", nil)
		return
	}

	u, err := h.UC.Users.UseCase.GetUserByEmail(r.Context(), email)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "User not found", err)
		return
	}

	writeJSON(w, http.StatusOK, u)
}

// ListUsers godoc
// @Summary List all users
// @Security BearerAuth
// @Description Get a list of all registered users
// @Tags users
// @Produce  json
// @Success 200 {array} user.User
// @Router /users/all_users [get]
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.UC.Users.UseCase.ListUsers(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Could not fetch users", err)
		return
	}

	writeJSON(w, http.StatusOK, users)
}

// LoginUser godoc
// @Summary Login user
// @Description Authenticates a user using email and password and returns a JWT token.
// @Tags public
// @Accept  json
// @Produce  json
// @Param user body user.LoginRequest true "User login credentials"
// @Success 200 {object} user.LoginResponse
// @Failure 400 {string} handlers.ErrorResponse "Invalid request"
// @Failure 401 {string} handlers.ErrorResponse "Invalid credentials"
// @Failure 500 {string} handlers.ErrorResponse "Internal server error"
// @Router /public/login [post]
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req user.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request", nil)
		return
	}

	u, err := h.UC.Users.UseCase.GetUserByEmail(r.Context(), req.Email)
	if err != nil || !u.ComparePassword(req.Password) {
		writeJSONError(w, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	// Update last login
	reqUpdate := &user.UpdateUserRequest{
		Column: "last_login",
		Value:  time.Now(),
	}
	if err := h.UC.Users.UseCase.UpdateUser(r.Context(), u.ID, reqUpdate); err != nil {
		log.Printf("failed to update last login for user %s: %v", u.ID, err)
	}

	// Load the JWT secret from env
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		writeJSONError(w, http.StatusInternalServerError, "JWT secret not configured", err)
		return
	}

	// Create the token
	claims := jwt.MapClaims{
		"iss":        "my-client",   // Kong
		"sub":        u.ID.String(), // subject
		"email":      u.Email,
		"role":       u.Role, // custom claim
		"name":       u.FullName,
		"phone":      u.Phone,
		"slug":       u.Slug,
		"status":     string(u.Status),
		"last_login": u.LastLogin.Unix(),
		"created_at": u.CreatedAt.Unix(),
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // expires in 24h
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign it using the secret
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to sign token", err)
		return
	}

	// Return the token in the response
	response := user.LoginResponse{
		ID:       u.ID.String(),
		FullName: u.FullName,
		Email:    u.Email,
		Role:     string(u.Role),
		Token:    signedToken,
	}

	writeJSON(w, http.StatusOK, response)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Permanently deletes a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]string "User profile deleted"
// @Failure 400 {object} handlers.ErrorResponse "Invalid user ID"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(idStr)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	if err := h.UC.Users.UseCase.DeleteUser(r.Context(), userID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete user", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "User profile deleted"})
}
