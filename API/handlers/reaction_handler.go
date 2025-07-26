package handlers

import (
	"encoding/json"
	"net/http"

	"forum/middleware"
	"forum/models"
	"forum/repository"
	"forum/utils"
)

// ReactionHandler handles like/dislike reactions
type ReactionHandler struct {
	Repo             *repository.ReactionRepository
	PostRepo         *repository.PostRepository
	NotificationRepo *repository.NotificationRepository
}

func NewReactionHandler(repo *repository.ReactionRepository, pRepo *repository.PostRepository, nRepo *repository.NotificationRepository) *ReactionHandler {
	return &ReactionHandler{Repo: repo, PostRepo: pRepo, NotificationRepo: nRepo}
}

// React toggles a reaction on a post or comment for the authenticated user
func (h *ReactionHandler) CreateReact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetCurrentUser(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		TargetID     string `json:"target_id"`
		TargetType   string `json:"target_type"`
		ReactionType int    `json:"reaction_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TargetID == "" || (req.TargetType != "post" && req.TargetType != "comment") || req.ReactionType == 0 {
		utils.ErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.Repo.ToggleReaction(user.ID, req.TargetType, req.TargetID, req.ReactionType); err != nil {
		utils.ErrorResponse(w, "Failed to react", http.StatusInternalServerError)
		return
	}

	// After toggling, check if reaction exists to notify
	if req.TargetType == "post" {
		if ownerID, err := h.PostRepo.GetAuthorID(req.TargetID); err == nil && ownerID != user.ID {
			if t, _ := h.Repo.GetReactionType(user.ID, "post", req.TargetID); t != 0 {
				msg := user.Username + " reacted to your post"
				n := models.Notification{UserID: ownerID, PostID: &req.TargetID, Type: "reaction", Message: msg}
				h.NotificationRepo.Create(n)
			}
		}
	}

	var (
		reactions []models.ReactionWithUser
		err       error
	)
	if req.TargetType == "post" {
		reactions, err = h.Repo.GetReactionsByPostWithUser(req.TargetID)
	} else {
		reactions, err = h.Repo.GetReactionsByCommentWithUser(req.TargetID)
	}
	if err != nil {
		utils.ErrorResponse(w, "Failed to load reactions", http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(w, reactions, http.StatusOK)
}
