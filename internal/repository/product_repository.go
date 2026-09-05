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

type ProductRepository struct {
	db    *sql.DB
	cache *cache.RedisClient
}

func NewProductRepository(db *sql.DB, cache *cache.RedisClient) *ProductRepository {
	return &ProductRepository{db: db, cache: cache}
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
	
	// Invalidate list cache on create
	if err == nil && r.cache != nil {
		_ = r.cache.Delete(ctx, "products:list:*")
	}
	
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
	// Build cache key from filters
	cacheKey := fmt.Sprintf("products:list:%v:%d:%d", filters, limit, offset)
	
	// Try cache first
	if r.cache != nil {
		cached, err := r.cache.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var result struct {
				Products []*domain.Product `json:"products"`
				Total    int               `json:"total"`
			}
			if json.Unmarshal([]byte(cached), &result) == nil {
				return result.Products, result.Total, nil
			}
		}
	}
	
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
	
	// Cache result for 5 minutes
	if r.cache != nil {
		result := struct {
			Products []*domain.Product `json:"products"`
			Total    int               `json:"total"`
		}{Products: products, Total: total}
		
		if cached, err := json.Marshal(result); err == nil {
			_ = r.cache.Set(ctx, cacheKey, string(cached), 5*time.Minute)
		}
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
	
	// Invalidate cache on update
	if err == nil && r.cache != nil {
		_ = r.cache.Delete(ctx, "products:list:*")
	}
	
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	
	// Invalidate cache on delete
	if err == nil && r.cache != nil {
		_ = r.cache.Delete(ctx, "products:list:*")
	}
	
	return err
}

// SaveEmbedding saves embedding for a product
func (r *ProductRepository) SaveEmbedding(ctx context.Context, productID, chunkType, chunkText string, embedding []float32) error {
	// Convert []float32 to string format for pgvector
	embeddingStr := "["
	for i, val := range embedding {
		if i > 0 {
			embeddingStr += ","
		}
		embeddingStr += fmt.Sprintf("%f", val)
	}
	embeddingStr += "]"

	query := `
		INSERT INTO product_embeddings (product_id, chunk_type, chunk_text, embedding)
		VALUES ($1, $2, $3, $4::vector)
		ON CONFLICT (product_id, chunk_type) 
		DO UPDATE SET chunk_text = $3, embedding = $4::vector, created_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, productID, chunkType, chunkText, embeddingStr)
	return err
}

// SemanticSearch performs semantic search using pgvector cosine similarity
func (r *ProductRepository) SemanticSearch(ctx context.Context, queryEmbedding []float32, limit int) ([]*domain.Product, error) {
	if limit <= 0 {
		limit = 5
	}

	// Convert []float32 to string format for pgvector
	embeddingStr := "["
	for i, val := range queryEmbedding {
		if i > 0 {
			embeddingStr += ","
		}
		embeddingStr += fmt.Sprintf("%f", val)
	}
	embeddingStr += "]"

	// Use pgvector cosine similarity (1 - cosine_distance = cosine_similarity)
	query := `
		SELECT DISTINCT p.id, p.name, p.slug, p.category, p.description, p.coverage_details, 
			p.min_sum_assured, p.max_sum_assured, p.min_payment_term, p.max_payment_term, 
			p.base_premium_rate, p.age_factor, p.is_active, p.created_at, p.updated_at,
			1 - (pe.embedding <=> $1::vector) as similarity
		FROM products p
		INNER JOIN product_embeddings pe ON p.id = pe.product_id
		WHERE p.is_active = true
		ORDER BY pe.embedding <=> $1::vector
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, embeddingStr, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*domain.Product{}
	for rows.Next() {
		product := &domain.Product{}
		var coverageDetailsJSON, ageFactorJSON []byte
		var similarity float64

		err := rows.Scan(
			&product.ID, &product.Name, &product.Slug, &product.Category, &product.Description,
			&coverageDetailsJSON, &product.MinSumAssured, &product.MaxSumAssured,
			&product.MinPaymentTerm, &product.MaxPaymentTerm, &product.BasePremiumRate,
			&ageFactorJSON, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
			&similarity,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal(coverageDetailsJSON, &product.CoverageDetails)
		json.Unmarshal(ageFactorJSON, &product.AgeFactor)

		products = append(products, product)
	}

	return products, nil
}

// GetAllProducts returns all products for embedding generation
func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]*domain.Product, error) {
	query := `
		SELECT id, name, slug, category, description, coverage_details, 
			min_sum_assured, max_sum_assured, min_payment_term, max_payment_term, 
			base_premium_rate, age_factor, is_active, created_at, updated_at
		FROM products
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		json.Unmarshal(coverageDetailsJSON, &product.CoverageDetails)
		json.Unmarshal(ageFactorJSON, &product.AgeFactor)

		products = append(products, product)
	}

	return products, nil
}
