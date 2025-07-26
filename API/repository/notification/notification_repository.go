package notification

import (
	"database/sql"
	"time"

	"forum/models"
	"forum/utils"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(n models.Notification) (*models.Notification, error) {
	n.ID = utils.GenerateUUID()
	n.CreatedAt = time.Now()
	_, err := r.db.Exec(`INSERT INTO notifications (notification_id, user_id, actor_id, post_id, comment_id, type, message, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		n.ID, n.UserID, n.ActorID, n.PostID, n.CommentID, n.Type, n.Message, n.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *Repository) GetByUser(userID string, limit, offset int) ([]models.Notification, error) {
	rows, err := r.db.Query(`SELECT n.notification_id, n.user_id, n.actor_id, n.post_id, n.comment_id, n.type, n.message,
                p.title, c.content,
                n.created_at, n.read_at, n.updated_at
                FROM notifications n
                LEFT JOIN posts p ON n.post_id = p.post_id
                LEFT JOIN comments c ON n.comment_id = c.comment_id
                WHERE n.user_id = ?
                ORDER BY n.created_at DESC
                LIMIT ? OFFSET ?`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ns []models.Notification
	for rows.Next() {
		var n models.Notification
		var title, content *string
		if err := rows.Scan(&n.ID, &n.UserID, &n.ActorID, &n.PostID, &n.CommentID, &n.Type, &n.Message, &title, &content, &n.CreatedAt, &n.ReadAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		n.PostTitle = title
		if content != nil {
			snippet := *content
			if len(snippet) > 50 {
				snippet = snippet[:50]
			}
			n.CommentSnippet = &snippet
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func (r *Repository) MarkRead(id string) error {
	_, err := r.db.Exec(`UPDATE notifications SET read_at = ? WHERE notification_id = ?`, time.Now(), id)
	return err
}

func (r *Repository) MarkAllRead(userID string) error {
	_, err := r.db.Exec(`UPDATE notifications SET read_at = ? WHERE user_id = ?`, time.Now(), userID)
	return err
}

func (r *Repository) SoftDelete(id string) error {
	_, err := r.db.Exec(`UPDATE notifications SET message = NULL, updated_at = ? WHERE notification_id = ?`, time.Now(), id)
	return err
}

func (r *Repository) SoftDeleteAll(userID string) error {
	_, err := r.db.Exec(`UPDATE notifications SET message = NULL, updated_at = ? WHERE user_id = ?`, time.Now(), userID)
	return err
}
