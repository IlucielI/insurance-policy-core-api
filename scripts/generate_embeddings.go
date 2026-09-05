package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/IlucielI/insurance-policy-core-api/config"
	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/cache"
	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/llm"
	"github.com/IlucielI/insurance-policy-core-api/internal/repository"
	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("🚀 Starting embeddings generation for products...")

	// Load config
	cfg := config.Load()

	// Database connection
	dbURL := cfg.DatabaseURL
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
		if dbURL == "" {
			dbURL = "postgres://postgres:postgres@localhost:5432/insurance_policy?sslmode=disable"
		}
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

	// Initialize embeddings client
	embeddingsClient := llm.NewEmbeddingsClient(cfg.LLMBaseURL, cfg.EmbeddingsModel)
	log.Printf("🔍 Embeddings Config: %s (model: %s)", cfg.LLMBaseURL, cfg.EmbeddingsModel)

	// Initialize repository and usecase
	productRepo := repository.NewProductRepository(db, nil) // No cache needed for script
	productUsecase := usecase.NewProductUsecase(productRepo, embeddingsClient)

	// Generate embeddings
	ctx := context.Background()
	err = productUsecase.GenerateProductEmbeddings(ctx)
	if err != nil {
		log.Fatal("Failed to generate embeddings:", err)
	}

	log.Println("✅ Embeddings generation completed successfully!")
}
