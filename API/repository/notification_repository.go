package repository

import (
	"database/sql"
	"time"

	"forum/models"
	"forum/utils"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(n models.Notification) (*models.Notification, error) {
	n.ID = utils.GenerateUUID()
	n.CreatedAt = time.Now()
	_, err := r.db.Exec(`INSERT INTO notifications (notification_id, user_id, post_id, comment_id, type, message, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		n.ID, n.UserID, n.PostID, n.CommentID, n.Type, n.Message, n.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *NotificationRepository) GetByUser(userID string) ([]models.Notification, error) {
	rows, err := r.db.Query(`SELECT notification_id, user_id, post_id, comment_id, type, message, created_at, read_at, deleted_at FROM notifications WHERE user_id = ? AND deleted_at IS NULL ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ns []models.Notification
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.PostID, &n.CommentID, &n.Type, &n.Message, &n.CreatedAt, &n.ReadAt, &n.DeletedAt); err != nil {
			return nil, err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func (r *NotificationRepository) MarkRead(id string) error {
	_, err := r.db.Exec(`UPDATE notifications SET read_at = ? WHERE notification_id = ?`, time.Now(), id)
	return err
}

func (r *NotificationRepository) MarkAllRead(userID string) error {
	_, err := r.db.Exec(`UPDATE notifications SET read_at = ? WHERE user_id = ? AND read_at IS NULL`, time.Now(), userID)
	return err
}

func (r *NotificationRepository) SoftDelete(id string) error {
	_, err := r.db.Exec(`UPDATE notifications SET deleted_at = ? WHERE notification_id = ?`, time.Now(), id)
	return err
}

func (r *NotificationRepository) SoftDeleteAll(userID string) error {
	_, err := r.db.Exec(`UPDATE notifications SET deleted_at = ? WHERE user_id = ? AND deleted_at IS NULL`, time.Now(), userID)
	return err
}
