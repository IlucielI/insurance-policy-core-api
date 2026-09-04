package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/google/uuid"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	product.ID = uuid.New().String()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	coverageDetailsJSON, _ := json.Marshal(product.CoverageDetails)
	ageFactorJSON, _ := json.Marshal(product.AgeFactor)

	query := `
		INSERT INTO products (id, name, slug, category, description, coverage_details, 
			min_sum_assured, max_sum_assured, min_payment_term, max_payment_term, 
			base_premium_rate, age_factor, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	_, err := r.db.ExecContext(ctx, query,
		product.ID, product.Name, product.Slug, product.Category, product.Description,
		coverageDetailsJSON, product.MinSumAssured, product.MaxSumAssured,
		product.MinPaymentTerm, product.MaxPaymentTerm, product.BasePremiumRate,
		ageFactorJSON, product.IsActive, product.CreatedAt, product.UpdatedAt,
	)
	return err
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	product := &domain.Product{}
	var coverageDetailsJSON, ageFactorJSON []byte

	query := `
		SELECT id, name, slug, category, description, coverage_details, 
			min_sum_assured, max_sum_assured, min_payment_term, max_payment_term, 
			base_premium_rate, age_factor, is_active, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.Name, &product.Slug, &product.Category, &product.Description,
		&coverageDetailsJSON, &product.MinSumAssured, &product.MaxSumAssured,
		&product.MinPaymentTerm, &product.MaxPaymentTerm, &product.BasePremiumRate,
		&ageFactorJSON, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(coverageDetailsJSON, &product.CoverageDetails)
	json.Unmarshal(ageFactorJSON, &product.AgeFactor)

	return product, nil
}

func (r *ProductRepository) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	product := &domain.Product{}
	var coverageDetailsJSON, ageFactorJSON []byte

	query := `
		SELECT id, name, slug, category, description, coverage_details, 
			min_sum_assured, max_sum_assured, min_payment_term, max_payment_term, 
			base_premium_rate, age_factor, is_active, created_at, updated_at
		FROM products
		WHERE slug = $1
	`
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&product.ID, &product.Name, &product.Slug, &product.Category, &product.Description,
		&coverageDetailsJSON, &product.MinSumAssured, &product.MaxSumAssured,
		&product.MinPaymentTerm, &product.MaxPaymentTerm, &product.BasePremiumRate,
		&ageFactorJSON, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(coverageDetailsJSON, &product.CoverageDetails)
	json.Unmarshal(ageFactorJSON, &product.AgeFactor)

	return product, nil
}

func (r *ProductRepository) List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*domain.Product, int, error) {
	// Build query with filters
	whereClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if category, ok := filters["category"].(string); ok && category != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("category = $%d", argPos))
		args = append(args, category)
		argPos++
	}

	if isActive, ok := filters["is_active"].(bool); ok {
		whereClauses = append(whereClauses, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, isActive)
		argPos++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch products
	query := fmt.Sprintf(`
		SELECT id, name, slug, category, description, coverage_details, 
			min_sum_assured, max_sum_assured, min_payment_term, max_payment_term, 
			base_premium_rate, age_factor, is_active, created_at, updated_at
		FROM products
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := []*domain.Product{}
	for rows.Next() {
		product := &domain.Product{}
		var coverageDetailsJSON, ageFactorJSON []byte

		err := rows.Scan(
			&product.ID, &product.Name, &product.Slug, &product.Category, &product.Description,
			&coverageDetailsJSON, &product.MinSumAssured, &product.MaxSumAssured,
			&product.MinPaymentTerm, &product.MaxPaymentTerm, &product.BasePremiumRate,
			&ageFactorJSON, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		json.Unmarshal(coverageDetailsJSON, &product.CoverageDetails)
		json.Unmarshal(ageFactorJSON, &product.AgeFactor)

		products = append(products, product)
	}

	return products, total, nil
}

func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	product.UpdatedAt = time.Now()

	coverageDetailsJSON, _ := json.Marshal(product.CoverageDetails)
	ageFactorJSON, _ := json.Marshal(product.AgeFactor)

	query := `
		UPDATE products
		SET name = $1, slug = $2, category = $3, description = $4, coverage_details = $5,
			min_sum_assured = $6, max_sum_assured = $7, min_payment_term = $8, 
			max_payment_term = $9, base_premium_rate = $10, age_factor = $11, 
			is_active = $12, updated_at = $13
		WHERE id = $14
	`
	_, err := r.db.ExecContext(ctx, query,
		product.Name, product.Slug, product.Category, product.Description,
		coverageDetailsJSON, product.MinSumAssured, product.MaxSumAssured,
		product.MinPaymentTerm, product.MaxPaymentTerm, product.BasePremiumRate,
		ageFactorJSON, product.IsActive, product.UpdatedAt, product.ID,
	)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
