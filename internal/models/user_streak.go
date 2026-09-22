package models

import (
	"time"

	"github.com/google/uuid"
)

type UserStreak struct {
	UserID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"userId"`
	User             *User      `gorm:"foreignKey:UserID" json:"-"`
	CurrentStreak    int        `gorm:"not null;default:0" json:"currentStreak"`
	LongestStreak    int        `gorm:"not null;default:0" json:"longestStreak"`
	LastActivityAt   *time.Time `json:"lastActivityAt,omitempty"`
	LastActivityDate string     `gorm:"size:10" json:"lastActivityDate,omitempty"`
	FreezesAvailable int        `gorm:"not null;default:0" json:"freezesAvailable"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
}
