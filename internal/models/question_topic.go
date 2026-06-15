package models

import (
	"time"

	"github.com/google/uuid"
)

type QuestionTopic struct {
	QuestionID uuid.UUID `gorm:"type:uuid;primaryKey" json:"question_id"`
	TopicID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"topic_id"`
	CreatedAt  time.Time `json:"created_at"`
}
