package handlers

import (
	"net/http"

	"forum/middleware"
	nrepo "forum/repository/notification"
	"forum/utils"
)

type NotificationHandler struct{ Repo *nrepo.Repository }

func NewNotificationHandler(r *nrepo.Repository) *NotificationHandler {
	return &NotificationHandler{Repo: r}
}

func (h *NotificationHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	notes, err := h.Repo.GetByUser(user.ID)
	if err != nil {
		utils.ErrorResponse(w, "failed to load notifications", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, notes, http.StatusOK)
}

func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id := utils.GetLastPathParam(r)
	if id == "" {
		utils.ErrorResponse(w, "missing id", http.StatusBadRequest)
		return
	}
	n, err := h.Repo.GetByID(id)
	if err != nil {
		utils.ErrorResponse(w, "not found", http.StatusNotFound)
		return
	}
	if n.UserID != user.ID {
		utils.ErrorResponse(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := h.Repo.MarkRead(id, user.ID); err != nil {
		utils.ErrorResponse(w, "failed", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{"status": "read"}, http.StatusOK)
}

func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := h.Repo.MarkAllRead(user.ID); err != nil {
		utils.ErrorResponse(w, "failed", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{"status": "all read"}, http.StatusOK)
}

func (h *NotificationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id := utils.GetLastPathParam(r)
	if id == "" {
		utils.ErrorResponse(w, "missing id", http.StatusBadRequest)
		return
	}
	n, err := h.Repo.GetByID(id)
	if err != nil {
		utils.ErrorResponse(w, "not found", http.StatusNotFound)
		return
	}
	if n.UserID != user.ID {
		utils.ErrorResponse(w, "forbidden", http.StatusForbidden)
		return
	}
	if err := h.Repo.SoftDelete(id, user.ID); err != nil {
		utils.ErrorResponse(w, "failed", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{"status": "deleted"}, http.StatusOK)
}

func (h *NotificationHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.ErrorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if err := h.Repo.SoftDeleteAll(user.ID); err != nil {
		utils.ErrorResponse(w, "failed", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{"status": "deleted"}, http.StatusOK)
}
