package questions

import (
	"time"

	"github.com/felipe/dev-test-api/internal/users"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Question struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	User       *users.User    `gorm:"foreignKey:UserID" json:"-"`
	Content    string         `gorm:"not null" json:"content"`
	Answer     string         `json:"answer,omitempty"`
	Topic      string         `json:"topic,omitempty"`
	ReviewedAt *time.Time     `json:"reviewed_at,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

type Store interface {
	Create(question *Question) error
	FindByUserID(userID uuid.UUID) ([]Question, error)
	FindByID(id uuid.UUID) (*Question, error)
	Update(question *Question) error
	SoftDelete(id uuid.UUID) error
}

type Service interface {
	// To be implemented when AI integration is ready
}

type Handler struct {
	// To be implemented
}
