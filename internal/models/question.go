package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Question struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID" json:"-"`
	Type        string         `gorm:"not null;size:30" json:"type"`
	Content     string         `gorm:"not null;type:text" json:"content"`
	Explanation string         `gorm:"type:text" json:"explanation,omitempty"`
	Difficulty  string         `gorm:"not null;default:intermediate;size:20" json:"difficulty"`
	IsPublic    bool           `gorm:"not null;default:false" json:"is_public"`
	Source      string         `gorm:"not null;default:ai_generated;size:20" json:"source"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Options        []QuestionOption `gorm:"foreignKey:QuestionID" json:"options,omitempty"`
	CodeChallenge  *CodeChallenge   `gorm:"foreignKey:QuestionID" json:"code_challenge,omitempty"`
	Topics         []Topic          `gorm:"many2many:question_topics;" json:"-"`
	QuestionTopics []QuestionTopic  `gorm:"foreignKey:QuestionID" json:"-"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}
