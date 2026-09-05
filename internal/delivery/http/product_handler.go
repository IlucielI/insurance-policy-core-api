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

// ListProducts godoc
// @Summary List produk asuransi
// @Description Mendapatkan daftar produk asuransi dengan filter dan pagination
// @Tags Products
// @Produce json
// @Param category query string false "Filter berdasarkan kategori (health, life, vehicle, property)"
// @Param is_active query boolean false "Filter produk aktif/nonaktif"
// @Param limit query int false "Jumlah data per halaman" default(10)
// @Param offset query int false "Offset pagination" default(0)
// @Success 200 {object} map[string]interface{} "List produk"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /products [get]
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

// GetProduct godoc
// @Summary Get detail produk
// @Description Mendapatkan detail lengkap produk asuransi berdasarkan ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]interface{} "Detail produk"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Router /products/{id} [get]
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

// CalculatePremium godoc
// @Summary Hitung premi asuransi
// @Description Menghitung estimasi premi berdasarkan usia, sum assured, dan payment term
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body CalculatePremiumRequest true "Parameter kalkulasi premi"
// @Success 200 {object} map[string]interface{} "Hasil kalkulasi premi"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /products/{id}/calculate-premium [post]
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

type SemanticSearchRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit"`
}

// SemanticSearch godoc
// @Summary Semantic search produk
// @Description Mencari produk asuransi menggunakan natural language query
// @Tags Products
// @Accept json
// @Produce json
// @Param request body SemanticSearchRequest true "Query pencarian"
// @Success 200 {object} map[string]interface{} "Hasil pencarian"
// @Failure 400 {object} map[string]interface{} "Query is required"
// @Router /products/search [post]
func (h *ProductHandler) SemanticSearch(c *fiber.Ctx) error {
	var req SemanticSearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query is required",
		})
	}

	if req.Limit <= 0 {
		req.Limit = 5
	}

	products, err := h.productUsecase.SemanticSearchProducts(c.Context(), req.Query, req.Limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"query":   req.Query,
		"results": products,
		"count":   len(products),
	})
}
