package config

import (
	"os"
)

type Config struct {
	Port              string
	OpenAIKey         string
	OpenAIModel       string
	DefaultUserID     string
	SessionSecret     string
	BrevoAPIKey       string
	BrevoSenderEmail  string
	BrevoSenderName   string
	AppBaseURL        string
}

func New() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		OpenAIKey:        os.Getenv("OPENAI_API_KEY"),
		OpenAIModel:      getEnv("OPENAI_MODEL", "gpt-4o"),
		DefaultUserID:    getEnv("DEFAULT_USER_ID", ""),
		SessionSecret:    getEnv("SESSION_SECRET", "change-me-in-production-32chars!!"),
		BrevoAPIKey:      os.Getenv("BREVO_API_KEY"),
		BrevoSenderEmail: getEnv("BREVO_SENDER_EMAIL", "noreply@kalorie.ai"),
		BrevoSenderName:  getEnv("BREVO_SENDER_NAME", "Kalorie AI"),
		AppBaseURL:       getEnv("APP_BASE_URL", "http://localhost:8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
