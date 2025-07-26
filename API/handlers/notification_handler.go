package handlers

import (
	"forum/middleware"
	"forum/models"
	"forum/repository"
	"forum/utils"
	"net/http"
)

// NotificationHandler handles notification endpoints
type NotificationHandler struct {
	Repo *repository.NotificationRepository
}

func NewNotificationHandler(repo *repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{Repo: repo}
}

// GetNotifications returns notifications for the current user
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	notifications, err := h.Repo.GetByUser(user.ID)
	if err != nil {
		utils.ErrorResponse(w, "Failed to load notifications", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, notifications, http.StatusOK)
}

// MarkRead marks a single notification as read
func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id := utils.GetLastPathParam(r)
	if id == "" {
		utils.ErrorResponse(w, "Missing notification ID", http.StatusBadRequest)
		return
	}
	if err := h.Repo.MarkRead(id); err != nil {
		utils.ErrorResponse(w, "Failed to update notification", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{"status": "read"}, http.StatusOK)
}

// MarkAllRead marks all notifications for the user as read
func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := h.Repo.MarkAllRead(user.ID); err != nil {
		utils.ErrorResponse(w, "Failed to update notifications", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{"status": "read all"}, http.StatusOK)
}

// Delete deletes a single notification (soft delete)
func (h *NotificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id := utils.GetLastPathParam(r)
	if id == "" {
		utils.ErrorResponse(w, "Missing notification ID", http.StatusBadRequest)
		return
	}
	if err := h.Repo.SoftDelete(id); err != nil {
		utils.ErrorResponse(w, "Failed to delete notification", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{"status": "deleted"}, http.StatusOK)
}

// DeleteAll deletes all notifications for the user (soft delete)
func (h *NotificationHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := h.Repo.SoftDeleteAll(user.ID); err != nil {
		utils.ErrorResponse(w, "Failed to delete notifications", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{"status": "deleted all"}, http.StatusOK)
}
