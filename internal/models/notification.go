package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification es un aviso in-app para un usuario. `creation_day` no se mapea:
// la calcula Postgres (columna generada) y sostiene la idempotencia del job.
type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	User      *User     `gorm:"foreignKey:UserID" json:"-"`
	Type      string    `gorm:"not null;size:30" json:"type"`
	Title     string    `gorm:"not null;size:200" json:"title"`
	Body      string    `gorm:"not null" json:"body"`
	IsRead    bool      `gorm:"not null;default:false" json:"isRead"`
	CreatedAt time.Time `json:"createdAt"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
