package domain

import "time"

type AuditLog struct {
	ID          string                 `json:"id"`
	UserID      *string                `json:"user_id,omitempty"`
	Action      string                 `json:"action"`
	EntityType  string                 `json:"entity_type"`
	EntityID    string                 `json:"entity_id"`
	ChangesJSON map[string]interface{} `json:"changes_json,omitempty"`
	IPAddress   string                 `json:"ip_address,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	
	// Joined data (not in DB)
	UserName    string `json:"user_name,omitempty"`
	UserEmail   string `json:"user_email,omitempty"`
}
