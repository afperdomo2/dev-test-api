package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	GenerationJobStatusPending   = "pending"
	GenerationJobStatusRunning   = "running"
	GenerationJobStatusCompleted = "completed"
	GenerationJobStatusFailed    = "failed"
)

// GenerationJob es una solicitud de generación de preguntas con IA en lote.
// Solo puede crearla un administrador; la ejecuta un worker del pool.
type GenerationJob struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	RequestedBy uuid.UUID  `gorm:"type:uuid;not null;index" json:"requestedBy"`
	Requester   *User      `gorm:"foreignKey:RequestedBy" json:"-"`
	Difficulty  string     `gorm:"not null;default:intermediate;size:20" json:"difficulty"`
	Quantity    int        `gorm:"not null" json:"quantity"`
	Status      string     `gorm:"not null;default:pending;size:20;index" json:"status"`
	Generated   int        `gorm:"not null;default:0" json:"generated"`
	Error       string     `gorm:"type:text" json:"error,omitempty"`
	StartedAt   *time.Time `json:"startedAt,omitempty"`
	FinishedAt  *time.Time `json:"finishedAt,omitempty"`
	CreatedAt   time.Time  `gorm:"index" json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`

	Topics    []Topic              `gorm:"many2many:generation_job_topics;" json:"-"`
	JobTopics []GenerationJobTopic `gorm:"foreignKey:JobID" json:"-"`
}

func (GenerationJob) TableName() string { return "generation_jobs" }

func (j *GenerationJob) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}
