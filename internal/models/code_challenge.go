package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CodeChallenge struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	QuestionID     uuid.UUID `gorm:"type:uuid;uniqueIndex;not null" json:"question_id"`
	StarterCode    string    `gorm:"type:text" json:"starter_code,omitempty"`
	ExpectedOutput string    `gorm:"type:text" json:"expected_output,omitempty"`
	Language       string    `gorm:"not null;size:50" json:"language"`
	TestCasesJSON  string    `gorm:"type:jsonb;default:'[]'::jsonb" json:"test_cases,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (c *CodeChallenge) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
