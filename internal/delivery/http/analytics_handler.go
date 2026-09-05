package http

import (
	"strconv"

	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	analyticsUsecase *usecase.AnalyticsUsecase
}

func NewAnalyticsHandler(analyticsUsecase *usecase.AnalyticsUsecase) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsUsecase: analyticsUsecase}
}

// GET /api/v1/admin/analytics/dashboard?months=12
func (h *AnalyticsHandler) GetDashboard(c *fiber.Ctx) error {
	monthsStr := c.Query("months", "12")
	months, err := strconv.Atoi(monthsStr)
	if err != nil || months < 1 || months > 36 {
		months = 12
	}

	data, err := h.analyticsUsecase.GetDashboardAnalytics(c.Context(), months)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch analytics data",
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

// GET /api/v1/admin/analytics/revenue?months=12
func (h *AnalyticsHandler) GetRevenue(c *fiber.Ctx) error {
	monthsStr := c.Query("months", "12")
	months, err := strconv.Atoi(monthsStr)
	if err != nil || months < 1 || months > 36 {
		months = 12
	}

	summary, err := h.analyticsUsecase.GetDashboardAnalytics(c.Context(), months)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch revenue data",
		})
	}

	return c.JSON(fiber.Map{
		"data": summary.MonthlyRevenue,
		"summary": fiber.Map{
			"total_premium_ytd": summary.Summary.TotalPremiumYTD,
			"total_approved":    summary.Summary.TotalApproved,
		},
	})
}

// GET /api/v1/admin/analytics/claims
func (h *AnalyticsHandler) GetClaimsStatus(c *fiber.Ctx) error {
	data, err := h.analyticsUsecase.GetDashboardAnalytics(c.Context(), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch claims status",
		})
	}

	return c.JSON(fiber.Map{
		"data": data.ClaimsStatus,
		"summary": fiber.Map{
			"total_claims":      data.Summary.TotalClaims,
			"total_claims_paid": data.Summary.TotalClaimsPaid,
		},
	})
}

// GET /api/v1/admin/analytics/products
func (h *AnalyticsHandler) GetTopProducts(c *fiber.Ctx) error {
	data, err := h.analyticsUsecase.GetDashboardAnalytics(c.Context(), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch product analytics",
		})
	}

	return c.JSON(fiber.Map{
		"data": data.TopProducts,
	})
}