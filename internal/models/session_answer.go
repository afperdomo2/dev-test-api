package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionAnswer struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SessionID       uuid.UUID `gorm:"type:uuid;not null;index" json:"session_id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	QuestionID      uuid.UUID `gorm:"type:uuid;not null;index" json:"question_id"`
	Session         *Session  `gorm:"foreignKey:SessionID" json:"-"`
	User            *User     `gorm:"foreignKey:UserID" json:"-"`
	Question        *Question `gorm:"foreignKey:QuestionID" json:"-"`
	AnswerText      string    `gorm:"type:text" json:"answer_text,omitempty"`
	SelectedOptions string    `gorm:"type:jsonb;default:'[]'::jsonb" json:"selected_options,omitempty"`
	IsCorrect       bool      `gorm:"not null;default:false" json:"is_correct"`
	AiFeedback      string    `gorm:"type:text" json:"ai_feedback,omitempty"`
	ResponseTimeMs  int       `json:"response_time_ms,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

func (a *SessionAnswer) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
