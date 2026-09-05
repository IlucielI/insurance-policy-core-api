package http

import (
	"strconv"

	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type DocumentHandler struct {
	documentUsecase *usecase.DocumentUsecase
}

func NewDocumentHandler(documentUsecase *usecase.DocumentUsecase) *DocumentHandler {
	return &DocumentHandler{documentUsecase: documentUsecase}
}

// GET /documents
func (h *DocumentHandler) ListDocuments(c *fiber.Ctx) error {
	// Get userID from JWT context
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	documents, total, err := h.documentUsecase.GetUserDocuments(c.Context(), userID.(string), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":   documents,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// POST /documents/upload
func (h *DocumentHandler) UploadDocument(c *fiber.Ctx) error {
	// Get userID from JWT context
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Parse multipart form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file is required",
		})
	}

	// Get form fields
	documentType := c.FormValue("document_type")
	title := c.FormValue("title")
	description := c.FormValue("description")
	policyIDStr := c.FormValue("policy_id")

	if documentType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "document_type is required",
		})
	}
	if title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "title is required",
		})
	}

	var policyID *string
	if policyIDStr != "" {
		policyID = &policyIDStr
	}

	// Upload document
	doc, err := h.documentUsecase.UploadDocument(
		c.Context(),
		userID.(string),
		file,
		documentType,
		title,
		description,
		policyID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Document uploaded successfully",
		"data":    doc,
	})
}

// GET /documents/:id/download
func (h *DocumentHandler) DownloadDocument(c *fiber.Ctx) error {
	id := c.Params("id")

	document, err := h.documentUsecase.GetDocumentByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Generate presigned download URL
	downloadURL, err := h.documentUsecase.GetDocumentDownloadURL(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"id":           document.ID,
		"file_name":    document.FileName,
		"file_size":    document.FileSize,
		"mime_type":    document.MimeType,
		"download_url": downloadURL,
		"title":        document.Title,
		"uploaded_at":  document.UploadedAt,
	})
}

// DELETE /documents/:id
func (h *DocumentHandler) DeleteDocument(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.documentUsecase.DeleteDocument(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Document deleted successfully",
	})
}
