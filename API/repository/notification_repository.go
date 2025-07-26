package repository

import (
	"database/sql"
	"forum/models"
	"forum/utils"
	"time"
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
	_, err := r.db.Exec(`INSERT INTO notifications (notification_id, user_id, type, post_id, comment_id, created_at, read, deleted) VALUES (?, ?, ?, ?, ?, ?, 0, 0)`,
		n.ID, n.UserID, n.Type, n.PostID, n.CommentID, n.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *NotificationRepository) GetByUser(userID string) ([]models.Notification, error) {
	rows, err := r.db.Query(`SELECT notification_id, user_id, type, post_id, comment_id, read, deleted, created_at FROM notifications WHERE user_id = ? AND deleted = 0 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var n []models.Notification
	for rows.Next() {
		var no models.Notification
		if err := rows.Scan(&no.ID, &no.UserID, &no.Type, &no.PostID, &no.CommentID, &no.Read, &no.Deleted, &no.CreatedAt); err != nil {
			return nil, err
		}
		n = append(n, no)
	}
	return n, nil
}

func (r *NotificationRepository) MarkRead(id string) error {
	_, err := r.db.Exec(`UPDATE notifications SET read = 1 WHERE notification_id = ?`, id)
	return err
}

func (r *NotificationRepository) MarkAllRead(userID string) error {
	_, err := r.db.Exec(`UPDATE notifications SET read = 1 WHERE user_id = ? AND deleted = 0`, userID)
	return err
}

func (r *NotificationRepository) SoftDelete(id string) error {
	_, err := r.db.Exec(`UPDATE notifications SET deleted = 1 WHERE notification_id = ?`, id)
	return err
}

func (r *NotificationRepository) SoftDeleteAll(userID string) error {
	_, err := r.db.Exec(`UPDATE notifications SET deleted = 1 WHERE user_id = ?`, userID)
	return err
}
