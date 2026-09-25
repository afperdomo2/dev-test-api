package models

import (
	"time"

	"github.com/google/uuid"
)

// Setting es un ajuste tunable de la aplicación (clave/valor). Es la única
// entidad del repo cuya PK es texto en vez de un UUID.
type Setting struct {
	Key         string     `gorm:"primaryKey;size:100" json:"key"`
	Value       string     `gorm:"not null" json:"value"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	UpdatedBy   *uuid.UUID `gorm:"type:uuid" json:"updatedBy,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func (Setting) TableName() string { return "settings" }
