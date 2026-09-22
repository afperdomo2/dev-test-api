package models

import (
	"time"

	"github.com/google/uuid"
)

type UserTopicXp struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"userId"`
	User      *User     `gorm:"foreignKey:UserID" json:"-"`
	TopicID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"topicId"`
	Topic     *Topic    `gorm:"foreignKey:TopicID" json:"-"`
	Xp        int       `gorm:"not null;default:0" json:"xp"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
