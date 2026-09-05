package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/google/uuid"
)

type BillingRepository struct {
	db *sql.DB
}

func NewBillingRepository(db *sql.DB) *BillingRepository {
	return &BillingRepository{db: db}
}

func (r *BillingRepository) GetInvoicesByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Invoice, int, error) {
	// Count total
	var total int
	countQuery := `SELECT COUNT(*) FROM invoices WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch invoices
	query := `
		SELECT id, invoice_number, policy_id, user_id, amount, due_date, status,
			paid_amount, paid_at, payment_method, payment_reference, created_at, updated_at
		FROM invoices
		WHERE user_id = $1
		ORDER BY due_date DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	invoices := []*domain.Invoice{}
	for rows.Next() {
		invoice := &domain.Invoice{}
		err := rows.Scan(
			&invoice.ID, &invoice.InvoiceNumber, &invoice.PolicyID, &invoice.UserID,
			&invoice.Amount, &invoice.DueDate, &invoice.Status, &invoice.PaidAmount,
			&invoice.PaidAt, &invoice.PaymentMethod, &invoice.PaymentReference,
			&invoice.CreatedAt, &invoice.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, total, nil
}

func (r *BillingRepository) GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error) {
	invoice := &domain.Invoice{}
	query := `
		SELECT id, invoice_number, policy_id, user_id, amount, due_date, status,
			paid_amount, paid_at, payment_method, payment_reference, created_at, updated_at
		FROM invoices
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&invoice.ID, &invoice.InvoiceNumber, &invoice.PolicyID, &invoice.UserID,
		&invoice.Amount, &invoice.DueDate, &invoice.Status, &invoice.PaidAmount,
		&invoice.PaidAt, &invoice.PaymentMethod, &invoice.PaymentReference,
		&invoice.CreatedAt, &invoice.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (r *BillingRepository) PayInvoice(ctx context.Context, invoice *domain.Invoice) error {
	now := time.Now()
	invoice.PaidAt = &now
	invoice.UpdatedAt = now
	invoice.Status = "paid"

	query := `
		UPDATE invoices
		SET status = $1, paid_amount = $2, paid_at = $3, payment_method = $4,
			payment_reference = $5, updated_at = $6
		WHERE id = $7
	`

	_, err := r.db.ExecContext(ctx, query,
		invoice.Status, invoice.PaidAmount, invoice.PaidAt, invoice.PaymentMethod,
		invoice.PaymentReference, invoice.UpdatedAt, invoice.ID,
	)
	return err
}

func (r *BillingRepository) GetPaymentHistory(ctx context.Context, userID string, limit, offset int) ([]*domain.Invoice, int, error) {
	// Count total paid invoices
	var total int
	countQuery := `SELECT COUNT(*) FROM invoices WHERE user_id = $1 AND status = 'paid'`
	err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch paid invoices
	query := `
		SELECT id, invoice_number, policy_id, user_id, amount, due_date, status,
			paid_amount, paid_at, payment_method, payment_reference, created_at, updated_at
		FROM invoices
		WHERE user_id = $1 AND status = 'paid'
		ORDER BY paid_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	invoices := []*domain.Invoice{}
	for rows.Next() {
		invoice := &domain.Invoice{}
		err := rows.Scan(
			&invoice.ID, &invoice.InvoiceNumber, &invoice.PolicyID, &invoice.UserID,
			&invoice.Amount, &invoice.DueDate, &invoice.Status, &invoice.PaidAmount,
			&invoice.PaidAt, &invoice.PaymentMethod, &invoice.PaymentReference,
			&invoice.CreatedAt, &invoice.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, total, nil
}

func (r *BillingRepository) Create(ctx context.Context, invoice *domain.Invoice) error {
	invoice.ID = uuid.New().String()
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()

	query := `
		INSERT INTO invoices (id, invoice_number, policy_id, user_id, amount, due_date,
			status, paid_amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		invoice.ID, invoice.InvoiceNumber, invoice.PolicyID, invoice.UserID,
		invoice.Amount, invoice.DueDate, invoice.Status, invoice.PaidAmount,
		invoice.CreatedAt, invoice.UpdatedAt,
	)
	return err
}
