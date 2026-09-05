package repository

import (
	"context"
	"database/sql"
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
