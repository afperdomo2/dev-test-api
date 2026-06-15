package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Question struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	User       *User          `gorm:"foreignKey:UserID" json:"-"`
	Content    string         `gorm:"not null" json:"content"`
	Answer     string         `json:"answer,omitempty"`
	Topic      string         `json:"topic,omitempty"`
	ReviewedAt *time.Time     `json:"reviewed_at,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}
