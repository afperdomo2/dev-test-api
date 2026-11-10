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

// Kind diferencia qué tipo de contenido genera el trabajo. Ambos tipos
// comparten la misma cola, worker y presupuesto diario (ver D-51).
const (
	GenerationKindQuestion = "question"
	GenerationKindTip      = "tip"
)

// Tipos de pregunta que puede pedir un trabajo de generación (D-54). "mixed"
// reparte la cantidad entre los tres tipos corregibles. code_completion queda
// fuera: ningún cliente lo corrige automáticamente todavía.
const (
	GenerationQuestionTypeMixed          = "mixed"
	GenerationQuestionTypeSingleChoice   = QuestionTypeSingleChoice
	GenerationQuestionTypeMultipleChoice = QuestionTypeMultipleChoice
	GenerationQuestionTypeTrueFalse      = QuestionTypeTrueFalse
)

// GenerationJob es una solicitud de generación de contenido con IA en lote.
// Solo puede crearla un administrador; la ejecuta un worker del pool.
type GenerationJob struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	RequestedBy  uuid.UUID  `gorm:"type:uuid;not null;index" json:"requestedBy"`
	Requester    *User      `gorm:"foreignKey:RequestedBy" json:"-"`
	Kind         string     `gorm:"not null;default:question;size:20;index" json:"kind"`
	Difficulty   string     `gorm:"not null;default:intermediate;size:20" json:"difficulty"`
	QuestionType string     `gorm:"not null;default:single_choice;size:20" json:"questionType"`
	Quantity     int        `gorm:"not null" json:"quantity"`
	Status       string     `gorm:"not null;default:pending;size:20;index" json:"status"`
	Generated    int        `gorm:"not null;default:0" json:"generated"`
	Error        string     `gorm:"type:text" json:"error,omitempty"`
	StartedAt    *time.Time `json:"startedAt,omitempty"`
	FinishedAt   *time.Time `json:"finishedAt,omitempty"`
	CreatedAt    time.Time  `gorm:"index" json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`

	Topics    []Topic              `gorm:"many2many:generation_job_topics;joinForeignKey:JobID;joinReferences:TopicID" json:"-"`
	JobTopics []GenerationJobTopic `gorm:"foreignKey:JobID" json:"-"`
}

func (GenerationJob) TableName() string { return "generation_jobs" }

func (j *GenerationJob) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}
