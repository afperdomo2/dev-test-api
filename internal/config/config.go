package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port    string
	GinMode string
	DB      DBConfig
	JWT     JWTConfig
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
	}

	os.Setenv("GIN_MODE", cfg.GinMode)

	if cfg.JWT.Secret == "" {
		log.Fatal("❌ JWT_SECRET environment variable is required")
	}
	if cfg.DB.Password == "" {
		log.Fatal("❌ DB_PASSWORD environment variable is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
