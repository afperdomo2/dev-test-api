package models

import (
	"time"

	"github.com/google/uuid"
)

type GenerationJobTopic struct {
	JobID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"jobId"`
	TopicID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"topicId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (GenerationJobTopic) TableName() string { return "generation_job_topics" }
