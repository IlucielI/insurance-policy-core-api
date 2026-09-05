package fraud

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
)

// LLMService interface for LLM completion
type LLMService interface {
	GenerateCompletion(ctx context.Context, messages []map[string]string) (string, error)
}

// FraudDetector analyzes applications for fraud risk
type FraudDetector struct {
	llm LLMService
}

// RiskAnalysisResult contains fraud detection results
type RiskAnalysisResult struct {
	RiskScore   int      `json:"risk_score"`
	RiskLevel   string   `json:"risk_level"` // low, medium, high
	FraudFlags  []string `json:"fraud_flags"`
	Analysis    string   `json:"analysis"`
	Confidence  float64  `json:"confidence"`
	AnalyzedAt  time.Time `json:"analyzed_at"`
}

// NewFraudDetector creates a new fraud detector
func NewFraudDetector(llm LLMService) *FraudDetector {
	return &FraudDetector{
		llm: llm,
	}
}

// AnalyzeApplication performs comprehensive fraud risk analysis
func (f *FraudDetector) AnalyzeApplication(ctx context.Context, app *domain.Application, product *domain.Product) (*RiskAnalysisResult, error) {
	// Build analysis prompt
	prompt := f.buildAnalysisPrompt(app, product)
	
	messages := []map[string]string{
		{
			"role":    "system",
			"content": "You are an expert insurance fraud detection AI. Analyze applications for suspicious patterns and assign risk scores (0-100). Be thorough but fair.",
		},
		{
			"role":    "content",
			"content": prompt,
		},
	}
	
	// Call LLM for analysis
	response, err := f.llm.GenerateCompletion(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("LLM analysis failed: %w", err)
	}
	
	// Parse LLM response
	result, err := f.parseAnalysisResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse analysis: %w", err)
	}
	
	result.AnalyzedAt = time.Now()
	
	return result, nil
}

// buildAnalysisPrompt constructs detailed prompt for LLM
func (f *FraudDetector) buildAnalysisPrompt(app *domain.Application, product *domain.Product) string {
	// Extract applicant data
	name, _ := app.ApplicantData["name"].(string)
	ageVal := app.ApplicantData["age"]
	dobVal := app.ApplicantData["date_of_birth"]
	address, _ := app.ApplicantData["address"].(string)
	phone, _ := app.ApplicantData["phone"].(string)
	email, _ := app.ApplicantData["email"].(string)
	occupation, _ := app.ApplicantData["occupation"].(string)
	
	// Calculate age if not provided
	age := 0
	if ageVal != nil {
		if ageFloat, ok := ageVal.(float64); ok {
			age = int(ageFloat)
		} else if ageInt, ok := ageVal.(int); ok {
			age = ageInt
		}
	}
	
	// Parse DOB if available
	dob := ""
	if dobVal != nil {
		dob, _ = dobVal.(string)
	}
	
	// Extract health data
	healthData := ""
	if len(app.HealthQuestions) > 0 {
		healthJSON, _ := json.MarshalIndent(app.HealthQuestions, "", "  ")
		healthData = string(healthJSON)
	}
	
	// Calculate coverage ratios
	coverageRatio := float64(app.SumAssured) / float64(app.PremiumAmount)
	maxCoverageRatio := float64(product.MaxSumAssured) / float64(product.MinSumAssured)
	
	prompt := fmt.Sprintf(`Analyze this insurance application for fraud risk and suspicious patterns:

APPLICATION DETAILS:
- Application ID: %s
- Product: %s (Category: %s)
- Status: %s

APPLICANT INFORMATION:
- Name: %s
- Age: %d
- Date of Birth: %s
- Address: %s
- Phone: %s
- Email: %s
- Occupation: %s

POLICY DETAILS:
- Sum Assured: Rp %d
- Payment Term: %d months
- Premium Amount: Rp %d
- Coverage Ratio: %.2f
- Product Min Coverage: Rp %d
- Product Max Coverage: Rp %d
- Max Coverage Ratio: %.2f

HEALTH INFORMATION:
%s

FRAUD RISK FACTORS TO CHECK:
1. **Age vs Coverage Mismatch**: Is the sum assured disproportionately high for the applicant's age?
2. **Unrealistic Coverage**: Is coverage amount near maximum limit with minimal health screening?
3. **Suspicious Patterns**: 
   - Incomplete or inconsistent personal information
   - Email/phone patterns suggesting fake identity
   - Address too generic or invalid format
   - High-risk occupation with very high coverage
4. **Health Declaration**: Pre-existing conditions hidden or inconsistent answers
5. **Premium vs Coverage**: Extremely high coverage for low premium (potential error or fraud)
6. **Timing Patterns**: Rushed application without proper documentation

RESPONSE FORMAT (JSON only):
{
  "risk_score": <0-100>,
  "risk_level": "<low|medium|high>",
  "fraud_flags": [
    "flag1: description",
    "flag2: description"
  ],
  "analysis": "Detailed explanation of risk assessment",
  "confidence": <0.0-1.0>
}

SCORING GUIDELINES:
- 0-30: Low risk (normal application)
- 31-60: Medium risk (requires review)
- 61-100: High risk (likely fraud)

Provide thorough but fair analysis. Return ONLY valid JSON.`,
		app.ID,
		product.Name,
		product.Category,
		app.Status,
		name,
		age,
		dob,
		address,
		phone,
		email,
		occupation,
		app.SumAssured,
		app.PaymentTerm,
		app.PremiumAmount,
		coverageRatio,
		product.MinSumAssured,
		product.MaxSumAssured,
		maxCoverageRatio,
		healthData,
	)
	
	return prompt
}

// parseAnalysisResponse parses LLM JSON response
func (f *FraudDetector) parseAnalysisResponse(response string) (*RiskAnalysisResult, error) {
	// Clean response - extract JSON if wrapped in markdown
	response = strings.TrimSpace(response)
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)
	}
	
	var result RiskAnalysisResult
	if err := json.Unmarshal([]byte(response), &result); err != nil {
		// Fallback: try to extract values manually if JSON parse fails
		result = f.fallbackParse(response)
	}
	
	// Validate and normalize
	if result.RiskScore < 0 {
		result.RiskScore = 0
	}
	if result.RiskScore > 100 {
		result.RiskScore = 100
	}
	
	// Ensure risk level matches score
	if result.RiskLevel == "" {
		result.RiskLevel = f.getRiskLevel(result.RiskScore)
	}
	
	if result.Confidence <= 0 || result.Confidence > 1 {
		result.Confidence = 0.8 // default confidence
	}
	
	return &result, nil
}

// fallbackParse attempts manual parsing if JSON fails
func (f *FraudDetector) fallbackParse(response string) RiskAnalysisResult {
	result := RiskAnalysisResult{
		RiskScore:  50, // default medium risk
		RiskLevel:  "medium",
		FraudFlags: []string{},
		Analysis:   response,
		Confidence: 0.5,
	}
	
	// Try to extract risk score
	if strings.Contains(response, "risk_score") {
		parts := strings.Split(response, "risk_score")
		if len(parts) > 1 {
			// Find number after risk_score
			numStr := ""
			for _, char := range parts[1] {
				if char >= '0' && char <= '9' {
					numStr += string(char)
				} else if numStr != "" {
					break
				}
			}
			if score, err := strconv.Atoi(numStr); err == nil {
				result.RiskScore = score
			}
		}
	}
	
	result.RiskLevel = f.getRiskLevel(result.RiskScore)
	
	return result
}

// getRiskLevel returns risk level based on score
func (f *FraudDetector) getRiskLevel(score int) string {
	if score <= 30 {
		return "low"
	} else if score <= 60 {
		return "medium"
	}
	return "high"
}
