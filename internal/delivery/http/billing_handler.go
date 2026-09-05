package http

import (
	"strconv"

	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type BillingHandler struct {
	billingUsecase *usecase.BillingUsecase
}

func NewBillingHandler(billingUsecase *usecase.BillingUsecase) *BillingHandler {
	return &BillingHandler{billingUsecase: billingUsecase}
}

// GET /billing/invoices
func (h *BillingHandler) GetInvoices(c *fiber.Ctx) error {
	// Get userID from JWT context
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	invoices, total, err := h.billingUsecase.GetInvoices(c.Context(), userID.(string), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":   invoices,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// POST /billing/pay
type PayInvoiceRequest struct {
	InvoiceID        string `json:"invoice_id"`
	PaymentMethod    string `json:"payment_method"`
	PaymentReference string `json:"payment_reference"`
}

func (h *BillingHandler) PayInvoice(c *fiber.Ctx) error {
	var req PayInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.billingUsecase.PayInvoice(c.Context(), req.InvoiceID, req.PaymentMethod, req.PaymentReference); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "payment processed successfully",
	})
}

// GET /billing/history
func (h *BillingHandler) GetPaymentHistory(c *fiber.Ctx) error {
	// Get userID from JWT context
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	invoices, total, err := h.billingUsecase.GetPaymentHistory(c.Context(), userID.(string), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":   invoices,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}
