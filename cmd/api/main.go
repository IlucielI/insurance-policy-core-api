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

	// Initialize usecases
	productUsecase := usecase.NewProductUsecase(productRepo)
	authUsecase := usecase.NewAuthUsecase(userRepo)
	applicationUsecase := usecase.NewApplicationUsecase(applicationRepo, productRepo)
	
	// Initialize LLM service (mock for now, will integrate real LLM later)
	llmService := service.NewMockLLMService()
	chatUsecase := usecase.NewChatUsecase(chatRepo, llmService)

	// Initialize handlers
	productHandler := http.NewProductHandler(productUsecase)
	authHandler := http.NewAuthHandler(authUsecase)
	chatHandler := http.NewChatHandler(chatUsecase)
	
	// TODO: Initialize application handler
	_ = applicationUsecase // Will be used when handler is implemented

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
		corsOrigins = "*"
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

	// Applications routes (TODO: implement handler)
	applications := api.Group("/applications")
	applications.Post("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Create application endpoint"})
	})
	applications.Get("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Get application endpoint"})
	})

	// Chat routes
	chat := api.Group("/chat")
	chat.Post("/", chatHandler.SendMessage)

	// Admin routes (TODO: add auth middleware)
	admin := api.Group("/admin")
	admin.Get("/applications", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "List applications endpoint"})
	})
	admin.Patch("/applications/:id/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Update application status endpoint"})
	})

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
