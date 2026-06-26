package questions

import (
	"github.com/felipe/dev-test-api/internal/common"
	"github.com/google/uuid"
)

var sortConfig = common.SortConfig{
	Allowed: []string{"type", "difficulty", "created_at", "updated_at"},
	Default: "created_at desc",
}

type ListQuestionsParams struct {
	common.PaginationParams
	Type       string
	Difficulty string
	TopicIDs   []uuid.UUID
}

type CreateQuestionRequest struct {
	Type           string            `json:"type" binding:"required,oneof=single_choice multiple_choice code_completion"`
	Content        string            `json:"content" binding:"required"`
	Explanation    string            `json:"explanation,omitempty"`
	Difficulty     string            `json:"difficulty" binding:"required,oneof=beginner intermediate advanced"`
	Source         string            `json:"source,omitempty" binding:"omitempty,oneof=ai_generated manual imported"`
	TopicIDs       []uuid.UUID       `json:"topicIds" binding:"required,min=1"`
	Options        []CreateOptionReq `json:"options,omitempty"`
	StarterCode    string            `json:"starterCode,omitempty"`
	ExpectedOutput string            `json:"expectedOutput,omitempty"`
	Language       string            `json:"language,omitempty"`
	TestCasesJSON  string            `json:"testCases,omitempty"`
}

type CreateOptionReq struct {
	Content   string `json:"content" binding:"required"`
	IsCorrect bool   `json:"isCorrect"`
}

type UpdateQuestionRequest struct {
	Content        string            `json:"content,omitempty" binding:"omitempty"`
	Explanation    string            `json:"explanation,omitempty"`
	Difficulty     string            `json:"difficulty,omitempty" binding:"omitempty,oneof=beginner intermediate advanced"`
	TopicIDs       []uuid.UUID       `json:"topicIds,omitempty"`
	Options        []CreateOptionReq `json:"options,omitempty"`
	StarterCode    string            `json:"starterCode,omitempty"`
	ExpectedOutput string            `json:"expectedOutput,omitempty"`
	Language       string            `json:"language,omitempty"`
	TestCasesJSON  string            `json:"testCases,omitempty"`
}
