package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Topic struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Slug      string     `gorm:"not null;size:100" json:"slug"`
	Name      string     `gorm:"not null;size:150" json:"name"`
	Category  string     `gorm:"not null;size:50" json:"category"`
	IsSystem  bool       `gorm:"not null;default:false" json:"isSystem"`
	CreatedBy *uuid.UUID `gorm:"type:uuid" json:"createdBy,omitempty"`
	Creator   *User      `gorm:"foreignKey:CreatedBy" json:"-"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func (t *Topic) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
