package http

import (
	"strconv"

	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/payment"
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

// POST /billing/pay - Create Midtrans payment
type CreatePaymentRequest struct {
	InvoiceID string `json:"invoice_id"`
}

func (h *BillingHandler) PayInvoice(c *fiber.Ctx) error {
	var req CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.InvoiceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invoice_id is required",
		})
	}

	// Create Midtrans payment transaction
	snapResp, err := h.billingUsecase.CreatePayment(c.Context(), req.InvoiceID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":      "payment transaction created",
		"snap_token":   snapResp.Token,
		"redirect_url": snapResp.RedirectURL,
	})
}

// GET /billing/verify/:order_id - Verify payment status (security: prevent URL param manipulation)
func (h *BillingHandler) VerifyPayment(c *fiber.Ctx) error {
	orderID := c.Params("order_id")
	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "order_id is required",
		})
	}

	// Verify from DATABASE, not URL params
	invoice, err := h.billingUsecase.GetInvoiceByOrderID(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "invoice not found",
		})
	}

	return c.JSON(fiber.Map{
		"order_id":          invoice.InvoiceNumber,
		"status":            invoice.Status, // 'pending', 'paid', 'failed', 'overdue'
		"amount":            invoice.Amount,
		"paid_amount":       invoice.PaidAmount,
		"paid_at":           invoice.PaidAt,
		"payment_method":    invoice.PaymentMethod,
		"payment_reference": invoice.PaymentReference,
	})
}

// POST /webhooks/payment - Handle Midtrans webhook
func (h *BillingHandler) HandlePaymentWebhook(c *fiber.Ctx) error {
	var notification payment.WebhookNotification
	if err := c.BodyParser(&notification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid webhook payload",
		})
	}

	// Handle webhook notification
	if err := h.billingUsecase.HandlePaymentWebhook(c.Context(), &notification); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
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

// GET /admin/billing/invoices - Get all invoices (admin)
func (h *BillingHandler) GetAllInvoices(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	status := c.Query("status", "")

	invoices, total, err := h.billingUsecase.GetAllInvoices(c.Context(), status, limit, offset)
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

// POST /admin/billing/invoices - Create new invoice (admin)
type CreateInvoiceRequest struct {
	PolicyID    string `json:"policy_id"`
	Amount      int64  `json:"amount"`
	DueDate     string `json:"due_date"`
	InvoiceType string `json:"invoice_type"`
	Description string `json:"description"`
}

func (h *BillingHandler) CreateInvoice(c *fiber.Ctx) error {
	var req CreateInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	invoice, err := h.billingUsecase.CreateInvoice(c.Context(), req.PolicyID, req.Amount, req.DueDate, req.InvoiceType, req.Description)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "invoice created",
		"invoice": invoice,
	})
}

// PUT /admin/billing/invoices/:id/status - Update invoice status (admin)
type UpdateInvoiceStatusRequest struct {
	Status string `json:"status"`
}

func (h *BillingHandler) UpdateInvoiceStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var req UpdateInvoiceStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.billingUsecase.UpdateInvoiceStatus(c.Context(), id, req.Status); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "invoice status updated",
	})
}
