package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type CorsConfig struct {
	AllowedOrigins string
}

type AIConfig struct {
	APIURL         string
	APIKey         string
	Model          string
	RequestTimeout int
}

type Config struct {
	Port    string
	GinMode string
	DB      DBConfig
	JWT     JWTConfig
	Cors    CorsConfig
	AI      AIConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret    string
	ExpiryHrs string
}

func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

func (j JWTConfig) SecretBytes() []byte {
	return []byte(j.Secret)
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "release"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "dev_test_api"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", ""),
			ExpiryHrs: getEnv("JWT_EXPIRY_HOURS", "24"),
		},
		Cors: CorsConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"),
		},
		AI: AIConfig{
			APIURL:         getEnv("AI_API_URL", ""),
			APIKey:         getEnv("AI_API_KEY", ""),
			Model:          getEnv("AI_MODEL", "gpt-4o-mini"),
			RequestTimeout: getIntEnv("AI_REQUEST_TIMEOUT", 120),
		},
	}

	os.Setenv("GIN_MODE", cfg.GinMode)

	if cfg.JWT.Secret == "" {
		log.Fatal("❌ JWT_SECRET environment variable is required")
	}
	if cfg.DB.Password == "" {
		log.Fatal("❌ DB_PASSWORD environment variable is required")
	}

	if cfg.AI.APIURL != "" && cfg.AI.APIKey != "" {
		log.Printf("🤖 AI configured: url=%s model=%s", cfg.AI.APIURL, cfg.AI.Model)
	} else {
		log.Println("⚠️ AI not configured — question generation disabled. Set AI_API_URL and AI_API_KEY.")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
