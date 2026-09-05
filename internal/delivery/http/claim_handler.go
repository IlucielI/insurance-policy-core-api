package http

import (
	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type ClaimHandler struct {
	claimUsecase *usecase.ClaimUsecase
}

func NewClaimHandler(claimUsecase *usecase.ClaimUsecase) *ClaimHandler {
	return &ClaimHandler{claimUsecase: claimUsecase}
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
