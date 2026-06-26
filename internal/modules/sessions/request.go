package sessions

import (
	"github.com/felipe/dev-test-api/internal/common"
	"github.com/google/uuid"
)

var sortConfig = common.SortConfig{
	Allowed: []string{"status", "score", "started_at", "created_at"},
	Default: "started_at DESC",
}

type CreateSessionRequest struct {
	Name          string      `json:"name" binding:"required,min=1,max=200"`
	Mode          string      `json:"mode" binding:"required,oneof=generate review"`
	Difficulty    string      `json:"difficulty" binding:"required,oneof=beginner intermediate advanced"`
	TopicIDs      []uuid.UUID `json:"topicIds" binding:"required,min=1"`
	QuestionLimit *int        `json:"questionLimit" binding:"omitempty,min=1,max=50"`
}

type AnswerRequest struct {
	QuestionID      uuid.UUID   `json:"questionId" binding:"required"`
	IsCorrect       bool        `json:"isCorrect"`
	AnswerText      string      `json:"answerText,omitempty"`
	SelectedOptions []uuid.UUID `json:"selectedOptions,omitempty"`
	ResponseTimeMs  int         `json:"responseTimeMs"`
}
