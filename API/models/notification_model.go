package models

import "time"

// Notification represents a user notification
// Type can be: comment, like, dislike, comment_edit, comment_delete
// Deleted implements soft-delete behavior
// Read indicates if the user has seen the notification

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`
	PostID    string    `json:"post_id"`
	CommentID *string   `json:"comment_id,omitempty"`
	Read      bool      `json:"read"`
	Deleted   bool      `json:"deleted"`
	CreatedAt time.Time `json:"created_at"`
}
