package models

import "time"

// Notification represents a user notification
type Notification struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	ActorID   string     `json:"actor_id,omitempty"`
	PostID    *string    `json:"post_id,omitempty"`
	CommentID *string    `json:"comment_id,omitempty"`
	Type      string     `json:"type"`
	Message   *string    `json:"message"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
