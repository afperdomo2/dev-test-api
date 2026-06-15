package database

import (
	"log"

	"github.com/felipe/dev-test-api/internal/config"
	"github.com/felipe/dev-test-api/internal/models"
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
		Logger:                                   logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: false,
	})
	if err != nil {
		log.Fatalf("❌ failed to connect to database: %v", err)
	}

	log.Println("🛢️ Connected to PostgreSQL")

	if err := db.AutoMigrate(
		&models.User{},
		&models.Topic{},
		&models.Question{},
		&models.QuestionOption{},
		&models.CodeChallenge{},
		&models.QuestionTopic{},
		&models.UserQuestionProgress{},
		&models.Session{},
		&models.SessionTopic{},
		&models.SessionAnswer{},
	); err != nil {
		log.Fatalf("❌ failed to run auto-migration: %v", err)
	}

	log.Println("✅ Database migration completed")

	if err := runDDL(db); err != nil {
		log.Fatalf("❌ failed to run DDL: %v", err)
	}

	seedDefaultTopics(db)

	return db
}
