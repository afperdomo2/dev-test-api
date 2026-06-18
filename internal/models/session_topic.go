package models

import (
	"time"

	"github.com/google/uuid"
)

type SessionTopic struct {
	SessionID uuid.UUID `gorm:"type:uuid;primaryKey" json:"sessionId"`
	TopicID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"topicId"`
	CreatedAt time.Time `json:"createdAt"`
}
