package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionAnswer struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SessionID       uuid.UUID `gorm:"type:uuid;not null;index" json:"sessionId"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	QuestionID      uuid.UUID `gorm:"type:uuid;not null;index" json:"questionId"`
	Session         *Session  `gorm:"foreignKey:SessionID" json:"-"`
	User            *User     `gorm:"foreignKey:UserID" json:"-"`
	Question        *Question `gorm:"foreignKey:QuestionID" json:"-"`
	AnswerText      string    `gorm:"type:text" json:"answerText,omitempty"`
	SelectedOptions string    `gorm:"type:jsonb" json:"selectedOptions,omitempty"`
	IsCorrect       bool      `gorm:"not null;default:false" json:"isCorrect"`
	AiFeedback      string    `gorm:"type:text" json:"aiFeedback,omitempty"`
	ResponseTimeMs  int       `json:"responseTimeMs,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
}

func (a *SessionAnswer) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.SelectedOptions == "" {
		a.SelectedOptions = "[]"
	}
	return nil
}
