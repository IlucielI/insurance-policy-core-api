package repository

import (
	"context"
	"database/sql"

	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/report"
)

// ReportRepository handles data queries for report generation
type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// GetAnalyticsSummary aggregates analytics data
func (r *ReportRepository) GetAnalyticsSummary(ctx context.Context) (*report.AnalyticsSummary, error) {
	s := &report.AnalyticsSummary{}

	// Total Policies
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM policies`).Scan(&s.TotalPolicies); err != nil {
		return nil, err
	}

	// Total Customers
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&s.TotalCustomers); err != nil {
		return nil, err
	}

	// Total Claims
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM claims`).Scan(&s.TotalClaims); err != nil {
		return nil, err
	}

	// Total Invoices
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM invoices`).Scan(&s.TotalInvoices); err != nil {
		return nil, err
	}

	// Total Premium (sum of all policy premium_amount)
	if err := r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(premium_amount), 0) FROM policies`).Scan(&s.TotalPremium); err != nil {
		return nil, err
	}

	// Total Billing (sum of all invoice amount)
	if err := r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(amount), 0) FROM invoices`).Scan(&s.TotalBilling); err != nil {
		return nil, err
	}

	// Total Paid Claims
	if err := r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(claim_amount), 0) FROM claims WHERE status = 'paid'`).Scan(&s.TotalPaidClaims); err != nil {
		return nil, err
	}

	// Policies by Status
	rows, err := r.db.QueryContext(ctx, `SELECT status, COUNT(*) as cnt FROM policies GROUP BY status ORDER BY cnt DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var sc report.StatusCount
		if err := rows.Scan(&sc.Status, &sc.Count); err != nil {
			return nil, err
		}
		s.PoliciesByStatus = append(s.PoliciesByStatus, sc)
	}

	// Claims by Status
	rows2, err := r.db.QueryContext(ctx, `SELECT status, COUNT(*) as cnt FROM claims GROUP BY status ORDER BY cnt DESC`)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()
	for rows2.Next() {
		var sc report.StatusCount
		if err := rows2.Scan(&sc.Status, &sc.Count); err != nil {
			return nil, err
		}
		s.ClaimsByStatus = append(s.ClaimsByStatus, sc)
	}

	return s, nil
}
