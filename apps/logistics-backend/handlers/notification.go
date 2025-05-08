package handlers

import (
	"encoding/json"
	"logistics-backend/internal/domain/notification"
	usecase "logistics-backend/internal/usecase/notification"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	NH *usecase.UseCase
}

func NewNotificationHandler(nh *usecase.UseCase) *NotificationHandler {
	return &NotificationHandler{NH: nh}
}

// CreateNotification godoc
// @Summary Create a new notification
// @Description Create a new notification with user_id, message, etc.
// @Tags notifications
// @Accept  json
// @Produce  json
// @Param user body notification.CreateNotificationRequest true "User Input"
// @Success 201 {object} notification.Notification
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Failed to create notification"
// @Router /notifications/create [post]
func (nh *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var req *notification.CreateNotificationRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	n := req.ToNotification()

	if err := nh.NH.CreateNotification(r.Context(), n); err != nil {
		http.Error(w, "could not create notification", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"id":      n.ID,
		"user_id": n.UserID,
		"message": n.Message,
		"type":    n.Type,
		"sent_at": n.SentAt,
	})
}

// GetNotificationByID godoc
// @Summary Get notification by ID
// @Description Retrieve a notification by their ID
// @Tags notifications
// @Produce  json
// @Param id path string true "Notification ID"
// @Success 200 {object} notification.Notification
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Notification not found"
// @Router /notifications/id/{id} [get]
func (nh *NotificationHandler) GetNotificationByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	notificationID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	n, err := nh.NH.GetNotification(r.Context(), notificationID)
	if err != nil {
		http.Error(w, "no notification found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(n)
}

// ListNotifications godoc
// @Summary List all notifications
// @Description Get a list of all notifications
// @Tags notifications
// @Produce  json
// @Success 200 {array} notification.Notification
// @Router /notifications/all_notifications [get]
func (nh *NotificationHandler) ListNotification(w http.ResponseWriter, r *http.Request) {
	n, err := nh.NH.ListNotification(r.Context())
	if err != nil {
		http.Error(w, "could not fetch notifications", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(n)
}
