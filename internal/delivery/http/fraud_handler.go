package http

import (
	"strconv"

	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type FraudHandler struct {
	fraudUsecase *usecase.FraudUsecase
}

func NewFraudHandler(fraudUsecase *usecase.FraudUsecase) *FraudHandler {
	return &FraudHandler{fraudUsecase: fraudUsecase}
}

// POST /api/v1/admin/applications/:id/analyze-risk - Analyze fraud risk for application
func (h *FraudHandler) AnalyzeApplicationRisk(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "application id is required",
		})
	}

	result, err := h.fraudUsecase.AnalyzeApplicationRisk(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Risk analysis completed",
		"data":    result,
	})
}

// GET /api/v1/admin/fraud/high-risk - Get high-risk applications
func (h *FraudHandler) GetHighRiskApplications(c *fiber.Ctx) error {
	minScore, _ := strconv.Atoi(c.Query("min_score", "61"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	results, err := h.fraudUsecase.GetHighRiskApplications(c.Context(), minScore, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":  results,
		"total": len(results),
	})
}
