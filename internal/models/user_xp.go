package models

import (
	"time"

	"github.com/google/uuid"
)

type UserXp struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"userId"`
	User      *User     `gorm:"foreignKey:UserID" json:"-"`
	Xp        int       `gorm:"not null;default:0" json:"xp"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
