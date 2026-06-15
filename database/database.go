package database

import (
	"log"

	"github.com/felipe/dev-test-api/config"
	"github.com/felipe/dev-test-api/internal/questions"
	"github.com/felipe/dev-test-api/internal/users"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) *gorm.DB {
	logLevel := logger.Info
	if cfg.GinMode == "release" {
		logLevel = logger.Error
	}

	db, err := gorm.Open(postgres.Open(cfg.DB.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("Connected to PostgreSQL")

	if err := db.AutoMigrate(
		&users.User{},
		&questions.Question{},
	); err != nil {
		log.Fatalf("failed to run auto-migration: %v", err)
	}

	log.Println("Database migration completed")

	return db
}
