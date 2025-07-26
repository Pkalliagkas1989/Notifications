package handlers

import (
	"encoding/json"
	"net/http"

	"forum/middleware"
	"forum/models"
	"forum/repository"
	nrepo "forum/repository/notification"
	"forum/repository/user"
	"forum/utils"
)

// ReactionHandler handles like/dislike reactions
type ReactionHandler struct {
	Repo             *repository.ReactionRepository
	PostRepo         *repository.PostRepository
	NotificationRepo *nrepo.Repository
	UserRepo         *user.UserRepository
}

func NewReactionHandler(repo *repository.ReactionRepository, postRepo *repository.PostRepository, nRepo *nrepo.Repository, uRepo *user.UserRepository) *ReactionHandler {
	return &ReactionHandler{Repo: repo, PostRepo: postRepo, NotificationRepo: nRepo, UserRepo: uRepo}
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

	newType, _ := h.Repo.GetReaction(user.ID, req.TargetType, req.TargetID)
	if req.TargetType == "post" && newType != 0 {
		if ownerID, err := h.PostRepo.GetPostOwner(req.TargetID); err == nil && ownerID != user.ID {
			if actor, err2 := h.UserRepo.GetByID(user.ID); err2 == nil {
				action := "liked"
				if newType == 2 {
					action = "disliked"
				}
				msg := actor.Username + " " + action + " your post"
				n := models.Notification{UserID: ownerID, ActorID: user.ID, PostID: &req.TargetID, Type: "reaction", Message: &msg}
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
