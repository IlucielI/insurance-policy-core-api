package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/google/uuid"
)

type PolicyRepository struct {
	db *sql.DB
}

func NewPolicyRepository(db *sql.DB) *PolicyRepository {
	return &PolicyRepository{db: db}
}

func (r *PolicyRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Policy, int, error) {
	// Count total
	var total int
	countQuery := `SELECT COUNT(*) FROM policies WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch policies
	query := `
		SELECT id, policy_number, application_id, user_id, product_id, sum_assured, 
			premium_amount, payment_frequency, status, issue_date, expiry_date,
			last_premium_paid_date, next_premium_due_date, created_at, updated_at
		FROM policies
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	policies := []*domain.Policy{}
	for rows.Next() {
		policy := &domain.Policy{}
		err := rows.Scan(
			&policy.ID, &policy.PolicyNumber, &policy.ApplicationID, &policy.UserID,
			&policy.ProductID, &policy.SumAssured, &policy.PremiumAmount,
			&policy.PaymentFrequency, &policy.Status, &policy.IssueDate,
			&policy.ExpiryDate, &policy.LastPremiumPaidDate, &policy.NextPremiumDueDate,
			&policy.CreatedAt, &policy.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		policies = append(policies, policy)
	}

	return policies, total, nil
}

func (r *PolicyRepository) GetByID(ctx context.Context, id string) (*domain.Policy, error) {
	policy := &domain.Policy{}
	query := `
		SELECT id, policy_number, application_id, user_id, product_id, sum_assured, 
			premium_amount, payment_frequency, status, issue_date, expiry_date,
			last_premium_paid_date, next_premium_due_date, created_at, updated_at
		FROM policies
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&policy.ID, &policy.PolicyNumber, &policy.ApplicationID, &policy.UserID,
		&policy.ProductID, &policy.SumAssured, &policy.PremiumAmount,
		&policy.PaymentFrequency, &policy.Status, &policy.IssueDate,
		&policy.ExpiryDate, &policy.LastPremiumPaidDate, &policy.NextPremiumDueDate,
		&policy.CreatedAt, &policy.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (r *PolicyRepository) CreateEndorsement(ctx context.Context, endorsement *domain.PolicyEndorsement) error {
	endorsement.ID = uuid.New().String()
	endorsement.CreatedAt = time.Now()
	endorsement.UpdatedAt = time.Now()

	query := `
		INSERT INTO policy_endorsements (id, endorsement_number, policy_id, endorsement_type,
			description, effective_date, old_values, new_values, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.ExecContext(ctx, query,
		endorsement.ID, endorsement.EndorsementNumber, endorsement.PolicyID,
		endorsement.EndorsementType, endorsement.Description, endorsement.EffectiveDate,
		nil, nil, endorsement.Status, endorsement.CreatedAt, endorsement.UpdatedAt,
	)
	return err
}

func (r *PolicyRepository) Update(ctx context.Context, policy *domain.Policy) error {
	policy.UpdatedAt = time.Now()

	query := `
		UPDATE policies
		SET status = $1, last_premium_paid_date = $2, next_premium_due_date = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := r.db.ExecContext(ctx, query,
		policy.Status, policy.LastPremiumPaidDate, policy.NextPremiumDueDate,
		policy.UpdatedAt, policy.ID,
	)
	return err
}
