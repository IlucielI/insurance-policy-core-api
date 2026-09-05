package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/cache"
	"github.com/google/uuid"
)

type PolicyRepository struct {
	db    *sql.DB
	cache *cache.RedisClient
}

func NewPolicyRepository(db *sql.DB, cache *cache.RedisClient) *PolicyRepository {
	return &PolicyRepository{db: db, cache: cache}
}

func (r *PolicyRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Policy, int, error) {
	// Build cache key
	cacheKey := fmt.Sprintf("policies:user:%s:%d:%d", userID, limit, offset)
	
	// Try cache first (1 min TTL)
	if r.cache != nil {
		cached, err := r.cache.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var result struct {
				Policies []*domain.Policy `json:"policies"`
				Total    int              `json:"total"`
			}
			if json.Unmarshal([]byte(cached), &result) == nil {
				return result.Policies, result.Total, nil
			}
		}
	}
	
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
	
	// Cache for 1 minute
	if r.cache != nil {
		result := struct {
			Policies []*domain.Policy `json:"policies"`
			Total    int              `json:"total"`
		}{Policies: policies, Total: total}
		
		if cached, err := json.Marshal(result); err == nil {
			_ = r.cache.Set(ctx, cacheKey, string(cached), 1*time.Minute)
		}
	}

	return policies, total, nil
}

func (r *PolicyRepository) GetByUserIDWithFilters(ctx context.Context, userID, search, status, product, dateFrom, dateTo string, limit, offset int) ([]*domain.Policy, int, error) {
	whereClauses := []string{"p.user_id = $1"}
	args := []interface{}{userID}
	argPos := 2

	// Text search: policy_number or product_name (via join)
	if search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(p.policy_number ILIKE $%d OR pr.name ILIKE $%d)", argPos, argPos+1))
		args = append(args, "%"+search+"%", "%"+search+"%")
		argPos += 2
	}

	// Status filter
	if status != "" && status != "all" {
		whereClauses = append(whereClauses, fmt.Sprintf("p.status = $%d", argPos))
		args = append(args, status)
		argPos++
	}

	// Product filter
	if product != "" && product != "all" {
		whereClauses = append(whereClauses, fmt.Sprintf("p.product_id = $%d", argPos))
		args = append(args, product)
		argPos++
	}

	// Date range filter
	if dateFrom != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("p.issue_date >= $%d", argPos))
		args = append(args, dateFrom)
		argPos++
	}
	if dateTo != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("p.issue_date <= $%d", argPos))
		args = append(args, dateTo)
		argPos++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM policies p
		LEFT JOIN products pr ON p.product_id = pr.id
		%s
	`, whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch policies
	query := fmt.Sprintf(`
		SELECT p.id, p.policy_number, p.application_id, p.user_id, p.product_id, p.sum_assured,
			p.premium_amount, p.payment_frequency, p.status, p.issue_date, p.expiry_date,
			p.last_premium_paid_date, p.next_premium_due_date, p.created_at, p.updated_at
		FROM policies p
		LEFT JOIN products pr ON p.product_id = pr.id
		%s
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
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
	// Build cache key for policy details
	cacheKey := fmt.Sprintf("policy:details:%s", id)
	
	// Try cache first (30 min TTL for policy details)
	if r.cache != nil {
		cached, err := r.cache.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var policy domain.Policy
			if json.Unmarshal([]byte(cached), &policy) == nil {
				return &policy, nil
			}
		}
	}
	
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

	// Cache policy details for 30 minutes
	if r.cache != nil {
		if cached, err := json.Marshal(policy); err == nil {
			_ = r.cache.Set(ctx, cacheKey, string(cached), 30*time.Minute)
		}
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
	
	// Invalidate both user policy list and policy details cache on update
	if err == nil && r.cache != nil {
		cachePattern := fmt.Sprintf("policies:user:%s:*", policy.UserID)
		detailsKey := fmt.Sprintf("policy:details:%s", policy.ID)
		_ = r.cache.Delete(ctx, cachePattern, detailsKey)
	}
	
	return err
}
