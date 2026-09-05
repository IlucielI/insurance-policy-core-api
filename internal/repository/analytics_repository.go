package repository

import (
	"context"
	"database/sql"
	"time"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

// RevenueByMonth returns monthly premium revenue from approved applications
func (r *AnalyticsRepository) RevenueByMonth(ctx context.Context, months int) ([]MonthlyRevenue, error) {
	query := `
		SELECT 
			DATE_TRUNC('month', approved_at) AS month,
			COALESCE(SUM(premium_amount), 0) AS revenue,
			COUNT(*) AS policy_count
		FROM applications
		WHERE status = 'approved'
			AND approved_at >= NOW() - ($1 || ' months')::INTERVAL
		GROUP BY DATE_TRUNC('month', approved_at)
		ORDER BY month ASC
	`
	rows, err := r.db.QueryContext(ctx, query, months)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MonthlyRevenue
	for rows.Next() {
		var m MonthlyRevenue
		if err := rows.Scan(&m.Month, &m.Revenue, &m.PolicyCount); err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, nil
}

// ClaimsStatusDistribution returns count of claims grouped by status
func (r *AnalyticsRepository) ClaimsStatusDistribution(ctx context.Context) ([]StatusCount, error) {
	query := `
		SELECT status, COUNT(*) AS count
		FROM claims
		GROUP BY status
		ORDER BY count DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []StatusCount
	for rows.Next() {
		var sc StatusCount
		if err := rows.Scan(&sc.Status, &sc.Count); err != nil {
			return nil, err
		}
		results = append(results, sc)
	}
	return results, nil
}

// PolicyGrowthByMonth returns policy count growth by month
func (r *AnalyticsRepository) PolicyGrowthByMonth(ctx context.Context, months int) ([]MonthlyGrowth, error) {
	query := `
		SELECT 
			DATE_TRUNC('month', created_at) AS month,
			COUNT(*) AS new_policies,
			COALESCE(SUM(premium_amount), 0) AS premium
		FROM applications
		WHERE status = 'approved'
			AND created_at >= NOW() - ($1 || ' months')::INTERVAL
		GROUP BY DATE_TRUNC('month', created_at)
		ORDER BY month ASC
	`
	rows, err := r.db.QueryContext(ctx, query, months)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MonthlyGrowth
	for rows.Next() {
		var mg MonthlyGrowth
		if err := rows.Scan(&mg.Month, &mg.NewPolicies, &mg.Premium); err != nil {
			return nil, err
		}
		results = append(results, mg)
	}
	return results, nil
}

// TopProducts returns top products by application count
func (r *AnalyticsRepository) TopProducts(ctx context.Context, limit int) ([]ProductRank, error) {
	query := `
		SELECT 
			p.name AS product_name,
			p.category,
			COUNT(a.id) AS application_count,
			COALESCE(SUM(a.premium_amount), 0) AS total_premium
		FROM applications a
		JOIN products p ON a.product_id = p.id
		WHERE a.status = 'approved'
		GROUP BY p.id, p.name, p.category
		ORDER BY application_count DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ProductRank
	for rows.Next() {
		var pr ProductRank
		if err := rows.Scan(&pr.ProductName, &pr.Category, &pr.ApplicationCount, &pr.TotalPremium); err != nil {
			return nil, err
		}
		results = append(results, pr)
	}
	return results, nil
}

// SummaryStats returns overall dashboard summary
func (r *AnalyticsRepository) SummaryStats(ctx context.Context) (*SummaryStats, error) {
	s := &SummaryStats{}

	// Total approved applications
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM applications WHERE status = 'approved'`).Scan(&s.TotalApproved)
	if err != nil {
		return nil, err
	}

	// Total premium YTD
	err = r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(premium_amount), 0) FROM applications 
		 WHERE status = 'approved' AND approved_at >= DATE_TRUNC('year', NOW())`).Scan(&s.TotalPremiumYTD)
	if err != nil {
		return nil, err
	}

	// Total claims submitted
	err = r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM claims`).Scan(&s.TotalClaims)
	if err != nil {
		return nil, err
	}

	// Total claims paid amount
	err = r.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(approved_amount), 0) FROM claims WHERE status = 'paid'`).Scan(&s.TotalClaimsPaid)
	if err != nil {
		return nil, err
	}

	// Active policies count
	err = r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM policies WHERE status = 'active'`).Scan(&s.ActivePolicies)
	if err != nil {
		return nil, err
	}

	return s, nil
}

type MonthlyRevenue struct {
	Month       time.Time `json:"month"`
	Revenue     int64     `json:"revenue"`
	PolicyCount int64     `json:"policy_count"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type MonthlyGrowth struct {
	Month       time.Time `json:"month"`
	NewPolicies int64     `json:"new_policies"`
	Premium     int64     `json:"premium"`
}

type ProductRank struct {
	ProductName      string `json:"product_name"`
	Category         string `json:"category"`
	ApplicationCount int64  `json:"application_count"`
	TotalPremium     int64  `json:"total_premium"`
}

type SummaryStats struct {
	TotalApproved   int64 `json:"total_approved"`
	TotalPremiumYTD int64 `json:"total_premium_ytd"`
	TotalClaims     int64 `json:"total_claims"`
	TotalClaimsPaid int64 `json:"total_claims_paid"`
	ActivePolicies  int64 `json:"active_policies"`
}