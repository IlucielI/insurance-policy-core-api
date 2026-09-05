package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id string) (*domain.Product, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Product, error)
	List(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*domain.Product, int, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id string) error
	SemanticSearch(ctx context.Context, queryEmbedding []float32, limit int) ([]*domain.Product, error)
	GetAllProducts(ctx context.Context) ([]*domain.Product, error)
	SaveEmbedding(ctx context.Context, productID, chunkType, chunkText string, embedding []float32) error
}

type EmbeddingsService interface {
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
}

type ProductUsecase struct {
	productRepo       ProductRepository
	embeddingsService EmbeddingsService
}

func NewProductUsecase(productRepo ProductRepository, embeddingsService EmbeddingsService) *ProductUsecase {
	return &ProductUsecase{
		productRepo:       productRepo,
		embeddingsService: embeddingsService,
	}
}

func (u *ProductUsecase) CreateProduct(ctx context.Context, product *domain.Product) error {
	// Validate
	if product.Name == "" || product.Slug == "" {
		return errors.New("name and slug are required")
	}
	if product.MinSumAssured <= 0 || product.MaxSumAssured <= 0 {
		return errors.New("invalid sum assured range")
	}
	if product.MinSumAssured > product.MaxSumAssured {
		return errors.New("min sum assured cannot be greater than max")
	}

	// Check slug unique
	existing, err := u.productRepo.GetBySlug(ctx, product.Slug)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("slug already exists")
	}

	product.IsActive = true
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	return u.productRepo.Create(ctx, product)
}

func (u *ProductUsecase) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	return u.productRepo.GetByID(ctx, id)
}

func (u *ProductUsecase) GetProductBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	return u.productRepo.GetBySlug(ctx, slug)
}

func (u *ProductUsecase) ListProducts(ctx context.Context, category string, isActive *bool, limit, offset int) ([]*domain.Product, int, error) {
	filters := make(map[string]interface{})
	if category != "" {
		filters["category"] = category
	}
	if isActive != nil {
		filters["is_active"] = *isActive
	}

	return u.productRepo.List(ctx, filters, limit, offset)
}

func (u *ProductUsecase) UpdateProduct(ctx context.Context, product *domain.Product) error {
	// Check exists
	existing, err := u.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("product not found")
	}

	// Check slug conflict (if changed)
	if product.Slug != existing.Slug {
		slugConflict, err := u.productRepo.GetBySlug(ctx, product.Slug)
		if err != nil {
			return err
		}
		if slugConflict != nil && slugConflict.ID != product.ID {
			return errors.New("slug already exists")
		}
	}

	product.UpdatedAt = time.Now()
	return u.productRepo.Update(ctx, product)
}

func (u *ProductUsecase) DeleteProduct(ctx context.Context, id string) error {
	return u.productRepo.Delete(ctx, id)
}

func (u *ProductUsecase) CalculatePremium(ctx context.Context, productID string, age int, sumAssured int64, paymentTerm int) (int64, error) {
	product, err := u.productRepo.GetByID(ctx, productID)
	if err != nil {
		return 0, err
	}
	if product == nil {
		return 0, errors.New("product not found")
	}

	// Base premium calculation
	basePremium := float64(sumAssured) * product.BasePremiumRate / 100

	// Apply age factor
	ageFactor := 1.0
	if product.AgeFactor != nil {
		// Simple age bracket logic (can be improved)
		if age < 30 {
			if factor, ok := product.AgeFactor["under_30"]; ok {
				ageFactor = factor
			}
		} else if age < 40 {
			if factor, ok := product.AgeFactor["30_39"]; ok {
				ageFactor = factor
			}
		} else if age < 50 {
			if factor, ok := product.AgeFactor["40_49"]; ok {
				ageFactor = factor
			}
		} else {
			if factor, ok := product.AgeFactor["50_plus"]; ok {
				ageFactor = factor
			}
		}
	}

	// Apply term factor (longer term = lower monthly premium)
	termFactor := float64(paymentTerm) / 12.0 // normalize to years

	finalPremium := (basePremium * ageFactor) / termFactor

	// Round to nearest 1000
	return int64(finalPremium/1000) * 1000, nil
}

// SemanticSearchProducts performs semantic search on products
func (u *ProductUsecase) SemanticSearchProducts(ctx context.Context, query string, limit int) ([]*domain.Product, error) {
	if query == "" {
		return nil, errors.New("search query is required")
	}

	if u.embeddingsService == nil {
		return nil, errors.New("embeddings service not available")
	}

	// Generate embedding for search query
	queryEmbedding, err := u.embeddingsService.GenerateEmbedding(ctx, query)
	if err != nil {
		return nil, errors.New("failed to generate query embedding: " + err.Error())
	}

	// Perform semantic search
	products, err := u.productRepo.SemanticSearch(ctx, queryEmbedding, limit)
	if err != nil {
		return nil, errors.New("semantic search failed: " + err.Error())
	}

	return products, nil
}

// GenerateProductEmbeddings generates embeddings for all products
func (u *ProductUsecase) GenerateProductEmbeddings(ctx context.Context) error {
	if u.embeddingsService == nil {
		return errors.New("embeddings service not available")
	}

	products, err := u.productRepo.GetAllProducts(ctx)
	if err != nil {
		return err
	}

	for _, product := range products {
		// Create searchable text from product data
		searchText := product.Name + ". " + product.Description
		if product.Category != "" {
			searchText = "Kategori: " + product.Category + ". " + searchText
		}

		// Generate embedding
		embedding, err := u.embeddingsService.GenerateEmbedding(ctx, searchText)
		if err != nil {
			return errors.New("failed to generate embedding for product " + product.Name + ": " + err.Error())
		}

		// Save embedding
		err = u.productRepo.SaveEmbedding(ctx, product.ID, "description", searchText, embedding)
		if err != nil {
			return errors.New("failed to save embedding for product " + product.Name + ": " + err.Error())
		}
	}

	return nil
}
