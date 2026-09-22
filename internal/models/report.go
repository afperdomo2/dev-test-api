package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Report struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ReporterUserID uuid.UUID `gorm:"type:uuid;not null;index" json:"reporterUserId"`
	Reporter       *User     `gorm:"foreignKey:ReporterUserID" json:"-"`
	TargetType     string    `gorm:"not null;size:20" json:"targetType"`
	TargetID       uuid.UUID `gorm:"type:uuid;not null;index" json:"targetId"`
	Reason         string    `gorm:"not null;type:text" json:"reason"`
	Status         string    `gorm:"not null;default:pending;size:20" json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (r *Report) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
