package http

import (
	"strconv"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type ApplicationHandler struct {
	applicationUsecase *usecase.ApplicationUsecase
}

func NewApplicationHandler(applicationUsecase *usecase.ApplicationUsecase) *ApplicationHandler {
	return &ApplicationHandler{applicationUsecase: applicationUsecase}
}

// CreateApplication godoc
// @Summary Create new insurance application
// @Description Submit a new insurance application for a product
// @Tags applications
// @Accept json
// @Produce json
// @Param body body CreateApplicationRequest true "Application data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /applications [post]
type CreateApplicationRequest struct {
	ProductID       string                 `json:"product_id"`
	ApplicantData   map[string]interface{} `json:"applicant_data"`
	SumAssured      int64                  `json:"sum_assured"`
	PaymentTerm     int                    `json:"payment_term"`
	PremiumAmount   int64                  `json:"premium_amount"`
	HealthQuestions map[string]interface{} `json:"health_questions"`
}

func (h *ApplicationHandler) CreateApplication(c *fiber.Ctx) error {
	var req CreateApplicationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Get userID from context (if authenticated) or use anonymous
	userID := "anonymous"
	if uid := c.Locals("userID"); uid != nil {
		userID = uid.(string)
	}

	app := &domain.Application{
		UserID:          userID,
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
		"message": "application created successfully",
		"data":    app,
	})
}

// GetApplication godoc
// @Summary Get application by ID
// @Description Get detailed information about an application
// @Tags applications
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /applications/{id} [get]
func (h *ApplicationHandler) GetApplication(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "application ID is required",
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
			"error": "application not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": app,
	})
}

// ListApplications godoc
// @Summary List applications (Admin)
// @Description Get list of all applications with filters
// @Tags admin
// @Produce json
// @Param user_id query string false "Filter by user ID"
// @Param status query string false "Filter by status"
// @Param product_id query string false "Filter by product ID"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/applications [get]
func (h *ApplicationHandler) ListApplications(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	status := c.Query("status")
	productID := c.Query("product_id")
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	apps, total, err := h.applicationUsecase.ListApplications(c.Context(), userID, status, productID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":   apps,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// UpdateStatus godoc
// @Summary Update application status (Admin)
// @Description Approve, reject, or update application status
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Param body body UpdateStatusRequest true "Status update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/applications/{id}/status [put]
type UpdateStatusRequest struct {
	Status          string `json:"status"` // submitted, under_review, approved, rejected
	UnderwriterNotes string `json:"underwriter_notes,omitempty"`
	RejectionReason string `json:"rejection_reason,omitempty"`
}

func (h *ApplicationHandler) UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "application ID is required",
		})
	}

	var req UpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Get underwriter ID from context
	var underwriterID *string
	if uid := c.Locals("userID"); uid != nil {
		uidStr := uid.(string)
		underwriterID = &uidStr
	}

	var err error
	switch req.Status {
	case "approved":
		err = h.applicationUsecase.ApproveApplication(c.Context(), id, underwriterID, req.UnderwriterNotes)
	case "rejected":
		err = h.applicationUsecase.RejectApplication(c.Context(), id, underwriterID, req.RejectionReason, req.UnderwriterNotes)
	case "under_review":
		err = h.applicationUsecase.ReviewApplication(c.Context(), id, underwriterID)
	case "submitted":
		err = h.applicationUsecase.SubmitApplication(c.Context(), id)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid status",
		})
	}

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "application status updated successfully",
	})
}

// BulkUpdateStatus godoc
// @Summary Bulk update application status (Admin)
// @Description Approve or reject multiple applications at once
// @Tags admin
// @Accept json
// @Produce json
// @Param body body BulkUpdateStatusRequest true "Bulk status update"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /admin/applications/bulk-status [post]
type BulkUpdateStatusRequest struct {
	ApplicationIDs   []string `json:"application_ids"`
	Status           string   `json:"status"` // approved, rejected
	UnderwriterNotes string   `json:"underwriter_notes,omitempty"`
	RejectionReason  string   `json:"rejection_reason,omitempty"`
}

func (h *ApplicationHandler) BulkUpdateStatus(c *fiber.Ctx) error {
	var req BulkUpdateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if len(req.ApplicationIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "application_ids is required",
		})
	}

	// Get underwriter ID from context
	var underwriterID *string
	if uid := c.Locals("userID"); uid != nil {
		uidStr := uid.(string)
		underwriterID = &uidStr
	}

	results := make(map[string]string) // id -> "success" or error message
	successCount := 0
	errorCount := 0

	for _, appID := range req.ApplicationIDs {
		var err error
		switch req.Status {
		case "approved":
			err = h.applicationUsecase.ApproveApplication(c.Context(), appID, underwriterID, req.UnderwriterNotes)
		case "rejected":
			err = h.applicationUsecase.RejectApplication(c.Context(), appID, underwriterID, req.RejectionReason, req.UnderwriterNotes)
		default:
			err = fiber.NewError(fiber.StatusBadRequest, "invalid status for bulk update")
		}

		if err != nil {
			results[appID] = err.Error()
			errorCount++
		} else {
			results[appID] = "success"
			successCount++
		}
	}

	return c.JSON(fiber.Map{
		"message":       "bulk update completed",
		"success_count": successCount,
		"error_count":   errorCount,
		"results":       results,
	})
}
