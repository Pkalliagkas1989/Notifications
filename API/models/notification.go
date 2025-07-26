package models

import "time"

// Notification represents a user notification
// Type can be: comment_created, comment_edited, comment_deleted, reaction
// PostID or CommentID reference the related object
// DeletedAt is used for soft deletion

type Notification struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	PostID    *string    `json:"post_id,omitempty"`
	CommentID *string    `json:"comment_id,omitempty"`
	Type      string     `json:"type"`
	Message   string     `json:"message"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	DeletedAt *time.Time `json:"-"`
}
