package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tip struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId"`
	User      *User          `gorm:"foreignKey:UserID" json:"-"`
	TopicID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"topicId"`
	Topic     *Topic         `gorm:"foreignKey:TopicID" json:"-"`
	Content   string         `gorm:"not null;type:text" json:"content"`
	Source    string         `gorm:"not null;default:manual;size:20" json:"source"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Tip) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
