package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/IlucielI/insurance-policy-core-api/internal/delivery/http"
	"github.com/IlucielI/insurance-policy-core-api/internal/repository"
	"github.com/IlucielI/insurance-policy-core-api/internal/service"
	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/insurance_policy?sslmode=disable"
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

	// Initialize repositories
	productRepo := repository.NewProductRepository(db)
	userRepo := repository.NewUserRepository(db)
	applicationRepo := repository.NewApplicationRepository(db)
	chatRepo := repository.NewChatRepository(db)
	policyRepo := repository.NewPolicyRepository(db)
	claimRepo := repository.NewClaimRepository(db)
	billingRepo := repository.NewBillingRepository(db)
	documentRepo := repository.NewDocumentRepository(db)

	// Initialize usecases
	productUsecase := usecase.NewProductUsecase(productRepo)
	authUsecase := usecase.NewAuthUsecase(userRepo)
	applicationUsecase := usecase.NewApplicationUsecase(applicationRepo, productRepo)
	policyUsecase := usecase.NewPolicyUsecase(policyRepo, productRepo)
	claimUsecase := usecase.NewClaimUsecase(claimRepo, policyRepo)
	billingUsecase := usecase.NewBillingUsecase(billingRepo, policyRepo)
	documentUsecase := usecase.NewDocumentUsecase(documentRepo)
	
	// Initialize LLM service (mock for now, will integrate real LLM later)
	llmService := service.NewMockLLMService()
	chatUsecase := usecase.NewChatUsecase(chatRepo, llmService)

	// Initialize handlers
	productHandler := http.NewProductHandler(productUsecase)
	authHandler := http.NewAuthHandler(authUsecase)
	chatHandler := http.NewChatHandler(chatUsecase)
	aiHandler := http.NewAIHandler(applicationUsecase, llmService)
	applicationHandler := http.NewApplicationHandler(applicationUsecase)
	policyHandler := http.NewPolicyHandler(policyUsecase)
	claimHandler := http.NewClaimHandler(claimUsecase)
	billingHandler := http.NewBillingHandler(billingUsecase)
	documentHandler := http.NewDocumentHandler(documentUsecase)

	// Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Insurance Policy API v1.0",
		ServerHeader: "Fiber",
	})

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

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "insurance-policy-api",
		})
	})

	// API routes
	api := app.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	// Products routes
	products := api.Group("/products")
	products.Get("/", productHandler.ListProducts)
	products.Get("/:id", productHandler.GetProduct)
	products.Post("/:id/calculate-premium", productHandler.CalculatePremium)

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

	// Customer App - Documents routes
	documents := api.Group("/documents")
	documents.Get("/", documentHandler.ListDocuments)
	documents.Get("/:id/download", documentHandler.DownloadDocument)

	// Admin routes (TODO: add auth middleware)
	admin := api.Group("/admin")
	admin.Get("/applications", applicationHandler.ListApplications)
	admin.Put("/applications/:id/status", applicationHandler.UpdateStatus)
	admin.Post("/applications/bulk-status", applicationHandler.BulkUpdateStatus)
	
	// AI routes
	admin.Post("/ai-review/:id", aiHandler.ReviewApplication)

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
