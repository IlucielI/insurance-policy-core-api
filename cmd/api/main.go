package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "Insurance Policy API v1.0",
		ServerHeader: "Fiber",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:3000,http://localhost:3001"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowCredentials: false,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"service": "insurance-policy-api",
		})
	})

	// API routes
	api := app.Group("/api/v1")
	
	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", func(c *fiber.Ctx) error {
		return c.SendString("Register endpoint")
	})
	auth.Post("/login", func(c *fiber.Ctx) error {
		return c.SendString("Login endpoint")
	})
	auth.Post("/logout", func(c *fiber.Ctx) error {
		return c.SendString("Logout endpoint")
	})

	// Products routes
	products := api.Group("/products")
	products.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("List products")
	})
	products.Get("/:id", func(c *fiber.Ctx) error {
		return c.SendString("Get product by ID")
	})

	// Applications routes
	applications := api.Group("/applications")
	applications.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("Create application")
	})
	applications.Get("/:id", func(c *fiber.Ctx) error {
		return c.SendString("Get application by ID")
	})

	// Chat routes
	chat := api.Group("/chat")
	chat.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("Chat with AI assistant")
	})

	// Admin routes (protected)
	admin := api.Group("/admin")
	admin.Get("/applications", func(c *fiber.Ctx) error {
		return c.SendString("List all applications (admin)")
	})
	admin.Patch("/applications/:id/status", func(c *fiber.Ctx) error {
		return c.SendString("Update application status")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
