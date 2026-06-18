package models

import (
	"time"

	"github.com/google/uuid"
)

func (UserQuestionProgress) TableName() string {
	return "user_question_progress"
}

type UserQuestionProgress struct {
	UserID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"userId"`
	QuestionID     uuid.UUID  `gorm:"type:uuid;primaryKey" json:"questionId"`
	User           *User      `gorm:"foreignKey:UserID" json:"-"`
	Question       *Question  `gorm:"foreignKey:QuestionID" json:"-"`
	Repetitions    int        `gorm:"not null;default:0" json:"repetitions"`
	EaseFactor     float64    `gorm:"not null;default:2.5" json:"easeFactor"`
	IntervalDays   int        `gorm:"not null;default:0" json:"intervalDays"`
	NextReviewAt   *time.Time `json:"nextReviewAt,omitempty"`
	LastReviewedAt *time.Time `json:"lastReviewedAt,omitempty"`
	IsSaved        bool       `gorm:"not null;default:false" json:"isSaved"`
	IsMastered     bool       `gorm:"not null;default:false" json:"isMastered"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}
