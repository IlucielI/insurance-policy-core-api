package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
)

type NotificationPreferencesRepository struct {
	db *sql.DB
}

func NewNotificationPreferencesRepository(db *sql.DB) *NotificationPreferencesRepository {
	return &NotificationPreferencesRepository{db: db}
}

func (r *NotificationPreferencesRepository) GetByUserID(userID string) (*domain.NotificationPreferences, error) {
	query := `
		SELECT id, user_id, promotional_emails, policy_update_emails, 
		       claim_notification_emails, newsletter_emails, created_at, updated_at
		FROM notification_preferences
		WHERE user_id = $1
	`

	var prefs domain.NotificationPreferences
	err := r.db.QueryRow(query, userID).Scan(
		&prefs.ID,
		&prefs.UserID,
		&prefs.PromotionalEmails,
		&prefs.PolicyUpdateEmails,
		&prefs.ClaimNotificationEmails,
		&prefs.NewsletterEmails,
		&prefs.CreatedAt,
		&prefs.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Create default preferences if not exists
		return r.CreateDefault(userID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get notification preferences: %w", err)
	}

	return &prefs, nil
}

func (r *NotificationPreferencesRepository) CreateDefault(userID string) (*domain.NotificationPreferences, error) {
	query := `
		INSERT INTO notification_preferences 
		(user_id, promotional_emails, policy_update_emails, claim_notification_emails, newsletter_emails)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, promotional_emails, policy_update_emails, 
		          claim_notification_emails, newsletter_emails, created_at, updated_at
	`

	var prefs domain.NotificationPreferences
	err := r.db.QueryRow(query, userID, true, true, true, true).Scan(
		&prefs.ID,
		&prefs.UserID,
		&prefs.PromotionalEmails,
		&prefs.PolicyUpdateEmails,
		&prefs.ClaimNotificationEmails,
		&prefs.NewsletterEmails,
		&prefs.CreatedAt,
		&prefs.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create default preferences: %w", err)
	}

	return &prefs, nil
}

func (r *NotificationPreferencesRepository) Update(userID string, prefs *domain.NotificationPreferences) error {
	query := `
		UPDATE notification_preferences
		SET promotional_emails = $1,
		    policy_update_emails = $2,
		    claim_notification_emails = $3,
		    newsletter_emails = $4,
		    updated_at = $5
		WHERE user_id = $6
	`

	result, err := r.db.Exec(query,
		prefs.PromotionalEmails,
		prefs.PolicyUpdateEmails,
		prefs.ClaimNotificationEmails,
		prefs.NewsletterEmails,
		time.Now(),
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to update notification preferences: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no preferences found for user")
	}

	return nil
}

func (r *NotificationPreferencesRepository) CheckPreference(userID, preferenceType string) (bool, error) {
	var enabled bool
	var query string

	switch preferenceType {
	case "promotional":
		query = "SELECT promotional_emails FROM notification_preferences WHERE user_id = $1"
	case "policy_update":
		query = "SELECT policy_update_emails FROM notification_preferences WHERE user_id = $1"
	case "claim_notification":
		query = "SELECT claim_notification_emails FROM notification_preferences WHERE user_id = $1"
	case "newsletter":
		query = "SELECT newsletter_emails FROM notification_preferences WHERE user_id = $1"
	default:
		return false, fmt.Errorf("invalid preference type: %s", preferenceType)
	}

	err := r.db.QueryRow(query, userID).Scan(&enabled)
	if err == sql.ErrNoRows {
		return true, nil // Default to enabled if no preferences found
	}

	if err != nil {
		return false, fmt.Errorf("failed to check preference: %w", err)
	}

	return enabled, nil
}
