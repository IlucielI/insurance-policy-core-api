package http

import (
	"strconv"

	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type PolicyHandler struct {
	policyUsecase *usecase.PolicyUsecase
}

func NewPolicyHandler(policyUsecase *usecase.PolicyUsecase) *PolicyHandler {
	return &PolicyHandler{policyUsecase: policyUsecase}
}

// GET /policies
func (h *PolicyHandler) ListPolicies(c *fiber.Ctx) error {
	// Get userID from JWT context (assume middleware sets it)
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Pagination
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	
	// Search & filters
	search := c.Query("search", "")
	status := c.Query("status", "")
	product := c.Query("product", "")
	dateFrom := c.Query("date_from", "")
	dateTo := c.Query("date_to", "")

	policies, total, err := h.policyUsecase.GetUserPoliciesWithFilters(
		c.Context(),
		userID.(string),
		search,
		status,
		product,
		dateFrom,
		dateTo,
		limit,
		offset,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":   policies,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GET /policies/:id
func (h *PolicyHandler) GetPolicy(c *fiber.Ctx) error {
	id := c.Params("id")

	policy, err := h.policyUsecase.GetPolicyByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(policy)
}

// POST /policies/:id/endorse
type EndorsePolicyRequest struct {
	EndorsementType string `json:"endorsement_type"`
	Description     string `json:"description"`
	EffectiveDate   string `json:"effective_date"`
}

func (h *PolicyHandler) EndorsePolicy(c *fiber.Ctx) error {
	id := c.Params("id")

	var req EndorsePolicyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.policyUsecase.EndorsePolicy(c.Context(), id, req.EndorsementType, req.Description, req.EffectiveDate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "endorsement request submitted successfully",
	})
}

// POST /policies/:id/renew
func (h *PolicyHandler) RenewPolicy(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.policyUsecase.RenewPolicy(c.Context(), id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "policy renewed successfully",
	})
}
