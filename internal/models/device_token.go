package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceToken es un dispositivo registrado para recibir push de Expo. El token
// es único global: si el mismo teléfono inicia sesión con otro usuario, el
// upsert lo reasigna.
type DeviceToken struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	User       *User     `gorm:"foreignKey:UserID" json:"-"`
	Token      string    `gorm:"not null;size:255" json:"token"`
	Platform   string    `gorm:"not null;size:10" json:"platform"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	LastSeenAt time.Time `json:"lastSeenAt"`
}

func (d *DeviceToken) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
