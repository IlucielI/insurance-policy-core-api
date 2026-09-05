package http

// Analytics, Reports, Fraud, OCR Handler Swagger Documentation

// GetDashboard godoc
// @Summary Get dashboard analytics (Admin)
// @Description Mendapatkan ringkasan analytics untuk dashboard admin
// @Tags Admin - Analytics
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.AnalyticsDashboard "Dashboard analytics data"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires admin roles"
// @Router /admin/analytics/dashboard [get]

// GetRevenue godoc
// @Summary Get revenue analytics (Admin)
// @Description Mendapatkan data revenue dengan breakdown per periode
// @Tags Admin - Analytics
// @Produce json
// @Security BearerAuth
// @Param from query string false "Tanggal mulai (YYYY-MM-DD)"
// @Param to query string false "Tanggal akhir (YYYY-MM-DD)"
// @Success 200 {object} object{revenue=[]object} "Revenue data"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /admin/analytics/revenue [get]

// GetClaimsStatus godoc
// @Summary Get claims status analytics (Admin)
// @Description Mendapatkan distribusi status klaim
// @Tags Admin - Analytics
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{claims=[]object} "Claims status distribution"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /admin/analytics/claims [get]

// GetTopProducts godoc
// @Summary Get top products (Admin)
// @Description Mendapatkan produk dengan penjualan terbanyak
// @Tags Admin - Analytics
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Jumlah produk" default(10)
// @Success 200 {object} object{products=[]object} "Top products"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /admin/analytics/products [get]

// BillingStatementPDF godoc
// @Summary Export billing statement PDF (Admin)
// @Description Generate dan download PDF billing statement
// @Tags Admin - Reports
// @Produce application/pdf
// @Security BearerAuth
// @Param id path string true "Billing ID"
// @Success 200 {file} file "PDF file"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or finance role"
// @Router /reports/billing/{id}/pdf [get]

// ClaimsReportExcel godoc
// @Summary Export claims report Excel (Admin)
// @Description Generate dan download Excel report untuk klaim
// @Tags Admin - Reports
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security BearerAuth
// @Param from query string false "Tanggal mulai (YYYY-MM-DD)"
// @Param to query string false "Tanggal akhir (YYYY-MM-DD)"
// @Param status query string false "Filter status klaim"
// @Success 200 {file} file "Excel file"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or claims_officer role"
// @Router /reports/claims/excel [get]

// CustomerListExcel godoc
// @Summary Export customer list Excel (Admin)
// @Description Generate dan download Excel report daftar customer
// @Tags Admin - Reports
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security BearerAuth
// @Success 200 {file} file "Excel file"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin role"
// @Router /reports/customers/excel [get]

// AnalyticsSummaryExcel godoc
// @Summary Export analytics summary Excel (Admin)
// @Description Generate dan download Excel report analytics summary
// @Tags Admin - Reports
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security BearerAuth
// @Param from query string false "Tanggal mulai (YYYY-MM-DD)"
// @Param to query string false "Tanggal akhir (YYYY-MM-DD)"
// @Success 200 {file} file "Excel file"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin, underwriter, or finance role"
// @Router /reports/analytics/excel [get]

// AnalyzeApplicationRisk godoc
// @Summary Analisis risiko fraud aplikasi (Admin)
// @Description Melakukan analisis fraud risk menggunakan AI untuk aplikasi tertentu
// @Tags Admin - Fraud Detection
// @Produce json
// @Security BearerAuth
// @Param id path string true "Application ID"
// @Success 200 {object} object{message=string,data=domain.FraudAnalysisResult} "Risk analysis completed"
// @Failure 400 {object} domain.ErrorResponse "application id is required"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or underwriter role"
// @Router /admin/applications/{id}/analyze-risk [post]

// GetHighRiskApplications godoc
// @Summary Get aplikasi berisiko tinggi (Admin)
// @Description Mendapatkan daftar aplikasi dengan risk score tinggi
// @Tags Admin - Fraud Detection
// @Produce json
// @Security BearerAuth
// @Param min_score query int false "Minimum risk score" default(61)
// @Param limit query int false "Jumlah data" default(20)
// @Success 200 {object} object{data=[]domain.FraudAnalysisResult,total=int} "High risk applications"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 403 {object} domain.ErrorResponse "Forbidden - requires super_admin or underwriter role"
// @Router /admin/fraud/high-risk [get]

// ExtractIDCard godoc
// @Summary OCR ekstraksi KTP
// @Description Ekstrak data dari foto KTP menggunakan OCR
// @Tags OCR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.OCRExtractRequest true "Base64 encoded image"
// @Success 200 {object} domain.OCRExtractResponse "Data KTP berhasil diekstrak"
// @Failure 400 {object} domain.ErrorResponse "Invalid request atau image tidak valid"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Router /ocr/extract [post]
