package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
)

type AuditRepository struct {
	db *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// Create creates a new audit log entry
func (r *AuditRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	log.CreatedAt = time.Now()

	changesJSON, err := json.Marshal(log.ChangesJSON)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO audit_logs (user_id, action, entity_type, entity_id, changes_json, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	return r.db.QueryRowContext(ctx, query,
		log.UserID,
		log.Action,
		log.EntityType,
		log.EntityID,
		changesJSON,
		log.IPAddress,
		log.UserAgent,
		log.CreatedAt,
	).Scan(&log.ID)
}

// List retrieves audit logs with filters
func (r *AuditRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*domain.AuditLog, int, error) {
	var whereConditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE conditions
	if userID, ok := filters["user_id"].(string); ok && userID != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("al.user_id = $%d", argIndex))
		args = append(args, userID)
		argIndex++
	}

	if action, ok := filters["action"].(string); ok && action != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("al.action = $%d", argIndex))
		args = append(args, action)
		argIndex++
	}

	if entityType, ok := filters["entity_type"].(string); ok && entityType != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("al.entity_type = $%d", argIndex))
		args = append(args, entityType)
		argIndex++
	}

	if entityID, ok := filters["entity_id"].(string); ok && entityID != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("al.entity_id = $%d", argIndex))
		args = append(args, entityID)
		argIndex++
	}

	if dateFrom, ok := filters["date_from"].(string); ok && dateFrom != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("al.created_at >= $%d", argIndex))
		args = append(args, dateFrom)
		argIndex++
	}

	if dateTo, ok := filters["date_to"].(string); ok && dateTo != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("al.created_at <= $%d", argIndex))
		args = append(args, dateTo)
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM audit_logs al
		%s
	`, whereClause)

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results with user join
	args = append(args, limit, offset)
	query := fmt.Sprintf(`
		SELECT 
			al.id, al.user_id, al.action, al.entity_type, al.entity_id, 
			al.changes_json, al.ip_address, al.user_agent, al.created_at,
			u.full_name, u.email
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		%s
		ORDER BY al.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		log := &domain.AuditLog{}
		var changesJSON []byte
		var userName, userEmail sql.NullString

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Action,
			&log.EntityType,
			&log.EntityID,
			&changesJSON,
			&log.IPAddress,
			&log.UserAgent,
			&log.CreatedAt,
			&userName,
			&userEmail,
		)
		if err != nil {
			return nil, 0, err
		}

		// Parse changes JSON
		if len(changesJSON) > 0 {
			json.Unmarshal(changesJSON, &log.ChangesJSON)
		}

		// Set user info
		if userName.Valid {
			log.UserName = userName.String
		}
		if userEmail.Valid {
			log.UserEmail = userEmail.String
		}

		logs = append(logs, log)
	}

	return logs, total, rows.Err()
}

// GetByEntityID retrieves audit logs for a specific entity
func (r *AuditRepository) GetByEntityID(ctx context.Context, entityType, entityID string, limit int) ([]*domain.AuditLog, error) {
	query := `
		SELECT 
			al.id, al.user_id, al.action, al.entity_type, al.entity_id, 
			al.changes_json, al.ip_address, al.user_agent, al.created_at,
			u.full_name, u.email
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE al.entity_type = $1 AND al.entity_id = $2
		ORDER BY al.created_at DESC
		LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, query, entityType, entityID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		log := &domain.AuditLog{}
		var changesJSON []byte
		var userName, userEmail sql.NullString

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Action,
			&log.EntityType,
			&log.EntityID,
			&changesJSON,
			&log.IPAddress,
			&log.UserAgent,
			&log.CreatedAt,
			&userName,
			&userEmail,
		)
		if err != nil {
			return nil, err
		}

		if len(changesJSON) > 0 {
			json.Unmarshal(changesJSON, &log.ChangesJSON)
		}

		if userName.Valid {
			log.UserName = userName.String
		}
		if userEmail.Valid {
			log.UserEmail = userEmail.String
		}

		logs = append(logs, log)
	}

	return logs, rows.Err()
}

// GetByID retrieves a single audit log entry by ID
func (r *AuditRepository) GetByID(ctx context.Context, id string) (*domain.AuditLog, error) {
	query := `
		SELECT 
			al.id, al.user_id, al.action, al.entity_type, al.entity_id, 
			al.changes_json, al.ip_address, al.user_agent, al.created_at,
			u.full_name, u.email
		FROM audit_logs al
		LEFT JOIN users u ON al.user_id = u.id
		WHERE al.id = $1
	`

	log := &domain.AuditLog{}
	var changesJSON []byte
	var userName, userEmail sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&log.ID,
		&log.UserID,
		&log.Action,
		&log.EntityType,
		&log.EntityID,
		&changesJSON,
		&log.IPAddress,
		&log.UserAgent,
		&log.CreatedAt,
		&userName,
		&userEmail,
	)
	if err != nil {
		return nil, err
	}

	if len(changesJSON) > 0 {
		json.Unmarshal(changesJSON, &log.ChangesJSON)
	}
	if userName.Valid {
		log.UserName = userName.String
	}
	if userEmail.Valid {
		log.UserEmail = userEmail.String
	}

	return log, nil
}
