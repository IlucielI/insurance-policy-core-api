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

// GET /documents/:id/download
func (h *DocumentHandler) DownloadDocument(c *fiber.Ctx) error {
	id := c.Params("id")

	document, err := h.documentUsecase.GetDocumentByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return document metadata and download URL
	// In production, you'd generate a signed URL or serve the file
	return c.JSON(fiber.Map{
		"id":            document.ID,
		"file_name":     document.FileName,
		"file_path":     document.FilePath,
		"file_size":     document.FileSize,
		"mime_type":     document.MimeType,
		"download_url":  "/files/" + document.FilePath, // Placeholder
		"document_type": document.DocumentType,
		"title":         document.Title,
	})
}
