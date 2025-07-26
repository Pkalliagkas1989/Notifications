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

func (r *Repository) GetByUser(userID string) ([]models.Notification, error) {
	rows, err := r.db.Query(`SELECT notification_id, user_id, actor_id, post_id, comment_id, type, message, created_at, read_at, updated_at FROM notifications WHERE user_id = ? ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ns []models.Notification
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.ActorID, &n.PostID, &n.CommentID, &n.Type, &n.Message, &n.CreatedAt, &n.ReadAt, &n.UpdatedAt); err != nil {
			return nil, err
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
