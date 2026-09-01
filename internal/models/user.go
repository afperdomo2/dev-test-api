package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const DefaultDailyImportLimit = 200

const DefaultDailyAiLimit = 20

type User struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email            string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash     string         `gorm:"not null" json:"-"`
	IsAdmin          bool           `gorm:"default:false" json:"isAdmin"`
	DailyImportLimit int            `gorm:"not null;default:200" json:"dailyImportLimit"`
	DailyAiLimit     int            `gorm:"not null;default:20" json:"dailyAiLimit"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
