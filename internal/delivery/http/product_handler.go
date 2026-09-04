package http

import (
	"strconv"

	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	productUsecase *usecase.ProductUsecase
}

func NewProductHandler(productUsecase *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: productUsecase}
}

func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	category := c.Query("category", "")
	isActiveStr := c.Query("is_active", "")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	var isActive *bool
	if isActiveStr != "" {
		val := isActiveStr == "true"
		isActive = &val
	}

	products, total, err := h.productUsecase.ListProducts(c.Context(), category, isActive, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":  products,
		"total": total,
		"limit": limit,
		"offset": offset,
	})
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	product, err := h.productUsecase.GetProductByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if product == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.JSON(product)
}

type CalculatePremiumRequest struct {
	Age        int   `json:"age"`
	SumAssured int64 `json:"sum_assured"`
	PaymentTerm int  `json:"payment_term"`
}

func (h *ProductHandler) CalculatePremium(c *fiber.Ctx) error {
	id := c.Params("id")

	var req CalculatePremiumRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	premium, err := h.productUsecase.CalculatePremium(c.Context(), id, req.Age, req.SumAssured, req.PaymentTerm)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"premium": premium,
		"details": fiber.Map{
			"age":          req.Age,
			"sum_assured":  req.SumAssured,
			"payment_term": req.PaymentTerm,
		},
	})
}
