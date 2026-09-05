package usecase

import (
	"context"
	"errors"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/fraud"
)

type FraudRepositoryInterface interface {
	UpdateRiskScore(ctx context.Context, applicationID string, riskScore int, fraudFlags []string, analysisDetail string) error
	GetHighRiskApplications(ctx context.Context, minRiskScore int, limit int) ([]map[string]interface{}, error)
}

type FraudUsecase struct {
	fraudDetector *fraud.FraudDetector
	fraudRepo     FraudRepositoryInterface
	appRepo       ApplicationRepository
	productRepo   interface {
		GetByID(ctx context.Context, id string) (*domain.Product, error)
	}
}

func NewFraudUsecase(
	fraudDetector *fraud.FraudDetector,
	fraudRepo FraudRepositoryInterface,
	appRepo ApplicationRepository,
	productRepo interface {
		GetByID(ctx context.Context, id string) (*domain.Product, error)
	},
) *FraudUsecase {
	return &FraudUsecase{
		fraudDetector: fraudDetector,
		fraudRepo:     fraudRepo,
		appRepo:       appRepo,
		productRepo:   productRepo,
	}
}

// AnalyzeApplicationRisk performs AI-powered fraud detection
func (u *FraudUsecase) AnalyzeApplicationRisk(ctx context.Context, applicationID string) (*fraud.RiskAnalysisResult, error) {
	// Get application
	app, err := u.appRepo.GetByID(ctx, applicationID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, errors.New("application not found")
	}

	// Get product details
	product, err := u.productRepo.GetByID(ctx, app.ProductID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	// Run fraud detection analysis
	result, err := u.fraudDetector.AnalyzeApplication(ctx, app, product)
	if err != nil {
		return nil, err
	}

	// Save risk score to database
	if err := u.fraudRepo.UpdateRiskScore(ctx, applicationID, result.RiskScore, result.FraudFlags, result.Analysis); err != nil {
		return nil, err
	}

	return result, nil
}

// GetHighRiskApplications returns list of high-risk applications
func (u *FraudUsecase) GetHighRiskApplications(ctx context.Context, minRiskScore int, limit int) ([]map[string]interface{}, error) {
	if minRiskScore < 0 || minRiskScore > 100 {
		return nil, errors.New("invalid risk score range (must be 0-100)")
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	return u.fraudRepo.GetHighRiskApplications(ctx, minRiskScore, limit)
}
