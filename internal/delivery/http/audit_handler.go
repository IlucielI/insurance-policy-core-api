package http

import (
	"strconv"

	"github.com/IlucielI/insurance-policy-core-api/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type AuditHandler struct {
	auditRepo *repository.AuditRepository
}

func NewAuditHandler(auditRepo *repository.AuditRepository) *AuditHandler {
	return &AuditHandler{auditRepo: auditRepo}
}

// GET /api/v1/admin/audit-logs - List audit logs with filters
func (h *AuditHandler) ListAuditLogs(c *fiber.Ctx) error {
	// Extract query parameters
	userID := c.Query("user_id", "")
	action := c.Query("action", "")
	entityType := c.Query("entity_type", "")
	entityID := c.Query("entity_id", "")
	dateFrom := c.Query("date_from", "")
	dateTo := c.Query("date_to", "")
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	// Validate limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	// Build filters
	filters := make(map[string]interface{})
	if userID != "" {
		filters["user_id"] = userID
	}
	if action != "" {
		filters["action"] = action
	}
	if entityType != "" {
		filters["entity_type"] = entityType
	}
	if entityID != "" {
		filters["entity_id"] = entityID
	}
	if dateFrom != "" {
		filters["date_from"] = dateFrom
	}
	if dateTo != "" {
		filters["date_to"] = dateTo
	}

	// Fetch logs
	logs, total, err := h.auditRepo.List(c.Context(), filters, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":   logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GET /api/v1/admin/audit-logs/entity/:type/:id - Get audit logs for specific entity
func (h *AuditHandler) GetEntityAuditLogs(c *fiber.Ctx) error {
	entityType := c.Params("type")
	entityID := c.Params("id")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	logs, err := h.auditRepo.GetByEntityID(c.Context(), entityType, entityID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": logs,
	})
}

// GET /api/v1/admin/audit-logs/:id - Get single audit log detail
func (h *AuditHandler) GetAuditLog(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id is required",
		})
	}

	log, err := h.auditRepo.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": log,
	})
}

// GET /api/v1/admin/audit-logs/actions - Get list of available actions
func (h *AuditHandler) GetAvailableActions(c *fiber.Ctx) error {
	actions := []string{
		"approve_application",
		"reject_application",
		"send_email",
		"approve_claim",
		"reject_claim",
		"claim_status_change",
	}

	return c.JSON(fiber.Map{
		"actions": actions,
	})
}

// GET /api/v1/admin/audit-logs/entity-types - Get list of entity types
func (h *AuditHandler) GetEntityTypes(c *fiber.Ctx) error {
	entityTypes := []string{
		"application",
		"claim",
		"policy",
		"product",
		"user",
		"email",
	}

	return c.JSON(fiber.Map{
		"entity_types": entityTypes,
	})
}
