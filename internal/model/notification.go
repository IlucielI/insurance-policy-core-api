package model

import (
	"time"
)

type NotificationType string

const (
	NotificationPolicyApproved  NotificationType = "policy_approved"
	NotificationClaimUpdated    NotificationType = "claim_updated"
	NotificationPaymentConfirmed NotificationType = "payment_confirmed"
)

type Notification struct {
	ID            string           `json:"id" db:"id"`
	UserID        string           `json:"user_id" db:"user_id"`
	Type          NotificationType `json:"type" db:"type"`
	Title         string           `json:"title" db:"title"`
	Message       string           `json:"message" db:"message"`
	ReferenceID   *string          `json:"reference_id,omitempty" db:"reference_id"`
	ReferenceType *string          `json:"reference_type,omitempty" db:"reference_type"`
	IsRead        bool             `json:"is_read" db:"is_read"`
	CreatedAt     time.Time        `json:"created_at" db:"created_at"`
	ReadAt        *time.Time       `json:"read_at,omitempty" db:"read_at"`
}

type NotificationCreateRequest struct {
	UserID        string           `json:"user_id"`
	Type          NotificationType `json:"type"`
	Title         string           `json:"title"`
	Message       string           `json:"message"`
	ReferenceID   *string          `json:"reference_id,omitempty"`
	ReferenceType *string          `json:"reference_type,omitempty"`
}
