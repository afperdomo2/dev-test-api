package models

import (
	"time"

	"github.com/google/uuid"
)

type QuestionTopic struct {
	QuestionID uuid.UUID `gorm:"type:uuid;primaryKey" json:"questionId"`
	TopicID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"topicId"`
	CreatedAt  time.Time `json:"createdAt"`
}
