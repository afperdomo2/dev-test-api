package generation

import (
	"github.com/felipe/dev-test-api/internal/common"
	"github.com/google/uuid"
)

var sortConfig = common.SortConfig{
	Allowed: []string{"created_at", "started_at", "finished_at"},
	Default: "created_at desc",
}

type ListJobsParams struct {
	common.PaginationParams
	Status string
	Kind   string
}

type CreateJobRequest struct {
	TopicIDs []uuid.UUID `json:"topicIds" binding:"required,min=1,dive,required"`
	// Kind elige qué se genera. Vacío = "question" (retrocompatible).
	Kind string `json:"kind,omitempty" binding:"omitempty,oneof=question tip"`
	// Difficulty solo aplica a preguntas; en tips se ignora.
	Difficulty string `json:"difficulty,omitempty" binding:"omitempty,oneof=beginner intermediate advanced random"`
	// QuestionType solo aplica a preguntas; en tips se ignora. "mixed" reparte
	// la cantidad entre selección única, múltiple y verdadero/falso.
	QuestionType string `json:"questionType,omitempty" binding:"omitempty,oneof=mixed single_choice multiple_choice true_false"`
	// Quantity es opcional: si se omite se usa el ajuste ai.default_quantity.
	Quantity int `json:"quantity,omitempty" binding:"omitempty,min=1"`
}
