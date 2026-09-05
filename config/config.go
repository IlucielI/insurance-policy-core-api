package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DatabaseURL       string
	RedisURL          string
	SessionSecret     string
	SessionMaxAge     int
	LLMEndpoint       string
	LLMModel          string
	LLMBaseURL        string
	EmbeddingsModel   string
	MidtransServerKey string
	MidtransClientKey string
	MidtransIsProd    bool
	RecaptchaSecret   string
	CORSOrigins       string
	SentryDSN         string
	SMTPHost          string
	SMTPPort          string
	SMTPUser          string
	SMTPPassword      string
	SMTPFrom          string
	MinIOEndpoint     string
	MinIOAccessKey    string
	MinIOSecretKey    string
	MinIOBucket       string
	MinIOUseSSL       bool
}

func Load() *Config {
	// Load .env if exists (local dev)
	_ = godotenv.Load()

	return &Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		RedisURL:          getEnv("REDIS_URL", ""),
		SessionSecret:     getEnv("SESSION_SECRET", "change-me-in-production"),
		SessionMaxAge:     86400,
		LLMEndpoint:       getEnv("LLM_ENDPOINT", ""),
		LLMModel:          getEnv("LLM_MODEL", "gpt-4o-mini"),
		LLMBaseURL:        getEnv("LLM_BASE_URL", "http://100.103.220.104:20128/v1"),
		EmbeddingsModel:   getEnv("EMBEDDINGS_MODEL", "bge-m3"),
		MidtransServerKey: getEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransClientKey: getEnv("MIDTRANS_CLIENT_KEY", ""),
		MidtransIsProd:    getEnv("MIDTRANS_IS_PRODUCTION", "false") == "true",
		RecaptchaSecret:   getEnv("RECAPTCHA_SECRET", ""),
		CORSOrigins:       getEnv("CORS_ORIGINS", "*"),
		SentryDSN:         getEnv("SENTRY_DSN", ""),
		SMTPHost:          getEnv("SMTP_HOST", ""),
		SMTPPort:          getEnv("SMTP_PORT", "587"),
		SMTPUser:          getEnv("SMTP_USER", ""),
		SMTPPassword:      getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:          getEnv("SMTP_FROM", ""),
		MinIOEndpoint:     getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey:    getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey:    getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinIOBucket:       getEnv("MINIO_BUCKET", "insurance-documents"),
		MinIOUseSSL:       getEnv("MINIO_USE_SSL", "false") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue == "" && key != "SENTRY_DSN" {
			log.Printf("Warning: %s not set, using empty value", key)
		}
		return defaultValue
	}
	return value
}
