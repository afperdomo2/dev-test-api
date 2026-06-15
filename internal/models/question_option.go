package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionOption struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	QuestionID uuid.UUID `gorm:"type:uuid;not null;index" json:"question_id"`
	Content    string    `gorm:"not null;type:text" json:"content"`
	IsCorrect  bool      `gorm:"not null;default:false" json:"is_correct"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (o *QuestionOption) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}
