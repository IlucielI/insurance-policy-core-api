package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type OCRService struct {
	apiKey string
	client *http.Client
}

type ExtractedIDData struct {
	NIK           string `json:"nik"`
	Name          string `json:"nama"`
	DateOfBirth   string `json:"tanggal_lahir"`
	Address       string `json:"alamat"`
	Province      string `json:"provinsi,omitempty"`
	City          string `json:"kota,omitempty"`
	District      string `json:"kecamatan,omitempty"`
	Village       string `json:"kelurahan,omitempty"`
	RT            string `json:"rt,omitempty"`
	RW            string `json:"rw,omitempty"`
	Gender        string `json:"jenis_kelamin,omitempty"`
	BloodType     string `json:"golongan_darah,omitempty"`
	Religion      string `json:"agama,omitempty"`
	MaritalStatus string `json:"status_perkawinan,omitempty"`
	Occupation    string `json:"pekerjaan,omitempty"`
	Nationality   string `json:"kewarganegaraan,omitempty"`
	ValidUntil    string `json:"berlaku_hingga,omitempty"`
	RawText       string `json:"raw_text"`
	Confidence    float64 `json:"confidence"`
}

type GoogleVisionRequest struct {
	Requests []GoogleVisionImageRequest `json:"requests"`
}

type GoogleVisionImageRequest struct {
	Image    GoogleVisionImage    `json:"image"`
	Features []GoogleVisionFeature `json:"features"`
}

type GoogleVisionImage struct {
	Content string `json:"content"`
}

type GoogleVisionFeature struct {
	Type string `json:"type"`
}

type GoogleVisionResponse struct {
	Responses []GoogleVisionAnnotateResponse `json:"responses"`
}

type GoogleVisionAnnotateResponse struct {
	TextAnnotations []GoogleVisionTextAnnotation `json:"textAnnotations"`
	Error           *GoogleVisionError           `json:"error,omitempty"`
}

type GoogleVisionTextAnnotation struct {
	Description string  `json:"description"`
	Confidence  float64 `json:"confidence,omitempty"`
}

type GoogleVisionError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewOCRService() *OCRService {
	apiKey := os.Getenv("GOOGLE_VISION_API_KEY")
	
	return &OCRService{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (s *OCRService) ExtractFromImage(ctx context.Context, imagePath string) (*ExtractedIDData, error) {
	// Check if API key configured
	if s.apiKey == "" {
		// Fallback to mock data for testing
		return s.mockExtraction(imagePath), nil
	}

	// Read image file
	fileData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}

	// Encode to base64
	base64Image := base64.StdEncoding.EncodeToString(fileData)

	// Prepare request
	reqBody := GoogleVisionRequest{
		Requests: []GoogleVisionImageRequest{
			{
				Image: GoogleVisionImage{
					Content: base64Image,
				},
				Features: []GoogleVisionFeature{
					{Type: "TEXT_DETECTION"},
				},
			},
		},
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Call Google Vision API
	url := fmt.Sprintf("https://vision.googleapis.com/v1/images:annotate?key=%s", s.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(reqJSON)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Vision API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Vision API error: %s", string(body))
	}

	// Parse response
	var visionResp GoogleVisionResponse
	if err := json.Unmarshal(body, &visionResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(visionResp.Responses) == 0 {
		return nil, fmt.Errorf("no response from Vision API")
	}

	annotateResp := visionResp.Responses[0]
	if annotateResp.Error != nil {
		return nil, fmt.Errorf("Vision API error: %s", annotateResp.Error.Message)
	}

	if len(annotateResp.TextAnnotations) == 0 {
		return nil, fmt.Errorf("no text detected in image")
	}

	// First annotation contains all text
	fullText := annotateResp.TextAnnotations[0].Description
	confidence := annotateResp.TextAnnotations[0].Confidence

	// Parse extracted text
	extracted := s.parseIDCardText(fullText)
	extracted.RawText = fullText
	extracted.Confidence = confidence

	return extracted, nil
}

func (s *OCRService) mockExtraction(imagePath string) *ExtractedIDData {
	// Return mock data when no API key
	return &ExtractedIDData{
		NIK:           "3275012345678901",
		Name:          "JOHN DOE EXAMPLE",
		DateOfBirth:   "15-08-1990",
		Address:       "JL. CONTOH NO. 123 RT 001 RW 005",
		Province:      "DKI JAKARTA",
		City:          "JAKARTA SELATAN",
		District:      "KEBAYORAN BARU",
		Village:       "SENAYAN",
		RT:            "001",
		RW:            "005",
		Gender:        "LAKI-LAKI",
		BloodType:     "O",
		Religion:      "ISLAM",
		MaritalStatus: "BELUM KAWIN",
		Occupation:    "KARYAWAN SWASTA",
		Nationality:   "WNI",
		ValidUntil:    "SEUMUR HIDUP",
		RawText:       "MOCK DATA - API KEY NOT CONFIGURED",
		Confidence:    0.95,
	}
}

func (s *OCRService) parseIDCardText(text string) *ExtractedIDData {
	data := &ExtractedIDData{}
	
	// Normalize text
	text = strings.TrimSpace(text)

	// Extract NIK (16 digits)
	nikRegex := regexp.MustCompile(`(?m)(?:NIK|Nik|nik)?\s*[:.]?\s*(\d{16})`)
	if match := nikRegex.FindStringSubmatch(text); len(match) > 1 {
		data.NIK = match[1]
	}

	// Extract Name
	nameRegex := regexp.MustCompile(`(?mi)(?:Nama|NAMA|nama)\s*[:.]?\s*([A-Z\s]+(?:[A-Z\s]*))`)
	if match := nameRegex.FindStringSubmatch(text); len(match) > 1 {
		name := strings.TrimSpace(match[1])
		name = regexp.MustCompile(`[:\d]+$`).ReplaceAllString(name, "")
		data.Name = strings.TrimSpace(name)
	}

	// Extract Date of Birth
	dobRegex := regexp.MustCompile(`(?mi)(?:Tempat.*Lahir|Tgl.*Lahir|TTL|Lahir)\s*[:.]?\s*([^,\n]+),?\s*(\d{1,2}[-/]\d{1,2}[-/]\d{2,4})`)
	if match := dobRegex.FindStringSubmatch(text); len(match) > 2 {
		data.DateOfBirth = strings.TrimSpace(match[2])
	} else {
		dobSimple := regexp.MustCompile(`(\d{1,2}[-/]\d{1,2}[-/]\d{2,4})`)
		if match := dobSimple.FindStringSubmatch(text); len(match) > 1 {
			data.DateOfBirth = match[1]
		}
	}

	// Extract Gender
	genderRegex := regexp.MustCompile(`(?mi)(?:Jenis.*Kelamin|Kelamin)\s*[:.]?\s*(LAKI-LAKI|PEREMPUAN)`)
	if match := genderRegex.FindStringSubmatch(text); len(match) > 1 {
		data.Gender = strings.TrimSpace(match[1])
	}

	// Extract Address
	addressRegex := regexp.MustCompile(`(?mi)(?:Alamat|ALAMAT)\s*[:.]?\s*([^\n]+(?:\n[^\n:]+)*)`)
	if match := addressRegex.FindStringSubmatch(text); len(match) > 1 {
		addr := strings.TrimSpace(match[1])
		addr = regexp.MustCompile(`(?mi)(RT/RW|KELURAHAN|KECAMATAN|Agama|Pekerjaan|Status).*`).ReplaceAllString(addr, "")
		data.Address = strings.TrimSpace(addr)
	}

	// Extract RT/RW
	rtRwRegex := regexp.MustCompile(`(?mi)RT[/\s]*RW\s*[:.]?\s*(\d+)\s*/\s*(\d+)`)
	if match := rtRwRegex.FindStringSubmatch(text); len(match) > 2 {
		data.RT = match[1]
		data.RW = match[2]
	}

	// Extract Village/Kelurahan
	kelurahanRegex := regexp.MustCompile(`(?mi)(?:Kel/Desa|KELURAHAN|KEL\.DESA)\s*[:.]?\s*([A-Z\s]+)`)
	if match := kelurahanRegex.FindStringSubmatch(text); len(match) > 1 {
		data.Village = strings.TrimSpace(match[1])
	}

	// Extract District/Kecamatan
	kecamatanRegex := regexp.MustCompile(`(?mi)(?:Kecamatan|KECAMATAN)\s*[:.]?\s*([A-Z\s]+)`)
	if match := kecamatanRegex.FindStringSubmatch(text); len(match) > 1 {
		data.District = strings.TrimSpace(match[1])
	}

	// Extract Religion
	religionRegex := regexp.MustCompile(`(?mi)(?:Agama|AGAMA)\s*[:.]?\s*(ISLAM|KRISTEN|KATOLIK|HINDU|BUDDHA|KONGHUCU)`)
	if match := religionRegex.FindStringSubmatch(text); len(match) > 1 {
		data.Religion = strings.TrimSpace(match[1])
	}

	// Extract Marital Status
	maritalRegex := regexp.MustCompile(`(?mi)(?:Status.*Perkawinan|Perkawinan)\s*[:.]?\s*(BELUM KAWIN|KAWIN|CERAI HIDUP|CERAI MATI)`)
	if match := maritalRegex.FindStringSubmatch(text); len(match) > 1 {
		data.MaritalStatus = strings.TrimSpace(match[1])
	}

	// Extract Occupation
	occupationRegex := regexp.MustCompile(`(?mi)(?:Pekerjaan|PEKERJAAN)\s*[:.]?\s*([A-Z\s/]+)`)
	if match := occupationRegex.FindStringSubmatch(text); len(match) > 1 {
		occ := strings.TrimSpace(match[1])
		occ = regexp.MustCompile(`(?mi)(Kewarganegaraan|Berlaku).*`).ReplaceAllString(occ, "")
		data.Occupation = strings.TrimSpace(occ)
	}

	// Extract Nationality
	nationalityRegex := regexp.MustCompile(`(?mi)(?:Kewarganegaraan|KEWARGANEGARAAN)\s*[:.]?\s*(WNI|WNA)`)
	if match := nationalityRegex.FindStringSubmatch(text); len(match) > 1 {
		data.Nationality = strings.TrimSpace(match[1])
	}

	// Extract Valid Until
	validRegex := regexp.MustCompile(`(?mi)(?:Berlaku.*Hingga|BERLAKU HINGGA)\s*[:.]?\s*(SEUMUR HIDUP|\d{1,2}[-/]\d{1,2}[-/]\d{2,4})`)
	if match := validRegex.FindStringSubmatch(text); len(match) > 1 {
		data.ValidUntil = strings.TrimSpace(match[1])
	}

	// Extract Blood Type
	bloodRegex := regexp.MustCompile(`(?mi)(?:Gol.*Darah|GOLONGAN DARAH)\s*[:.]?\s*([ABO]{1,2}[+-]?)`)
	if match := bloodRegex.FindStringSubmatch(text); len(match) > 1 {
		data.BloodType = strings.TrimSpace(match[1])
	}

	return data
}
