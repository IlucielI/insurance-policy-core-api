package http

import (
	"strconv"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/IlucielI/insurance-policy-core-api/internal/middleware"
	"github.com/IlucielI/insurance-policy-core-api/internal/repository"
	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type ApplicationHandler struct {
	applicationUsecase *usecase.ApplicationUsecase
	auditRepo          *repository.AuditRepository
}

func NewApplicationHandler(applicationUsecase *usecase.ApplicationUsecase, auditRepo *repository.AuditRepository) *ApplicationHandler {
	return &ApplicationHandler{applicationUsecase: applicationUsecase, auditRepo: auditRepo}
}

type CreateApplicationRequest struct {
	UserID          string                 `json:"user_id"`
	ProductID       string                 `json:"product_id"`
	ApplicantData   map[string]interface{} `json:"applicant_data"`
	SumAssured      int64                  `json:"sum_assured"`
	PaymentTerm     int                    `json:"payment_term"`
	PremiumAmount   int64                  `json:"premium_amount"`
	HealthQuestions map[string]interface{} `json:"health_questions"`
}

type UpdateStatusRequest struct {
	Status           string  `json:"status"`
	UnderwriterID    *string `json:"underwriter_id"`
	UnderwriterNotes string  `json:"underwriter_notes"`
	RejectionReason  string  `json:"rejection_reason"`
}

// POST /api/v1/applications - Create new application
func (h *ApplicationHandler) CreateApplication(c *fiber.Ctx) error {
	var req CreateApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if req.UserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}
	if req.ProductID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "product_id is required",
		})
	}
	if req.SumAssured <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "sum_assured must be greater than 0",
		})
	}
	if req.PaymentTerm <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "payment_term must be greater than 0",
		})
	}
	if req.PremiumAmount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "premium_amount must be greater than 0",
		})
	}

	app := &domain.Application{
		UserID:          req.UserID,
		ProductID:       req.ProductID,
		ApplicantData:   req.ApplicantData,
		SumAssured:      req.SumAssured,
		PaymentTerm:     req.PaymentTerm,
		PremiumAmount:   req.PremiumAmount,
		HealthQuestions: req.HealthQuestions,
	}

	if err := h.applicationUsecase.CreateApplication(c.Context(), app); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":     "Application created successfully",
		"application": app,
	})
}

// GET /api/v1/applications/:id - Get application detail
func (h *ApplicationHandler) GetApplication(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "application id is required",
		})
	}

	app, err := h.applicationUsecase.GetApplicationByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if app == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Application not found",
		})
	}

	return c.JSON(app)
}

// PUT /api/v1/applications/:id/status - Update application status (admin only)
func (h *ApplicationHandler) UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "application id is required",
		})
	}

	var req UpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "status is required",
		})
	}

	validStatuses := map[string]bool{
		"draft":        true,
		"submitted":    true,
		"under_review": true,
		"approved":     true,
		"rejected":     true,
	}
	if !validStatuses[req.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid status. Valid values: draft, submitted, under_review, approved, rejected",
		})
	}

	// For rejected status, rejection_reason is required
	if req.Status == "rejected" && req.RejectionReason == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "rejection_reason is required when status is rejected",
		})
	}

	err := h.applicationUsecase.UpdateStatus(c.Context(), id, req.Status, req.UnderwriterID, req.UnderwriterNotes, req.RejectionReason)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get updated application
	app, _ := h.applicationUsecase.GetApplicationByID(c.Context(), id)

	// Audit log
	action := req.Status
	if req.Status == "approved" {
		action = "approve_application"
	} else if req.Status == "rejected" {
		action = "reject_application"
	}
	changes := map[string]interface{}{
		"new_status":        req.Status,
		"underwriter_notes": req.UnderwriterNotes,
	}
	if req.RejectionReason != "" {
		changes["rejection_reason"] = req.RejectionReason
	}
	if h.auditRepo != nil {
		middleware.LogAudit(c, h.auditRepo, action, "application", id, changes)
	}

	return c.JSON(fiber.Map{
		"message":     "Application status updated successfully",
		"application": app,
	})
}

// BulkUpdateStatusRequest for bulk approve/reject
type BulkUpdateStatusRequest struct {
	IDs              []string `json:"ids"`
	Status           string   `json:"status"`
	UnderwriterNotes string   `json:"underwriter_notes"`
	RejectionReason  string   `json:"rejection_reason"`
}

// POST /api/v1/admin/applications/bulk-status - Bulk approve/reject applications
func (h *ApplicationHandler) BulkUpdateStatus(c *fiber.Ctx) error {
	var req BulkUpdateStatusRequest
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

	if req.Status == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "status is required",
		})
	}

	validStatuses := map[string]bool{
		"approved": true,
		"rejected": true,
	}
	if !validStatuses[req.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid status. Bulk operations support: approved, rejected",
		})
	}

	if req.Status == "rejected" && req.RejectionReason == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "rejection_reason is required when status is rejected",
		})
	}

	results := make([]fiber.Map, 0, len(req.IDs))
	successCount := 0
	failCount := 0

	for _, id := range req.IDs {
		err := h.applicationUsecase.UpdateStatus(c.Context(), id, req.Status, nil, req.UnderwriterNotes, req.RejectionReason)
		if err != nil {
			failCount++
			results = append(results, fiber.Map{
				"id":    id,
				"error": err.Error(),
			})
		} else {
			successCount++
			results = append(results, fiber.Map{
				"id":     id,
				"status": req.Status,
			})
		}
	}

	return c.JSON(fiber.Map{
		"message":       "bulk status update completed",
		"total":         len(req.IDs),
		"success_count": successCount,
		"fail_count":    failCount,
		"results":       results,
	})
}

// GET /api/v1/admin/applications - List all applications (admin only)
func (h *ApplicationHandler) ListApplications(c *fiber.Ctx) error {
	// Search & filters
	search := c.Query("search", "")
	userID := c.Query("user_id", "")
	status := c.Query("status", "")
	productID := c.Query("product_id", "")
	priority := c.Query("priority", "")
	dateFrom := c.Query("date_from", "")
	dateTo := c.Query("date_to", "")
	
	// Pagination
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	// Validate limit
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	applications, total, err := h.applicationUsecase.ListApplicationsWithFilters(
		c.Context(),
		search,
		userID,
		status,
		productID,
		priority,
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
		"data":   applications,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}
