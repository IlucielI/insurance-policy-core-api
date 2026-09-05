package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/model"
	"github.com/google/uuid"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(req *model.NotificationCreateRequest) (*model.Notification, error) {
	notification := &model.Notification{
		ID:            uuid.New().String(),
		UserID:        req.UserID,
		Type:          req.Type,
		Title:         req.Title,
		Message:       req.Message,
		ReferenceID:   req.ReferenceID,
		ReferenceType: req.ReferenceType,
		IsRead:        false,
		CreatedAt:     time.Now(),
	}

	query := `
		INSERT INTO notifications (id, user_id, type, title, message, reference_id, reference_type, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, user_id, type, title, message, reference_id, reference_type, is_read, created_at
	`

	err := r.db.QueryRow(
		query,
		notification.ID,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.ReferenceID,
		notification.ReferenceType,
		notification.IsRead,
		notification.CreatedAt,
	).Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Type,
		&notification.Title,
		&notification.Message,
		&notification.ReferenceID,
		&notification.ReferenceType,
		&notification.IsRead,
		&notification.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return notification, nil
}

func (r *NotificationRepository) GetByUserID(userID string, limit, offset int) ([]*model.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, reference_id, reference_type, is_read, created_at, read_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*model.Notification
	for rows.Next() {
		n := &model.Notification{}
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Type,
			&n.Title,
			&n.Message,
			&n.ReferenceID,
			&n.ReferenceType,
			&n.IsRead,
			&n.CreatedAt,
			&n.ReadAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

func (r *NotificationRepository) MarkAsRead(notificationID string) error {
	query := `
		UPDATE notifications
		SET is_read = TRUE, read_at = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(query, time.Now(), notificationID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	return nil
}

func (r *NotificationRepository) MarkAllAsRead(userID string) error {
	query := `
		UPDATE notifications
		SET is_read = TRUE, read_at = $1
		WHERE user_id = $2 AND is_read = FALSE
	`

	_, err := r.db.Exec(query, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	return nil
}

func (r *NotificationRepository) GetUnreadCount(userID string) (int, error) {
	query := `
		SELECT COUNT(*) FROM notifications
		WHERE user_id = $1 AND is_read = FALSE
	`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

func (r *NotificationRepository) Delete(notificationID string) error {
	query := `DELETE FROM notifications WHERE id = $1`
	_, err := r.db.Exec(query, notificationID)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}
