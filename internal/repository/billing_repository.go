package repository

import (
	"context"
	"database/sql"
	"fmt"
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

// Common invoice SELECT fields including new columns
const invoiceSelectFields = `id, invoice_number, policy_id, user_id, invoice_type, amount, due_date, description, status,
	paid_amount, paid_at, payment_method, payment_reference, created_at, updated_at`

func scanInvoice(row interface{ Scan(...interface{}) error }) (*domain.Invoice, error) {
	i := &domain.Invoice{}
	err := row.Scan(
		&i.ID, &i.InvoiceNumber, &i.PolicyID, &i.UserID,
		&i.InvoiceType, &i.Amount, &i.DueDate, &i.Description, &i.Status, &i.PaidAmount,
		&i.PaidAt, &i.PaymentMethod, &i.PaymentReference,
		&i.CreatedAt, &i.UpdatedAt,
	)
	return i, err
}

func (r *BillingRepository) GetInvoicesByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Invoice, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM invoices WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`SELECT %s FROM invoices WHERE user_id = $1 ORDER BY due_date DESC LIMIT $2 OFFSET $3`, invoiceSelectFields)
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var invoices []*domain.Invoice
	for rows.Next() {
		i, err := scanInvoice(rows)
		if err != nil {
			return nil, 0, err
		}
		invoices = append(invoices, i)
	}
	return invoices, total, rows.Err()
}

func (r *BillingRepository) GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error) {
	query := fmt.Sprintf(`SELECT %s FROM invoices WHERE id = $1`, invoiceSelectFields)
	row := r.db.QueryRowContext(ctx, query, id)
	i, err := scanInvoice(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return i, err
}

// GetInvoiceByOrderID gets invoice by invoice_number (used as order_id in payment)
func (r *BillingRepository) GetInvoiceByOrderID(ctx context.Context, orderID string) (*domain.Invoice, error) {
	query := fmt.Sprintf(`SELECT %s FROM invoices WHERE invoice_number = $1`, invoiceSelectFields)
	row := r.db.QueryRowContext(ctx, query, orderID)
	i, err := scanInvoice(row)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return i, err
}

func (r *BillingRepository) PayInvoice(ctx context.Context, invoice *domain.Invoice) error {
	now := time.Now()
	invoice.PaidAt = &now
	invoice.UpdatedAt = now
	invoice.Status = "paid"

	query := `UPDATE invoices SET status = $1, paid_amount = $2, paid_at = $3, payment_method = $4,
		payment_reference = $5, updated_at = $6 WHERE id = $7`
	_, err := r.db.ExecContext(ctx, query,
		invoice.Status, invoice.PaidAmount, invoice.PaidAt, invoice.PaymentMethod,
		invoice.PaymentReference, invoice.UpdatedAt, invoice.ID,
	)
	return err
}

func (r *BillingRepository) GetPaymentHistory(ctx context.Context, userID string, limit, offset int) ([]*domain.Invoice, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM invoices WHERE user_id = $1 AND status = 'paid'`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`SELECT %s FROM invoices WHERE user_id = $1 AND status = 'paid' ORDER BY paid_at DESC LIMIT $2 OFFSET $3`, invoiceSelectFields)
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var invoices []*domain.Invoice
	for rows.Next() {
		i, err := scanInvoice(rows)
		if err != nil {
			return nil, 0, err
		}
		invoices = append(invoices, i)
	}
	return invoices, total, rows.Err()
}

func (r *BillingRepository) Create(ctx context.Context, invoice *domain.Invoice) error {
	invoice.ID = uuid.New().String()
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()

	query := `INSERT INTO invoices (id, invoice_number, policy_id, user_id, invoice_type, amount, due_date,
		description, status, paid_amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.ExecContext(ctx, query,
		invoice.ID, invoice.InvoiceNumber, invoice.PolicyID, invoice.UserID,
		invoice.InvoiceType, invoice.Amount, invoice.DueDate,
		invoice.Description, invoice.Status, invoice.PaidAmount,
		invoice.CreatedAt, invoice.UpdatedAt,
	)
	return err
}

func (r *BillingRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	query := `INSERT INTO payments (id, application_id, order_id, gross_amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query,
		payment.ID, payment.ApplicationID, payment.OrderID, payment.GrossAmount,
		payment.Status, payment.CreatedAt, payment.UpdatedAt,
	)
	return err
}

func (r *BillingRepository) GetPaymentByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	payment := &domain.Payment{}
	query := `SELECT id, application_id, order_id, midtrans_transaction_id, payment_type,
		gross_amount, status, paid_at, expired_at, created_at, updated_at
		FROM payments WHERE order_id = $1`

	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&payment.ID, &payment.ApplicationID, &payment.OrderID,
		&payment.MidtransTransactionID, &payment.PaymentType,
		&payment.GrossAmount, &payment.Status, &payment.PaidAt,
		&payment.ExpiredAt, &payment.CreatedAt, &payment.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return payment, err
}

func (r *BillingRepository) UpdatePaymentStatus(ctx context.Context, payment *domain.Payment) error {
	query := `UPDATE payments SET midtrans_transaction_id = $1, payment_type = $2, status = $3,
		paid_at = $4, updated_at = $5 WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query,
		payment.MidtransTransactionID, payment.PaymentType, payment.Status,
		payment.PaidAt, payment.UpdatedAt, payment.ID,
	)
	return err
}

// GetAllInvoices returns all invoices with optional status filter (admin)
func (r *BillingRepository) GetAllInvoices(ctx context.Context, status string, limit, offset int) ([]*domain.Invoice, int, error) {
	var total int
	if status != "" {
		err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM invoices WHERE status = $1`, status).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM invoices`).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
	}

	var query string
	var rows *sql.Rows
	var err error
	if status != "" {
		query = fmt.Sprintf(`SELECT %s FROM invoices WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, invoiceSelectFields)
		rows, err = r.db.QueryContext(ctx, query, status, limit, offset)
	} else {
		query = fmt.Sprintf(`SELECT %s FROM invoices ORDER BY created_at DESC LIMIT $1 OFFSET $2`, invoiceSelectFields)
		rows, err = r.db.QueryContext(ctx, query, limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var invoices []*domain.Invoice
	for rows.Next() {
		i, err := scanInvoice(rows)
		if err != nil {
			return nil, 0, err
		}
		invoices = append(invoices, i)
	}
	return invoices, total, rows.Err()
}

// CreateInvoice inserts a new invoice (admin)
func (r *BillingRepository) CreateInvoice(ctx context.Context, invoice *domain.Invoice) error {
	invoice.ID = uuid.New().String()
	invoice.InvoiceNumber = fmt.Sprintf("INV-%d", time.Now().Unix())
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()

	query := `INSERT INTO invoices (id, invoice_number, policy_id, user_id, invoice_type, amount, due_date,
		description, status, paid_amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := r.db.ExecContext(ctx, query,
		invoice.ID, invoice.InvoiceNumber, invoice.PolicyID, invoice.UserID,
		invoice.InvoiceType, invoice.Amount, invoice.DueDate,
		invoice.Description, invoice.Status, invoice.PaidAmount,
		invoice.CreatedAt, invoice.UpdatedAt,
	)
	return err
}

// UpdateInvoiceStatus updates invoice status (admin)
func (r *BillingRepository) UpdateInvoiceStatus(ctx context.Context, invoiceID, status string) error {
	query := `UPDATE invoices SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), invoiceID)
	return err
}
