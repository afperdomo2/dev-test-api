package questions

import (
	"time"

	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
)

var sortConfig = common.SortConfig{
	Allowed: []string{"type", "difficulty", "created_at", "updated_at"},
	Default: "created_at desc",
}

type CreateQuestionRequest struct {
	Type           string            `json:"type" binding:"required,oneof=single_choice multiple_choice code_completion"`
	Content        string            `json:"content" binding:"required"`
	Explanation    string            `json:"explanation,omitempty"`
	Difficulty     string            `json:"difficulty" binding:"required,oneof=beginner intermediate advanced"`
	Source         string            `json:"source,omitempty" binding:"omitempty,oneof=ai_generated manual imported"`
	TopicIDs       []uuid.UUID       `json:"topic_ids" binding:"required,min=1"`
	Options        []CreateOptionReq `json:"options,omitempty"`
	StarterCode    string            `json:"starter_code,omitempty"`
	ExpectedOutput string            `json:"expected_output,omitempty"`
	Language       string            `json:"language,omitempty"`
	TestCasesJSON  string            `json:"test_cases,omitempty"`
}

type CreateOptionReq struct {
	Content   string `json:"content" binding:"required"`
	IsCorrect bool   `json:"is_correct"`
}

type UpdateQuestionRequest struct {
	Content        string            `json:"content,omitempty" binding:"omitempty"`
	Explanation    string            `json:"explanation,omitempty"`
	Difficulty     string            `json:"difficulty,omitempty" binding:"omitempty,oneof=beginner intermediate advanced"`
	TopicIDs       []uuid.UUID       `json:"topic_ids,omitempty"`
	Options        []CreateOptionReq `json:"options,omitempty"`
	StarterCode    string            `json:"starter_code,omitempty"`
	ExpectedOutput string            `json:"expected_output,omitempty"`
	Language       string            `json:"language,omitempty"`
	TestCasesJSON  string            `json:"test_cases,omitempty"`
}

type QuestionResponse struct {
	ID            uuid.UUID              `json:"id"`
	UserID        uuid.UUID              `json:"user_id"`
	Type          string                 `json:"type"`
	Content       string                 `json:"content"`
	Explanation   string                 `json:"explanation,omitempty"`
	Difficulty    string                 `json:"difficulty"`
	IsPublic      bool                   `json:"is_public"`
	Source        string                 `json:"source"`
	Options       []OptionResponse       `json:"options,omitempty"`
	CodeChallenge *CodeChallengeResponse `json:"code_challenge,omitempty"`
	Topics        []TopicInfo            `json:"topics"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type OptionResponse struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	IsCorrect bool      `json:"is_correct"`
}

type CodeChallengeResponse struct {
	ID             uuid.UUID `json:"id"`
	StarterCode    string    `json:"starter_code,omitempty"`
	ExpectedOutput string    `json:"expected_output,omitempty"`
	Language       string    `json:"language"`
	TestCasesJSON  string    `json:"test_cases,omitempty"`
}

type TopicInfo struct {
	ID   uuid.UUID `json:"id"`
	Slug string    `json:"slug"`
	Name string    `json:"name"`
}

func ToQuestionResponse(q models.Question) QuestionResponse {
	resp := QuestionResponse{
		ID:          q.ID,
		UserID:      q.UserID,
		Type:        q.Type,
		Content:     q.Content,
		Explanation: q.Explanation,
		Difficulty:  q.Difficulty,
		IsPublic:    q.IsPublic,
		Source:      q.Source,
		CreatedAt:   q.CreatedAt,
		UpdatedAt:   q.UpdatedAt,
	}

	if q.Options != nil {
		resp.Options = make([]OptionResponse, len(q.Options))
		for i, o := range q.Options {
			resp.Options[i] = OptionResponse{
				ID:        o.ID,
				Content:   o.Content,
				IsCorrect: o.IsCorrect,
			}
		}
	}

	if q.CodeChallenge != nil {
		resp.CodeChallenge = &CodeChallengeResponse{
			ID:             q.CodeChallenge.ID,
			StarterCode:    q.CodeChallenge.StarterCode,
			ExpectedOutput: q.CodeChallenge.ExpectedOutput,
			Language:       q.CodeChallenge.Language,
			TestCasesJSON:  q.CodeChallenge.TestCasesJSON,
		}
	}

	if q.Topics != nil {
		resp.Topics = make([]TopicInfo, len(q.Topics))
		for i, t := range q.Topics {
			resp.Topics[i] = TopicInfo{
				ID:   t.ID,
				Slug: t.Slug,
				Name: t.Name,
			}
		}
	}

	return resp
}
