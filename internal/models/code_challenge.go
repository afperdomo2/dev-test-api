package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CodeChallenge struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	QuestionID     uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"questionId"`
	StarterCode    string    `gorm:"type:text" json:"starterCode,omitempty"`
	ExpectedOutput string    `gorm:"type:text" json:"expectedOutput,omitempty"`
	Language       string    `gorm:"not null;size:50" json:"language"`
	TestCasesJSON  string    `gorm:"type:jsonb" json:"testCases,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (c *CodeChallenge) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	if c.TestCasesJSON == "" {
		c.TestCasesJSON = "[]"
	}
	return nil
}
