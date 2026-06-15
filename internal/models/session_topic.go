package models

import (
	"time"

	"github.com/google/uuid"
)

type SessionTopic struct {
	SessionID uuid.UUID `gorm:"type:uuid;primaryKey" json:"session_id"`
	TopicID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"topic_id"`
	CreatedAt time.Time `json:"created_at"`
}
