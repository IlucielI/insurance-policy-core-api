package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/google/uuid"
)

type ClaimRepository struct {
	db *sql.DB
}

func NewClaimRepository(db *sql.DB) *ClaimRepository {
	return &ClaimRepository{db: db}
}

func (r *ClaimRepository) Create(ctx context.Context, claim *domain.Claim) error {
	claim.ID = uuid.New().String()
	claim.SubmittedAt = time.Now()
	claim.CreatedAt = time.Now()
	claim.UpdatedAt = time.Now()

	query := `
		INSERT INTO claims (id, claim_number, policy_id, user_id, claim_type, claim_amount,
			incident_date, incident_description, status, submitted_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		claim.ID, claim.ClaimNumber, claim.PolicyID, claim.UserID, claim.ClaimType,
		claim.ClaimAmount, claim.IncidentDate, claim.IncidentDescription, claim.Status,
		claim.SubmittedAt, claim.CreatedAt, claim.UpdatedAt,
	)
	return err
}

func (r *ClaimRepository) GetByID(ctx context.Context, id string) (*domain.Claim, error) {
	claim := &domain.Claim{}
	query := `
		SELECT id, claim_number, policy_id, user_id, claim_type, claim_amount,
			incident_date, incident_description, status, reviewer_id, reviewer_notes,
			rejection_reason, approved_amount, submitted_at, reviewed_at, approved_at,
			paid_at, created_at, updated_at
		FROM claims
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&claim.ID, &claim.ClaimNumber, &claim.PolicyID, &claim.UserID, &claim.ClaimType,
		&claim.ClaimAmount, &claim.IncidentDate, &claim.IncidentDescription, &claim.Status,
		&claim.ReviewerID, &claim.ReviewerNotes, &claim.RejectionReason, &claim.ApprovedAmount,
		&claim.SubmittedAt, &claim.ReviewedAt, &claim.ApprovedAt, &claim.PaidAt,
		&claim.CreatedAt, &claim.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return claim, nil
}

func (r *ClaimRepository) AddDocument(ctx context.Context, doc *domain.ClaimDocument) error {
	doc.ID = uuid.New().String()
	doc.UploadedAt = time.Now()

	query := `
		INSERT INTO claim_documents (id, claim_id, document_type, file_name, file_path,
			file_size, mime_type, uploaded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		doc.ID, doc.ClaimID, doc.DocumentType, doc.FileName, doc.FilePath,
		doc.FileSize, doc.MimeType, doc.UploadedAt,
	)
	return err
}

func (r *ClaimRepository) AddTimelineEntry(ctx context.Context, entry *domain.ClaimTimeline) error {
	entry.ID = uuid.New().String()
	entry.CreatedAt = time.Now()

	query := `
		INSERT INTO claim_timeline (id, claim_id, action, description, actor_id, actor_name, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		entry.ID, entry.ClaimID, entry.Action, entry.Description, entry.ActorID,
		entry.ActorName, entry.CreatedAt,
	)
	return err
}

func (r *ClaimRepository) GetTimeline(ctx context.Context, claimID string) ([]*domain.ClaimTimeline, error) {
	query := `
		SELECT id, claim_id, action, description, actor_id, actor_name, created_at
		FROM claim_timeline
		WHERE claim_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, claimID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	timeline := []*domain.ClaimTimeline{}
	for rows.Next() {
		entry := &domain.ClaimTimeline{}
		err := rows.Scan(
			&entry.ID, &entry.ClaimID, &entry.Action, &entry.Description,
			&entry.ActorID, &entry.ActorName, &entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		timeline = append(timeline, entry)
	}

	return timeline, nil
}

func (r *ClaimRepository) Update(ctx context.Context, claim *domain.Claim) error {
	claim.UpdatedAt = time.Now()

	query := `
		UPDATE claims
		SET status = $1, reviewer_id = $2, reviewer_notes = $3, rejection_reason = $4,
			approved_amount = $5, reviewed_at = $6, approved_at = $7, paid_at = $8, updated_at = $9
		WHERE id = $10
	`

	_, err := r.db.ExecContext(ctx, query,
		claim.Status, claim.ReviewerID, claim.ReviewerNotes, claim.RejectionReason,
		claim.ApprovedAmount, claim.ReviewedAt, claim.ApprovedAt, claim.PaidAt,
		claim.UpdatedAt, claim.ID,
	)
	return err
}

// ListClaimsWithFilters returns claims with search and filter support
func (r *ClaimRepository) ListClaimsWithFilters(ctx context.Context, search, status, claimType, dateFrom, dateTo, amountMin, amountMax string, limit, offset int) ([]*domain.Claim, int, error) {
	where := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		where = append(where, fmt.Sprintf("(c.claim_number ILIKE $%d OR c.incident_description ILIKE $%d)", argIdx, argIdx+1))
		args = append(args, "%"+search+"%", "%"+search+"%")
		argIdx += 2
	}
	if status != "" {
		where = append(where, fmt.Sprintf("c.status = $%d", argIdx))
		args = append(args, status)
		argIdx++
	}
	if claimType != "" {
		where = append(where, fmt.Sprintf("c.claim_type = $%d", argIdx))
		args = append(args, claimType)
		argIdx++
	}
	if dateFrom != "" {
		where = append(where, fmt.Sprintf("c.incident_date >= $%d", argIdx))
		args = append(args, dateFrom)
		argIdx++
	}
	if dateTo != "" {
		where = append(where, fmt.Sprintf("c.incident_date <= $%d", argIdx))
		args = append(args, dateTo)
		argIdx++
	}
	if amountMin != "" {
		where = append(where, fmt.Sprintf("c.claim_amount >= $%d", argIdx))
		args = append(args, amountMin)
		argIdx++
	}
	if amountMax != "" {
		where = append(where, fmt.Sprintf("c.claim_amount <= $%d", argIdx))
		args = append(args, amountMax)
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")

	// Count
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM claims c WHERE %s", whereClause)
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch
	query := fmt.Sprintf(`
		SELECT c.id, c.claim_number, c.policy_id, c.user_id, c.claim_type, c.claim_amount,
			c.incident_date, c.incident_description, c.status, c.reviewer_id, c.reviewer_notes,
			c.rejection_reason, c.approved_amount, c.submitted_at, c.reviewed_at, c.approved_at,
			c.paid_at, c.created_at, c.updated_at
		FROM claims c
		WHERE %s
		ORDER BY c.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	claims := []*domain.Claim{}
	for rows.Next() {
		c := &domain.Claim{}
		err := rows.Scan(
			&c.ID, &c.ClaimNumber, &c.PolicyID, &c.UserID, &c.ClaimType, &c.ClaimAmount,
			&c.IncidentDate, &c.IncidentDescription, &c.Status, &c.ReviewerID, &c.ReviewerNotes,
			&c.RejectionReason, &c.ApprovedAmount, &c.SubmittedAt, &c.ReviewedAt, &c.ApprovedAt,
			&c.PaidAt, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		claims = append(claims, c)
	}
	return claims, total, nil
}

// ListAllClaimsForExport returns all claims without pagination for report export
func (r *ClaimRepository) ListAllClaimsForExport(ctx context.Context, status, claimType string) ([]*domain.Claim, error) {
	where := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if status != "" {
		where = append(where, fmt.Sprintf("c.status = $%d", argIdx))
		args = append(args, status)
		argIdx++
	}
	if claimType != "" {
		where = append(where, fmt.Sprintf("c.claim_type = $%d", argIdx))
		args = append(args, claimType)
		argIdx++
	}

	whereClause := strings.Join(where, " AND ")
	query := fmt.Sprintf(`
		SELECT c.id, c.claim_number, c.policy_id, c.user_id, c.claim_type, c.claim_amount,
			c.incident_date, c.incident_description, c.status, c.reviewer_id, c.reviewer_notes,
			c.rejection_reason, c.approved_amount, c.submitted_at, c.reviewed_at, c.approved_at,
			c.paid_at, c.created_at, c.updated_at,
			COALESCE(p.policy_number, '') as policy_number
		FROM claims c
		LEFT JOIN policies p ON c.policy_id = p.id
		WHERE %s
		ORDER BY c.created_at DESC
	`, whereClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	claims := []*domain.Claim{}
	for rows.Next() {
		c := &domain.Claim{}
		policy := &domain.Policy{}
		err := rows.Scan(
			&c.ID, &c.ClaimNumber, &c.PolicyID, &c.UserID, &c.ClaimType, &c.ClaimAmount,
			&c.IncidentDate, &c.IncidentDescription, &c.Status, &c.ReviewerID, &c.ReviewerNotes,
			&c.RejectionReason, &c.ApprovedAmount, &c.SubmittedAt, &c.ReviewedAt, &c.ApprovedAt,
			&c.PaidAt, &c.CreatedAt, &c.UpdatedAt, &policy.PolicyNumber,
		)
		if err != nil {
			return nil, err
		}
		c.Policy = policy
		claims = append(claims, c)
	}
	return claims, nil
}

func (r *ClaimRepository) ListWithFilters(ctx context.Context, search, userID, status, claimType, dateFrom, dateTo string, limit, offset int) ([]*domain.Claim, int, error) {
	whereClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(c.claim_number ILIKE $%d OR c.incident_description ILIKE $%d OR u.email ILIKE $%d OR u.full_name ILIKE $%d)", argPos, argPos+1, argPos+2, argPos+3))
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
		argPos += 4
	}

	if userID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.user_id = $%d", argPos))
		args = append(args, userID)
		argPos++
	}

	if status != "" && status != "all" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.status = $%d", argPos))
		args = append(args, status)
		argPos++
	}

	if claimType != "" && claimType != "all" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.claim_type = $%d", argPos))
		args = append(args, claimType)
		argPos++
	}

	if dateFrom != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.incident_date >= $%d", argPos))
		args = append(args, dateFrom)
		argPos++
	}
	if dateTo != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("c.incident_date <= $%d", argPos))
		args = append(args, dateTo)
		argPos++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM claims c LEFT JOIN users u ON c.user_id = u.id %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT c.id, c.claim_number, c.policy_id, c.user_id, c.claim_type, c.claim_amount,
			c.incident_date, c.incident_description, c.status, c.reviewer_id, c.reviewer_notes,
			c.rejection_reason, c.approved_amount, c.submitted_at, c.reviewed_at, c.approved_at,
			c.paid_at, c.created_at, c.updated_at
		FROM claims c
		LEFT JOIN users u ON c.user_id = u.id
		%s
		ORDER BY c.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	claims := []*domain.Claim{}
	for rows.Next() {
		claim := &domain.Claim{}
		err := rows.Scan(
			&claim.ID, &claim.ClaimNumber, &claim.PolicyID, &claim.UserID, &claim.ClaimType,
			&claim.ClaimAmount, &claim.IncidentDate, &claim.IncidentDescription, &claim.Status,
			&claim.ReviewerID, &claim.ReviewerNotes, &claim.RejectionReason, &claim.ApprovedAmount,
			&claim.SubmittedAt, &claim.ReviewedAt, &claim.ApprovedAt, &claim.PaidAt,
			&claim.CreatedAt, &claim.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		claims = append(claims, claim)
	}

	return claims, total, nil
}
