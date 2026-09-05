package http

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/IlucielI/insurance-policy-core-api/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OCRHandler struct {
	ocrService *service.OCRService
}

func NewOCRHandler(ocrService *service.OCRService) *OCRHandler {
	return &OCRHandler{ocrService: ocrService}
}

// POST /api/v1/ocr/extract
func (h *OCRHandler) ExtractIDCard(c *fiber.Ctx) error {
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

	// Validate file type (images only)
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}

	if !allowedTypes[file.Header.Get("Content-Type")] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "only image files (JPEG, PNG, WebP) are allowed",
		})
	}

	// Check file size (max 10MB)
	if file.Size > 10*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file size must be less than 10MB",
		})
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to open file",
		})
	}
	defer src.Close()

	// Create temp file
	tempDir := os.TempDir()
	tempFileName := fmt.Sprintf("ocr_%s_%s", uuid.New().String(), filepath.Ext(file.Filename))
	tempFilePath := filepath.Join(tempDir, tempFileName)

	dst, err := os.Create(tempFilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create temp file",
		})
	}
	defer dst.Close()
	defer os.Remove(tempFilePath) // Clean up

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save file",
		})
	}

	// Extract data using OCR
	extracted, err := h.ocrService.ExtractFromImage(c.Context(), tempFilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to extract data from image",
			"details": err.Error(),
		})
	}

	// Return extracted data
	return c.JSON(fiber.Map{
		"message": "data extracted successfully",
		"data":    extracted,
	})
}
