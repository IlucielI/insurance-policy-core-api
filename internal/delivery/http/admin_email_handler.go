package http

import (
	"github.com/IlucielI/insurance-policy-core-api/internal/middleware"
	"github.com/IlucielI/insurance-policy-core-api/internal/repository"
	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type AdminEmailHandler struct {
	applicationUsecase *usecase.ApplicationUsecase
	claimUsecase       *usecase.ClaimUsecase
	userRepo           usecase.UserRepositoryInterface
	auditRepo          *repository.AuditRepository
}

func NewAdminEmailHandler(
	applicationUsecase *usecase.ApplicationUsecase,
	claimUsecase *usecase.ClaimUsecase,
	userRepo usecase.UserRepositoryInterface,
	auditRepo *repository.AuditRepository,
) *AdminEmailHandler {
	return &AdminEmailHandler{
		applicationUsecase: applicationUsecase,
		claimUsecase:       claimUsecase,
		userRepo:           userRepo,
		auditRepo:          auditRepo,
	}
}

type SendApplicationEmailRequest struct {
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	PolicyNumber string `json:"policy_number"`
	ProductName  string `json:"product_name"`
	SumAssured   int64  `json:"sum_assured"`
}

type SendClaimEmailRequest struct {
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	ClaimNumber string `json:"claim_number"`
	Status      string `json:"status"`
	Notes       string `json:"notes"`
}

// POST /api/v1/admin/email/application-approved/:id
func (h *AdminEmailHandler) SendApplicationApprovedEmail(c *fiber.Ctx) error {
	appID := c.Params("id")
	if appID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "application id is required",
		})
	}

	// Get application with user details
	app, err := h.applicationUsecase.GetApplicationByID(c.Context(), appID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "application not found",
		})
	}

	// Get user details
	user, err := h.userRepo.GetByID(c.Context(), app.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	// Trigger email send (will call the usecase method to send)
	var req SendApplicationEmailRequest
	if err := c.BodyParser(&req); err == nil && req.Email != "" {
		// Use provided data from request
		req.Email = user.Email
		req.FullName = user.FullName
	} else {
		// Use data from DB
		req.Email = user.Email
		req.FullName = user.FullName
	}

	if req.PolicyNumber == "" {
		req.PolicyNumber = "POL-" + appID[:8]
	}

	// Get product info
	if req.ProductName == "" || req.SumAssured == 0 {
		req.SumAssured = app.SumAssured
		// ProductName would come from app details if we join
		req.ProductName = "Insurance Policy"
	}

	// Send via usecase
	err = h.applicationUsecase.SendPolicyIssuedEmailManual(
		c.Context(),
		req.Email,
		req.FullName,
		req.PolicyNumber,
		req.ProductName,
		req.SumAssured,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to send email: " + err.Error(),
		})
	}

	// Audit log
	if h.auditRepo != nil {
		middleware.LogAudit(c, h.auditRepo, "send_email", "application", appID, map[string]interface{}{
			"email_type": "application_approved",
			"to":         req.Email,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Email sent successfully",
		"to":      req.Email,
	})
}

// POST /api/v1/admin/email/claim-status/:id
func (h *AdminEmailHandler) SendClaimStatusEmail(c *fiber.Ctx) error {
	claimID := c.Params("id")
	if claimID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "claim id is required",
		})
	}

	var req SendClaimEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Get claim details
	claim, err := h.claimUsecase.GetClaimByID(c.Context(), claimID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "claim not found",
		})
	}

	// Get user details
	user, err := h.userRepo.GetByID(c.Context(), claim.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	// Use user email and claim number from DB
	req.Email = user.Email
	req.FullName = user.FullName
	req.ClaimNumber = claim.ClaimNumber

	// Status from request or claim current status
	if req.Status == "" {
		req.Status = claim.Status
	}

	// Send via usecase
	err = h.claimUsecase.SendClaimStatusEmailManual(
		c.Context(),
		req.Email,
		req.FullName,
		req.ClaimNumber,
		req.Status,
		req.Notes,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to send email: " + err.Error(),
		})
	}

	// Audit log
	if h.auditRepo != nil {
		middleware.LogAudit(c, h.auditRepo, "send_email", "claim", claimID, map[string]interface{}{
			"email_type": "claim_status",
			"status":     req.Status,
			"to":         req.Email,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Email sent successfully",
		"to":      req.Email,
	})
}

// BulkSendEmailRequest for bulk email sending
type BulkSendEmailRequest struct {
	IDs          []string `json:"ids"`
	EntityType   string   `json:"entity_type"` // "application" or "claim"
	EmailType    string   `json:"email_type"`  // "approved", "rejected", "claim_status"
	StatusNotes  string   `json:"status_notes"` // claim status notes (for claim emails)
	TemplateVars map[string]interface{} `json:"template_vars"` // optional template variables
}

// POST /api/v1/admin/email/bulk-send - Bulk send emails
func (h *AdminEmailHandler) BulkSendEmail(c *fiber.Ctx) error {
	var req BulkSendEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if len(req.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ids array is required and must not be empty",
		})
	}

	if req.EntityType != "application" && req.EntityType != "claim" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "entity_type must be 'application' or 'claim'",
		})
	}

	results := make([]fiber.Map, 0, len(req.IDs))
	successCount := 0
	failCount := 0

	if req.EntityType == "application" {
		for _, id := range req.IDs {
			app, err := h.applicationUsecase.GetApplicationByID(c.Context(), id)
			if err != nil || app == nil {
				failCount++
				results = append(results, fiber.Map{
					"id":    id,
					"error": "application not found",
				})
				continue
			}

			user, err := h.userRepo.GetByID(c.Context(), app.UserID)
			if err != nil || user == nil {
				failCount++
				results = append(results, fiber.Map{
					"id":    id,
					"error": "user not found",
				})
				continue
			}

			policyNumber := "POL-" + app.ID[:8]
			productName := "Insurance Policy"
			sumAssured := app.SumAssured

			err = h.applicationUsecase.SendPolicyIssuedEmailManual(
				c.Context(),
				user.Email,
				user.FullName,
				policyNumber,
				productName,
				sumAssured,
			)

			if err != nil {
				failCount++
				results = append(results, fiber.Map{
					"id":    id,
					"error": err.Error(),
				})
			} else {
				successCount++
				results = append(results, fiber.Map{
					"id":    id,
					"email": user.Email,
				})
				// Audit log
				if h.auditRepo != nil {
					middleware.LogAudit(c, h.auditRepo, "send_email", "application", id, map[string]interface{}{
						"email_type": "application_approved",
						"to":         user.Email,
						"bulk":       true,
					})
				}
			}
		}
	} else { // claim
		for _, id := range req.IDs {
			claim, err := h.claimUsecase.GetClaimByID(c.Context(), id)
			if err != nil || claim == nil {
				failCount++
				results = append(results, fiber.Map{
					"id":    id,
					"error": "claim not found",
				})
				continue
			}

			user, err := h.userRepo.GetByID(c.Context(), claim.UserID)
			if err != nil || user == nil {
				failCount++
				results = append(results, fiber.Map{
					"id":    id,
					"error": "user not found",
				})
				continue
			}

			status := claim.Status
			notes := req.StatusNotes
			err = h.claimUsecase.SendClaimStatusEmailManual(
				c.Context(),
				user.Email,
				user.FullName,
				claim.ClaimNumber,
				status,
				notes,
			)

			if err != nil {
				failCount++
				results = append(results, fiber.Map{
					"id":    id,
					"error": err.Error(),
				})
			} else {
				successCount++
				results = append(results, fiber.Map{
					"id":    id,
					"email": user.Email,
				})
				// Audit log for claim bulk email
				if h.auditRepo != nil {
					middleware.LogAudit(c, h.auditRepo, "send_email", "claim", id, map[string]interface{}{
						"email_type": "claim_status",
						"status":     status,
						"to":         user.Email,
						"bulk":       true,
					})
				}
			}
		}
	}

	return c.JSON(fiber.Map{
		"message":       "bulk email send completed",
		"total":         len(req.IDs),
		"success_count": successCount,
		"fail_count":    failCount,
		"results":       results,
	})
}

// POST /api/v1/admin/email/preview/:type
func (h *AdminEmailHandler) PreviewEmail(c *fiber.Ctx) error {
	templateType := c.Params("type")
	if templateType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "template type is required",
		})
	}

	// Validate template type
	validTypes := map[string]bool{
		"welcome": true,
		"policy":  true,
		"claim":   true,
		"reset":   true,
	}
	if !validTypes[templateType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid template type, valid: welcome, policy, claim, reset",
		})
	}

	// Use application usecase to get preview (it has access to email service)
	subject, html, err := h.applicationUsecase.PreviewEmail(templateType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate preview: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"subject": subject,
		"html":    html,
		"type":    templateType,
	})
}
