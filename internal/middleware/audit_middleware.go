package middleware

import (
	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/IlucielI/insurance-policy-core-api/internal/repository"
	"github.com/gofiber/fiber/v2"
)

// AuditContext stores audit info from request
type AuditContext struct {
	UserID    *string
	IPAddress string
	UserAgent string
}

// ExtractAuditContext extracts audit information from request
func ExtractAuditContext(c *fiber.Ctx) *AuditContext {
	ctx := &AuditContext{
		IPAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
	}

	// Extract user ID from locals (set by auth middleware)
	if userID := c.Locals("user_id"); userID != nil {
		if uid, ok := userID.(string); ok {
			ctx.UserID = &uid
		}
	}

	return ctx
}

// LogAudit creates audit log entry asynchronously
func LogAudit(c *fiber.Ctx, auditRepo *repository.AuditRepository, action, entityType, entityID string, changes map[string]interface{}) {
	auditCtx := ExtractAuditContext(c)

	log := &domain.AuditLog{
		UserID:      auditCtx.UserID,
		Action:      action,
		EntityType:  entityType,
		EntityID:    entityID,
		ChangesJSON: changes,
		IPAddress:   auditCtx.IPAddress,
		UserAgent:   auditCtx.UserAgent,
	}

	// Log asynchronously to not block request
	go func() {
		_ = auditRepo.Create(c.Context(), log)
	}()
}