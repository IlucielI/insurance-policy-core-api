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
}

type ProductUsecase struct {
	productRepo ProductRepository
}

func NewProductUsecase(productRepo ProductRepository) *ProductUsecase {
	return &ProductUsecase{productRepo: productRepo}
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
