package usecase

import (
	"context"

	"github.com/IlucielI/insurance-policy-core-api/internal/repository"
)

type AnalyticsUsecase struct {
	analyticsRepo *repository.AnalyticsRepository
}

func NewAnalyticsUsecase(analyticsRepo *repository.AnalyticsRepository) *AnalyticsUsecase {
	return &AnalyticsUsecase{analyticsRepo: analyticsRepo}
}

type AnalyticsDashboardResponse struct {
	Summary          *repository.SummaryStats      `json:"summary"`
	MonthlyRevenue   []repository.MonthlyRevenue   `json:"monthly_revenue"`
	ClaimsStatus     []repository.StatusCount      `json:"claims_status"`
	PolicyGrowth     []repository.MonthlyGrowth    `json:"policy_growth"`
	TopProducts      []repository.ProductRank      `json:"top_products"`
}

func (u *AnalyticsUsecase) GetDashboardAnalytics(ctx context.Context, months int) (*AnalyticsDashboardResponse, error) {
	if months <= 0 {
		months = 12
	}

	summary, err := u.analyticsRepo.SummaryStats(ctx)
	if err != nil {
		return nil, err
	}

	revenue, err := u.analyticsRepo.RevenueByMonth(ctx, months)
	if err != nil {
		return nil, err
	}

	claimsStatus, err := u.analyticsRepo.ClaimsStatusDistribution(ctx)
	if err != nil {
		return nil, err
	}

	policyGrowth, err := u.analyticsRepo.PolicyGrowthByMonth(ctx, months)
	if err != nil {
		return nil, err
	}

	topProducts, err := u.analyticsRepo.TopProducts(ctx, 10)
	if err != nil {
		return nil, err
	}

	return &AnalyticsDashboardResponse{
		Summary:        summary,
		MonthlyRevenue:       revenue,
		ClaimsStatus:   claimsStatus,
		PolicyGrowth:   policyGrowth,
		TopProducts:    topProducts,
	}, nil
}