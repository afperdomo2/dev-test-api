package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID             uuid.UUID  `gorm:"type:uuid;not null;index" json:"userId"`
	User               *User      `gorm:"foreignKey:UserID" json:"-"`
	Name               string     `gorm:"not null;size:200" json:"name"`
	Status             string     `gorm:"not null;default:in_progress;size:20" json:"status"`
	Mode               string     `gorm:"not null;default:generate;size:20" json:"mode"`
	Difficulty         string     `gorm:"not null;default:intermediate;size:20" json:"difficulty"`
	QuestionLimit      *int       `json:"questionLimit,omitempty"`
	QuestionsGenerated int        `gorm:"not null;default:0" json:"questionsGenerated"`
	Score              *float64   `json:"score,omitempty"`
	StartedAt          time.Time  `json:"startedAt"`
	FinishedAt         *time.Time `json:"finishedAt,omitempty"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`

	Topics        []Topic         `gorm:"many2many:session_topics;" json:"-"`
	SessionTopics []SessionTopic  `gorm:"foreignKey:SessionID" json:"-"`
	Answers       []SessionAnswer `gorm:"foreignKey:SessionID" json:"-"`
}

func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.StartedAt.IsZero() {
		s.StartedAt = time.Now()
	}
	return nil
}
