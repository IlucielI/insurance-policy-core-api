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
	MidtransServerKey string
	MidtransClientKey string
	MidtransIsProd    bool
	RecaptchaSecret   string
	CORSOrigins       string
	SentryDSN         string
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
		MidtransServerKey: getEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransClientKey: getEnv("MIDTRANS_CLIENT_KEY", ""),
		MidtransIsProd:    getEnv("MIDTRANS_IS_PRODUCTION", "false") == "true",
		RecaptchaSecret:   getEnv("RECAPTCHA_SECRET", ""),
		CORSOrigins:       getEnv("CORS_ORIGINS", "*"),
		SentryDSN:         getEnv("SENTRY_DSN", ""),
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
