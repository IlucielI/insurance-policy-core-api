package http

import (
	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	reportUsecase *usecase.ReportUsecase
}

func NewReportHandler(reportUsecase *usecase.ReportUsecase) *ReportHandler {
	return &ReportHandler{reportUsecase: reportUsecase}
}

// GET /api/v1/reports/billing/:id/pdf
func (h *ReportHandler) BillingStatementPDF(c *fiber.Ctx) error {
	id := c.Params("id")
	pdf, err := h.reportUsecase.GenerateBillingStatementPDF(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", "attachment; filename=billing-statement-"+id+".pdf")
	return c.Send(pdf)
}

// GET /api/v1/reports/claims/excel
func (h *ReportHandler) ClaimsReportExcel(c *fiber.Ctx) error {
	status := c.Query("status", "")
	claimType := c.Query("claim_type", "")

	xlsx, err := h.reportUsecase.GenerateClaimsReportExcel(c.Context(), status, claimType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=laporan-klaim.xlsx")
	return c.Send(xlsx)
}

// GET /api/v1/reports/customers/excel
func (h *ReportHandler) CustomerListExcel(c *fiber.Ctx) error {
	xlsx, err := h.reportUsecase.GenerateCustomerListExcel(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=daftar-pelanggan.xlsx")
	return c.Send(xlsx)
}

// GET /api/v1/reports/analytics/excel
func (h *ReportHandler) AnalyticsSummaryExcel(c *fiber.Ctx) error {
	xlsx, err := h.reportUsecase.GenerateAnalyticsSummaryExcel(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=ringkasan-analitik.xlsx")
	return c.Send(xlsx)
}
