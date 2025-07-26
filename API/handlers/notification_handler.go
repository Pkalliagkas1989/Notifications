package handlers

import (
	"net/http"

	"forum/middleware"
	"forum/repository"
	"forum/utils"
)

// NotificationHandler handles notification related endpoints

type NotificationHandler struct {
	Repo *repository.NotificationRepository
}

func NewNotificationHandler(repo *repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{Repo: repo}
}

// List returns notifications for the authenticated user
func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ns, err := h.Repo.GetByUser(user.ID)
	if err != nil {
		utils.ErrorResponse(w, "Failed to load notifications", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, ns, http.StatusOK)
}

// MarkRead marks a single notification as read
func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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
		utils.ErrorResponse(w, "Failed to mark read", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{"status": "read"}, http.StatusOK)
}

// MarkAllRead marks all notifications as read
func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := h.Repo.MarkAllRead(user.ID); err != nil {
		utils.ErrorResponse(w, "Failed to mark all read", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{"status": "all read"}, http.StatusOK)
}

// Delete soft deletes a single notification
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
		utils.ErrorResponse(w, "Failed to delete", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{"status": "deleted"}, http.StatusOK)
}

// DeleteAll soft deletes all notifications for the user
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
		utils.ErrorResponse(w, "Failed to delete all", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{"status": "all deleted"}, http.StatusOK)
}
