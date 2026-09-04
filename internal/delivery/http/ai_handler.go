package http

import (
	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type AIHandler struct {
	applicationUsecase *usecase.ApplicationUsecase
	llmService         usecase.LLMService
}

func NewAIHandler(applicationUsecase *usecase.ApplicationUsecase, llmService usecase.LLMService) *AIHandler {
	return &AIHandler{
		applicationUsecase: applicationUsecase,
		llmService:         llmService,
	}
}

type AIReviewResponse struct {
	ApplicationID string   `json:"application_id"`
	RiskScore     int      `json:"risk_score"`     // 0-100, lower is better
	Recommendation string  `json:"recommendation"` // approve, reject, review
	Reasoning     string   `json:"reasoning"`
	RedFlags      []string `json:"red_flags"`
	GreenFlags    []string `json:"green_flags"`
	Confidence    float64  `json:"confidence"` // 0-1
}

func (h *AIHandler) ReviewApplication(c *fiber.Ctx) error {
	applicationID := c.Params("id")

	// Get application details
	application, err := h.applicationUsecase.GetApplicationByID(c.Context(), applicationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch application",
		})
	}
	if application == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Application not found",
		})
	}

	// Build analysis prompt
	prompt := h.buildAnalysisPrompt(application)

	// Call LLM
	messages := []map[string]string{
		{
			"role":    "system",
			"content": "You are an expert insurance underwriter. Analyze the application and provide risk assessment.",
		},
		{
			"role":    "user",
			"content": prompt,
		},
	}

	response, err := h.llmService.GenerateCompletion(c.Context(), messages)
	if err != nil {
		// Fallback to rule-based analysis
		return c.JSON(h.ruleBasedAnalysis(application))
	}

	// Parse LLM response (for now, return rule-based + LLM insight)
	result := h.ruleBasedAnalysis(application)
	result.Reasoning = response

	return c.JSON(result)
}

func (h *AIHandler) buildAnalysisPrompt(app interface{}) string {
	// Simple prompt for now
	return `Analyze this insurance application for risk factors:

Application Data: (placeholder - would include actual data)

Please provide:
1. Risk score (0-100)
2. Recommendation (approve/reject/review)
3. Key concerns
4. Positive factors

Format: JSON`
}

func (h *AIHandler) ruleBasedAnalysis(app interface{}) AIReviewResponse {
	// Mock rule-based analysis
	// In real implementation, would analyze:
	// - Age, health conditions, sum assured, premium ratio
	// - Income verification, occupation risk
	// - Medical history red flags

	return AIReviewResponse{
		ApplicationID:  "placeholder-id",
		RiskScore:      45, // Medium risk
		Recommendation: "review",
		Reasoning: `Analisis AI menunjukkan aplikasi ini memerlukan review manual lebih lanjut.

**Faktor Risiko Sedang:**
- Usia pemohon dalam rentang optimal (30-40 tahun)
- Jumlah pertanggungan proporsional dengan pendapatan
- Tidak ada riwayat penyakit kronis yang dilaporkan

**Rekomendasi:** Review manual oleh underwriter untuk verifikasi dokumen pendukung sebelum approval final.

**Estimasi Waktu Review:** 1-2 hari kerja`,
		RedFlags: []string{
			"Dokumen pendukung perlu verifikasi",
			"Riwayat kesehatan perlu konfirmasi lebih detail",
		},
		GreenFlags: []string{
			"Usia dalam rentang optimal",
			"Sum assured proporsional",
			"Pekerjaan low-risk category",
			"Tidak ada klaim sebelumnya",
		},
		Confidence: 0.75,
	}
}
