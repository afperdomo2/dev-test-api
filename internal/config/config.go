package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// defaultMaxBodyBytes es el límite por defecto para bodies JSON (1 MiB).
const defaultMaxBodyBytes int64 = 1 << 20

type CorsConfig struct {
	AllowedOrigins string
}

type AIConfig struct {
	APIURL         string
	APIKey         string
	Model          string
	RequestTimeout int
}

type OpenCodeConfig struct {
	UsageURL     string
	UsageTimeout int
}

type PushConfig struct {
	ExpoURL     string
	AccessToken string
}

type Config struct {
	Port     string
	GinMode  string
	DB       DBConfig
	JWT      JWTConfig
	Cors     CorsConfig
	AI       AIConfig
	OpenCode OpenCodeConfig
	Push     PushConfig

	// HTTP hardening
	MaxBodyBytes   int64
	TrustedProxies []string
	HSTSEnabled    bool
	SwaggerEnabled bool
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
		OpenCode: OpenCodeConfig{
			UsageURL:     getEnv("OPENCODE_USAGE_URL", "https://opencode.ai/zen/go/v1/usage"),
			UsageTimeout: getIntEnv("OPENCODE_USAGE_TIMEOUT", 10),
		},
		Push: PushConfig{
			ExpoURL:     getEnv("EXPO_PUSH_URL", "https://exp.host/--/api/v2/push/send"),
			AccessToken: getEnv("EXPO_ACCESS_TOKEN", ""),
		},
	}

	cfg.MaxBodyBytes = getInt64Env("MAX_BODY_BYTES", defaultMaxBodyBytes)
	if cfg.MaxBodyBytes <= 0 {
		cfg.MaxBodyBytes = defaultMaxBodyBytes
	}
	cfg.TrustedProxies = splitList(getEnv("TRUSTED_PROXIES", ""))
	cfg.HSTSEnabled = getBoolEnv("HSTS_ENABLED", true)
	cfg.SwaggerEnabled = swaggerEnabled(cfg.GinMode)

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

func getInt64Env(key string, fallback int64) int64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

// splitList convierte "a, b,,c" en ["a","b","c"], descartando vacíos.
func splitList(raw string) []string {
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

// swaggerEnabled resuelve SWAGGER_ENABLED: vacío = habilitado salvo en release.
func swaggerEnabled(ginMode string) bool {
	raw := os.Getenv("SWAGGER_ENABLED")
	if raw == "" {
		return ginMode != "release"
	}

	if enabled, err := strconv.ParseBool(raw); err == nil {
		return enabled
	}
	return ginMode != "release"
}
