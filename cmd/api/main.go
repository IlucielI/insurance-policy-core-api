package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/IlucielI/insurance-policy-core-api/config"
	_ "github.com/IlucielI/insurance-policy-core-api/docs"
	"github.com/IlucielI/insurance-policy-core-api/internal/delivery/http"
	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/cache"
	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/email"
	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/fraud"
	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/llm"
	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/payment"
	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/storage"
	ws "github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/websocket"
	"github.com/IlucielI/insurance-policy-core-api/internal/middleware"
	"github.com/IlucielI/insurance-policy-core-api/internal/repository"
	"github.com/IlucielI/insurance-policy-core-api/internal/service"
	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Insurance Policy Core API
// @version 1.0
// @description API komprehensif untuk sistem manajemen asuransi dengan fitur aplikasi, polis, klaim, billing, dokumen, notifikasi, fraud detection, OCR, dan analytics
// @description
// @description **Fitur Utama:**
// @description - Autentikasi & Otorisasi berbasis JWT dengan RBAC (Role-Based Access Control)
// @description - Manajemen Produk Asuransi dengan semantic search
// @description - Aplikasi Asuransi dengan AI review dan fraud detection
// @description - Manajemen Polis (endorsement, renewal)
// @description - Klaim dengan timeline tracking
// @description - Billing & Payment terintegrasi Midtrans
// @description - Document Management dengan MinIO
// @description - OCR untuk ekstraksi KTP
// @description - Real-time Notifications via WebSocket
// @description - Analytics & Reporting
// @description - Audit Logging
// @description
// @description **RBAC Roles:**
// @description - `customer`: Akses dasar untuk pengguna
// @description - `underwriter`: Approve aplikasi, risk analysis
// @description - `claims_officer`: Manajemen klaim
// @description - `finance`: Billing & invoice management
// @description - `super_admin`: Full access
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@insurance.example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Authorization header menggunakan Bearer scheme. Format: "Bearer {token}"
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name auth_token
// @description JWT token stored in HTTP-only cookie

func main() {
	// Load config
	cfg := config.Load()

	// Database connection
	dbURL := cfg.DatabaseURL
	if dbURL == "" {
		dbURL = "postgres://postgres:***@localhost:5432/insurance_policy?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("✅ Database connected")

	// Initialize Redis client (optional for caching)
	var redisClient *cache.RedisClient
	if cfg.RedisURL != "" {
		var err error
		redisClient, err = cache.NewRedisClient(cfg.RedisURL)
		if err != nil {
			log.Printf("⚠️  Redis connection failed: %v - running without cache", err)
			redisClient = nil
		} else {
			log.Println("✅ Redis cache connected")
		}
	} else {
		log.Println("⚠️  REDIS_URL not set, caching disabled")
	}

	// Initialize Midtrans client
	if cfg.MidtransServerKey == "" {
		log.Println("⚠️  MIDTRANS_SERVER_KEY not set, payment gateway disabled")
	}
	midtransClient := payment.NewMidtransClient(cfg.MidtransServerKey, cfg.MidtransIsProd)

	// Initialize SMTP email client
	var emailClient *email.SMTPClient
	if cfg.SMTPHost != "" && cfg.SMTPUser != "" {
		smtpPort, _ := strconv.Atoi(cfg.SMTPPort)
		if smtpPort == 0 {
			smtpPort = 587
		}
		emailClient = email.NewSMTPClient(&email.SMTPConfig{
			Host:     cfg.SMTPHost,
			Port:     smtpPort,
			User:     cfg.SMTPUser,
			Password: cfg.SMTPPassword,
			From:     cfg.SMTPFrom,
		})
		log.Printf("✅ SMTP email service configured: %s:%d", cfg.SMTPHost, smtpPort)
	} else {
		log.Println("⚠️  SMTP not configured, email notifications disabled")
	}

	// Initialize MinIO client
	minioClient, err := storage.NewMinIOClient(storage.MinIOConfig{
		Endpoint:  cfg.MinIOEndpoint,
		AccessKey: cfg.MinIOAccessKey,
		SecretKey: cfg.MinIOSecretKey,
		Bucket:    cfg.MinIOBucket,
		UseSSL:    cfg.MinIOUseSSL,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	// Initialize repositories
	productRepo := repository.NewProductRepository(db, redisClient)
	userRepo := repository.NewUserRepository(db, redisClient)
	roleRepo := repository.NewRoleRepository(db)
	applicationRepo := repository.NewApplicationRepository(db)
	chatRepo := repository.NewChatRepository(db)
	policyRepo := repository.NewPolicyRepository(db, redisClient)
	claimRepo := repository.NewClaimRepository(db)
	billingRepo := repository.NewBillingRepository(db)
	documentRepo := repository.NewDocumentRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	notificationPreferencesRepo := repository.NewNotificationPreferencesRepository(db)
	fraudRepo := repository.NewFraudRepository(db)
	reportRepo := repository.NewReportRepository(db)

	// Initialize WebSocket hub
	notificationHub := ws.NewHub()
	go notificationHub.Run()
	log.Println("✅ WebSocket notification hub running")

	// Initialize embeddings client for semantic search
	llmAPIKey := os.Getenv("LLM_API_KEY")
	if llmAPIKey == "" {
		llmAPIKey = "sk-omniroute-default" // Fallback
	}
	embeddingsClient := llm.NewEmbeddingsClient(cfg.LLMBaseURL, cfg.EmbeddingsModel, llmAPIKey)
	log.Printf("🔍 Embeddings Config: %s (model: %s)", cfg.LLMBaseURL, cfg.EmbeddingsModel)

	// Initialize usecases
	productUsecase := usecase.NewProductUsecase(productRepo, embeddingsClient)
	authUsecase := usecase.NewAuthUsecase(userRepo, roleRepo, emailClient)
	applicationUsecase := usecase.NewApplicationUsecase(applicationRepo, productRepo, emailClient)
	policyUsecase := usecase.NewPolicyUsecase(policyRepo, productRepo)
	claimUsecase := usecase.NewClaimUsecase(claimRepo, policyRepo, emailClient)
	billingUsecase := usecase.NewBillingUsecase(billingRepo, policyRepo, userRepo, midtransClient)
	documentUsecase := usecase.NewDocumentUsecase(documentRepo, minioClient)
	notificationUsecase := usecase.NewNotificationUsecase(notificationRepo, notificationHub)
	notificationPreferencesUsecase := usecase.NewNotificationPreferencesUsecase(notificationPreferencesRepo)
	analyticsUsecase := usecase.NewAnalyticsUsecase(analyticsRepo)
	reportUsecase := usecase.NewReportUsecase(reportRepo, userRepo, billingRepo, policyRepo, claimRepo)

	// Inject notification service into usecases for real-time push
	applicationUsecase.SetNotificationService(notificationUsecase)
	claimUsecase.SetNotificationService(notificationUsecase)
	billingUsecase.SetNotificationService(notificationUsecase)
	log.Println("🔔 Real-time notification hooks wired")
	
	// Initialize LLM service with OmniRoute
	llmBaseURL := os.Getenv("LLM_BASE_URL")
	if llmBaseURL == "" {
		llmBaseURL = "http://100.103.220.104:20128/v1"
	}
	
	llmModel := os.Getenv("LLM_MODEL")
	if llmModel == "" {
		llmModel = "claude-sonnet-4.5"
	}
	
	log.Printf("🤖 LLM Config: %s (model: %s)", llmBaseURL, llmModel)
	// Use OmniRouteClient for chat (implements full LLMService interface)
	llmChatClient := llm.NewOmniRouteClient(llmBaseURL, llmModel)
	chatUsecase := usecase.NewChatUsecase(chatRepo, llmChatClient)

	// Initialize fraud detection
	llmClient := llm.NewOmniRouteClient(llmBaseURL, llmModel)
	fraudDetector := fraud.NewFraudDetector(llmClient)
	fraudUsecase := usecase.NewFraudUsecase(fraudDetector, fraudRepo, applicationRepo, productRepo)
	log.Println("🛡️  Fraud detection service initialized")

	// Initialize analytics
	analyticsHandler := http.NewAnalyticsHandler(analyticsUsecase)
	reportHandler := http.NewReportHandler(reportUsecase)

	// Initialize OCR service
	ocrService := service.NewOCRService()
	log.Println("🔍 OCR service initialized")

	// Initialize handlers
	productHandler := http.NewProductHandler(productUsecase)
	authHandler := http.NewAuthHandler(authUsecase)
	applicationHandler := http.NewApplicationHandler(applicationUsecase, auditRepo)
	chatHandler := http.NewChatHandler(chatUsecase)
	aiHandler := http.NewAIHandler(applicationUsecase, llmChatClient)
	fraudHandler := http.NewFraudHandler(fraudUsecase)
	policyHandler := http.NewPolicyHandler(policyUsecase)
	claimHandler := http.NewClaimHandler(claimUsecase, auditRepo)
	billingHandler := http.NewBillingHandler(billingUsecase)
	documentHandler := http.NewDocumentHandler(documentUsecase)
	ocrHandler := http.NewOCRHandler(ocrService)
	auditHandler := http.NewAuditHandler(auditRepo)
	notificationHandler := http.NewNotificationHandler(notificationUsecase, notificationHub)
	notificationPreferencesHandler := http.NewNotificationPreferencesHandler(notificationPreferencesUsecase)

	// Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Insurance Policy API v1.0",
		ServerHeader: "Fiber",
	})

	// Initialize rate limiter (100 req/min for public, 1000 req/min for authenticated)
	middleware.InitRateLimiter(100, 1000)

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	
	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:3000,http://localhost:3001,https://insurance-app.bayuanugerah.my.id,https://insurance-app-cms.bayuanugerah.my.id"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))
	app.Use(middleware.RateLimit()) // Apply rate limiting to all routes

	// Swagger UI documentation
	app.Get("/api/docs/*", fiberSwagger.WrapHandler)

	// Health check endpoints
	healthHandler := http.NewHealthHandler(db, redisClient, minioClient, cfg.LLMBaseURL)
	app.Get("/health", healthHandler.Liveness)           // Liveness probe
	app.Get("/health/ready", healthHandler.Readiness)    // Readiness probe
	app.Get("/metrics", healthHandler.Metrics)           // Metrics endpoint

	// API routes
	api := app.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)
	auth.Post("/forgot-password", authHandler.ForgotPassword)
	auth.Post("/reset-password", authHandler.ResetPassword)
	auth.Get("/me", authHandler.Me)

	// Products routes
	products := api.Group("/products")
	products.Get("/", productHandler.ListProducts)
	products.Get("/:id", productHandler.GetProduct)
	products.Post("/:id/calculate-premium", productHandler.CalculatePremium)
	products.Post("/search", productHandler.SemanticSearch)

	// Applications routes
	applications := api.Group("/applications")
	applications.Post("/", applicationHandler.CreateApplication)
	applications.Get("/:id", applicationHandler.GetApplication)

	// Chat routes
	chat := api.Group("/chat")
	chat.Post("/", chatHandler.SendMessage)

	// Customer App - Policy Management routes
	policies := api.Group("/policies")
	policies.Get("/", policyHandler.ListPolicies)
	policies.Get("/:id", policyHandler.GetPolicy)
	policies.Post("/:id/endorse", policyHandler.EndorsePolicy)
	policies.Post("/:id/renew", policyHandler.RenewPolicy)

	// Customer App - Claims routes
	claims := api.Group("/claims")
	claims.Post("/", claimHandler.CreateClaim)
	claims.Get("/:id", claimHandler.GetClaim)
	claims.Put("/:id/documents", claimHandler.UploadDocument)
	claims.Get("/:id/timeline", claimHandler.GetClaimTimeline)

	// Customer App - Billing routes
	billing := api.Group("/billing")
	billing.Get("/invoices", billingHandler.GetInvoices)
	billing.Post("/pay", billingHandler.PayInvoice)
	billing.Get("/history", billingHandler.GetPaymentHistory)
	billing.Get("/verify/:order_id", billingHandler.VerifyPayment) // Security: verify payment from DB, not URL params

	// Webhook routes (no auth required)
	webhooks := api.Group("/webhooks")
	webhooks.Post("/payment", billingHandler.HandlePaymentWebhook)

	// Customer App - Documents routes
	documents := api.Group("/documents")
	documents.Get("/", documentHandler.ListDocuments)
	documents.Post("/upload", documentHandler.UploadDocument)
	documents.Get("/:id/download", documentHandler.DownloadDocument)
	documents.Delete("/:id", documentHandler.DeleteDocument)

	// OCR routes
	ocr := api.Group("/ocr")
	ocr.Post("/extract", middleware.AuthRequired(), ocrHandler.ExtractIDCard)

	// Reports routes (admin-only, accessible by super_admin, finance, claims_officer, underwriter)
	reports := api.Group("/reports", middleware.AuthRequired())
	reports.Get("/billing/:id/pdf", middleware.RequireRole("super_admin", "finance"), reportHandler.BillingStatementPDF)
	reports.Get("/claims/excel", middleware.RequireRole("super_admin", "claims_officer"), reportHandler.ClaimsReportExcel)
	reports.Get("/customers/excel", middleware.RequireRole("super_admin"), reportHandler.CustomerListExcel)
	reports.Get("/analytics/excel", middleware.RequireRole("super_admin", "underwriter", "finance"), reportHandler.AnalyticsSummaryExcel)

	// Admin routes with RBAC protection
	admin := api.Group("/admin", middleware.AuthRequired())

	// Super admin + underwriter: applications management
	admin.Get("/applications", middleware.RequireRole("super_admin", "underwriter"), applicationHandler.ListApplications)
	admin.Put("/applications/:id/status", middleware.RequireRole("super_admin", "underwriter"), applicationHandler.UpdateStatus)
	admin.Post("/applications/bulk-status", middleware.RequireRole("super_admin", "underwriter"), applicationHandler.BulkUpdateStatus)

	// Admin email notification routes
	adminEmailHandler := http.NewAdminEmailHandler(applicationUsecase, claimUsecase, userRepo, auditRepo)
	admin.Post("/email/application-approved/:id", middleware.RequireRole("super_admin", "underwriter"), adminEmailHandler.SendApplicationApprovedEmail)
	admin.Post("/email/claim-status/:id", middleware.RequireRole("super_admin", "claims_officer"), adminEmailHandler.SendClaimStatusEmail)
	admin.Post("/email/bulk-send", middleware.RequireRole("super_admin"), adminEmailHandler.BulkSendEmail)
	admin.Get("/email/preview/:type", middleware.RequireRole("super_admin", "underwriter", "claims_officer"), adminEmailHandler.PreviewEmail)

	// AI routes (super_admin + underwriter)
	admin.Post("/ai-review/:id", middleware.RequireRole("super_admin", "underwriter"), aiHandler.ReviewApplication)

	// Fraud detection routes (super_admin + underwriter)
	admin.Post("/applications/:id/analyze-risk", middleware.RequireRole("super_admin", "underwriter"), fraudHandler.AnalyzeApplicationRisk)
	admin.Get("/fraud/high-risk", middleware.RequireRole("super_admin", "underwriter"), fraudHandler.GetHighRiskApplications)

	// Claims management (super_admin + claims_officer)
	admin.Get("/claims", middleware.RequireRole("super_admin", "claims_officer"), claimHandler.ListClaims)
	admin.Put("/claims/:id/status", middleware.RequireRole("super_admin", "claims_officer"), claimHandler.UpdateClaimStatus)
	admin.Put("/claims/:id/approve", middleware.RequireRole("super_admin", "claims_officer"), claimHandler.ApproveClaim)

	// Analytics routes (accessible by super_admin, underwriter, claims_officer, finance)
	analytics := admin.Group("/analytics")
	analytics.Get("/dashboard", analyticsHandler.GetDashboard)
	analytics.Get("/revenue", analyticsHandler.GetRevenue)
	analytics.Get("/claims", analyticsHandler.GetClaimsStatus)
	analytics.Get("/products", analyticsHandler.GetTopProducts)
	admin.Post("/claims/bulk-status", middleware.RequireRole("super_admin", "claims_officer"), claimHandler.BulkUpdateClaimStatus)

	// Audit logs routes (super_admin only)
	auditRoutes := admin.Group("/audit-logs")
	auditRoutes.Get("/actions", middleware.RequireRole("super_admin", "underwriter", "claims_officer"), auditHandler.GetAvailableActions)
	auditRoutes.Get("/entity-types", middleware.RequireRole("super_admin", "underwriter", "claims_officer"), auditHandler.GetEntityTypes)
	auditRoutes.Get("/entity/:type/:id", middleware.RequireRole("super_admin", "underwriter", "claims_officer"), auditHandler.GetEntityAuditLogs)
	auditRoutes.Get("/:id", middleware.RequireRole("super_admin"), auditHandler.GetAuditLog)
	auditRoutes.Get("/", middleware.RequireRole("super_admin"), auditHandler.ListAuditLogs)

	// Notifications routes
	notifications := api.Group("/notifications")
	notifications.Get("/ws", notificationHandler.WebSocketHandler)
	notifications.Get("/", notificationHandler.GetNotifications)
	notifications.Get("/unread-count", notificationHandler.GetUnreadCount)
	notifications.Put("/:id/read", notificationHandler.MarkAsRead)
	notifications.Put("/read-all", notificationHandler.MarkAllAsRead)
	notifications.Delete("/:id", notificationHandler.DeleteNotification)

	// Notification Preferences routes
	api.Get("/notification-preferences", middleware.AuthRequired(), notificationPreferencesHandler.GetPreferences)
	api.Put("/notification-preferences", middleware.AuthRequired(), notificationPreferencesHandler.UpdatePreferences)

	// Billing management (super_admin + finance)
	admin.Get("/billing/invoices", middleware.RequireRole("super_admin", "finance"), billingHandler.GetAllInvoices)
	admin.Post("/billing/invoices", middleware.RequireRole("super_admin", "finance"), billingHandler.CreateInvoice)
	admin.Put("/billing/invoices/:id/status", middleware.RequireRole("super_admin", "finance"), billingHandler.UpdateInvoiceStatus)

	// User & role management (super_admin only)
	admin.Get("/users", middleware.RequireRole("super_admin"), authHandler.ListUsers)
	admin.Post("/users/:id/roles", middleware.RequireRole("super_admin"), authHandler.AssignRole)
	admin.Delete("/users/:id/roles", middleware.RequireRole("super_admin"), authHandler.RemoveRole)
	admin.Get("/roles", middleware.RequireRole("super_admin"), authHandler.GetRoles)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
