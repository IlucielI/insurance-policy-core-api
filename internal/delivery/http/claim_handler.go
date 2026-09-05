package http

import (
	"strconv"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/IlucielI/insurance-policy-core-api/internal/middleware"
	"github.com/IlucielI/insurance-policy-core-api/internal/repository"
	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type ClaimHandler struct {
	claimUsecase *usecase.ClaimUsecase
	auditRepo    *repository.AuditRepository
}

func NewClaimHandler(claimUsecase *usecase.ClaimUsecase, auditRepo *repository.AuditRepository) *ClaimHandler {
	return &ClaimHandler{claimUsecase: claimUsecase, auditRepo: auditRepo}
}

// POST /claims
type CreateClaimRequest struct {
	PolicyID            string `json:"policy_id"`
	ClaimType           string `json:"claim_type"`
	ClaimAmount         int64  `json:"claim_amount"`
	IncidentDate        string `json:"incident_date"`
	IncidentDescription string `json:"incident_description"`
}

func (h *ClaimHandler) CreateClaim(c *fiber.Ctx) error {
	var req CreateClaimRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	claim := &domain.Claim{
		PolicyID:            req.PolicyID,
		ClaimType:           req.ClaimType,
		ClaimAmount:         req.ClaimAmount,
		IncidentDate:        req.IncidentDate,
		IncidentDescription: req.IncidentDescription,
	}

	if err := h.claimUsecase.CreateClaim(c.Context(), claim); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      "claim created successfully",
		"claim_id":     claim.ID,
		"claim_number": claim.ClaimNumber,
	})
}

// GET /claims/:id
func (h *ClaimHandler) GetClaim(c *fiber.Ctx) error {
	id := c.Params("id")

	claim, err := h.claimUsecase.GetClaimByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(claim)
}

// PUT /claims/:id/documents
type UploadClaimDocumentRequest struct {
	DocumentType string `json:"document_type"`
	FileName     string `json:"file_name"`
	FilePath     string `json:"file_path"`
	FileSize     int64  `json:"file_size"`
	MimeType     string `json:"mime_type"`
}

func (h *ClaimHandler) UploadDocument(c *fiber.Ctx) error {
	id := c.Params("id")

	var req UploadClaimDocumentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.claimUsecase.UploadDocument(
		c.Context(),
		id,
		req.DocumentType,
		req.FileName,
		req.FilePath,
		req.MimeType,
		req.FileSize,
	); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "document uploaded successfully",
	})
}

// GET /claims/:id/timeline
func (h *ClaimHandler) GetClaimTimeline(c *fiber.Ctx) error {
	id := c.Params("id")

	timeline, err := h.claimUsecase.GetClaimTimeline(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": timeline,
	})
}

// GET /admin/claims - List all claims with search and filters (admin only)
func (h *ClaimHandler) ListClaims(c *fiber.Ctx) error {
	// Search & filters
	search := c.Query("search", "")
	status := c.Query("status", "")
	claimType := c.Query("claim_type", "")
	priority := c.Query("priority", "")
	dateFrom := c.Query("date_from", "")
	dateTo := c.Query("date_to", "")
	amountMin := c.Query("amount_min", "")
	amountMax := c.Query("amount_max", "")
	
	// Pagination
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	claims, total, err := h.claimUsecase.ListClaimsWithFilters(
		c.Context(),
		search,
		status,
		claimType,
		priority,
		dateFrom,
		dateTo,
		amountMin,
		amountMax,
		limit,
		offset,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":   claims,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// PUT /admin/claims/:id/status - Update claim status
type UpdateClaimStatusRequest struct {
	Status string `json:"status"`
	Notes  string `json:"notes"`
}

func (h *ClaimHandler) UpdateClaimStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var req UpdateClaimStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.claimUsecase.UpdateClaimStatus(c.Context(), id, req.Status, req.Notes); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Audit log
	if h.auditRepo != nil {
		action := "claim_status_change"
		if req.Status == "approved" {
			action = "approve_claim"
		} else if req.Status == "rejected" {
			action = "reject_claim"
		}
		middleware.LogAudit(c, h.auditRepo, action, "claim", id, map[string]interface{}{
			"new_status": req.Status,
			"notes":      req.Notes,
		})
	}

	return c.JSON(fiber.Map{
		"message": "claim status updated",
	})
}

// PUT /admin/claims/:id/approve - Approve claim
type ApproveClaimRequest struct {
	ApprovedAmount int64  `json:"approved_amount"`
	Notes          string `json:"notes"`
}

func (h *ClaimHandler) ApproveClaim(c *fiber.Ctx) error {
	id := c.Params("id")
	var req ApproveClaimRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.claimUsecase.UpdateClaimStatus(c.Context(), id, "approved", req.Notes); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Audit log
	if h.auditRepo != nil {
		middleware.LogAudit(c, h.auditRepo, "approve_claim", "claim", id, map[string]interface{}{
			"approved_amount": req.ApprovedAmount,
			"notes":           req.Notes,
		})
	}

	return c.JSON(fiber.Map{
		"message": "claim approved",
	})
}

// BulkUpdateClaimStatusRequest for bulk claim status update
type BulkUpdateClaimStatusRequest struct {
	IDs    []string `json:"ids"`
	Status string   `json:"status"`
	Notes  string   `json:"notes"`
}

// POST /api/v1/admin/claims/bulk-status - Bulk update claim statuses
func (h *ClaimHandler) BulkUpdateClaimStatus(c *fiber.Ctx) error {
	var req BulkUpdateClaimStatusRequest
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
		"new":                 true,
		"under_investigation": true,
		"assigned":            true,
		"pending_approval":    true,
		"approved":            true,
		"rejected":            true,
		"paid":                true,
		"partially_paid":      true,
	}
	if !validStatuses[req.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid status. Valid: new, under_investigation, assigned, pending_approval, approved, rejected, paid, partially_paid",
		})
	}

	results := make([]fiber.Map, 0, len(req.IDs))
	successCount := 0
	failCount := 0

	for _, id := range req.IDs {
		err := h.claimUsecase.UpdateClaimStatus(c.Context(), id, req.Status, req.Notes)
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
			// Audit log each successful update
			if h.auditRepo != nil {
				action := "claim_status_change"
				if req.Status == "approved" {
					action = "approve_claim"
				} else if req.Status == "rejected" {
					action = "reject_claim"
				}
				middleware.LogAudit(c, h.auditRepo, action, "claim", id, map[string]interface{}{
					"new_status": req.Status,
					"notes":      req.Notes,
					"bulk":       true,
				})
			}
		}
	}

	return c.JSON(fiber.Map{
		"message":       "bulk claim status update completed",
		"total":         len(req.IDs),
		"success_count": successCount,
		"fail_count":    failCount,
		"results":       results,
	})
}
