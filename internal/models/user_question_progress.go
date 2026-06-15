package models

import (
	"time"

	"github.com/google/uuid"
)

type UserQuestionProgress struct {
	UserID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"user_id"`
	QuestionID     uuid.UUID  `gorm:"type:uuid;primaryKey" json:"question_id"`
	User           *User      `gorm:"foreignKey:UserID" json:"-"`
	Question       *Question  `gorm:"foreignKey:QuestionID" json:"-"`
	Repetitions    int        `gorm:"not null;default:0" json:"repetitions"`
	EaseFactor     float64    `gorm:"not null;default:2.5" json:"ease_factor"`
	IntervalDays   int        `gorm:"not null;default:0" json:"interval_days"`
	NextReviewAt   *time.Time `json:"next_review_at,omitempty"`
	LastReviewedAt *time.Time `json:"last_reviewed_at,omitempty"`
	IsSaved        bool       `gorm:"not null;default:false" json:"is_saved"`
	IsMastered     bool       `gorm:"not null;default:false" json:"is_mastered"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
