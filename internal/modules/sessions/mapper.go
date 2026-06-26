package sessions

import (
	"encoding/json"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/modules/questions"
	"github.com/google/uuid"
)

func ToSessionResponse(s models.Session) SessionResponse {
	topics := make([]TopicInfo, len(s.Topics))
	for i, t := range s.Topics {
		topics[i] = TopicInfo{ID: t.ID, Slug: t.Slug, Name: t.Name}
	}
	return SessionResponse{
		ID:            s.ID,
		UserID:        s.UserID,
		Name:          s.Name,
		Status:        s.Status,
		Mode:          s.Mode,
		Difficulty:    s.Difficulty,
		QuestionLimit: s.QuestionLimit,
		Score:         s.Score,
		StartedAt:     s.StartedAt,
		FinishedAt:    s.FinishedAt,
		Topics:        topics,
		AnswerCount:   len(s.Answers),
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

func ToSessionListResponse(s models.Session) SessionListResponse {
	topics := make([]TopicInfo, len(s.Topics))
	for i, t := range s.Topics {
		topics[i] = TopicInfo{ID: t.ID, Slug: t.Slug, Name: t.Name}
	}
	return SessionListResponse{
		ID:            s.ID,
		UserID:        s.UserID,
		Name:          s.Name,
		Status:        s.Status,
		Mode:          s.Mode,
		Difficulty:    s.Difficulty,
		QuestionLimit: s.QuestionLimit,
		Score:         s.Score,
		StartedAt:     s.StartedAt,
		FinishedAt:    s.FinishedAt,
		Topics:        topics,
		AnswerCount:   len(s.Answers),
		CreatedAt:     s.CreatedAt,
	}
}

func ToAnswerResponse(a models.SessionAnswer) SessionAnswerResponse {
	var selected []uuid.UUID
	if a.SelectedOptions != "" {
		json.Unmarshal([]byte(a.SelectedOptions), &selected)
	}

	resp := SessionAnswerResponse{
		ID:              a.ID,
		QuestionID:      a.QuestionID,
		AnswerText:      a.AnswerText,
		SelectedOptions: selected,
		IsCorrect:       a.IsCorrect,
		AiFeedback:      a.AiFeedback,
		ResponseTimeMs:  a.ResponseTimeMs,
		CreatedAt:       a.CreatedAt,
	}

	if a.Question != nil {
		resp.Question = questions.ToQuestionResponse(*a.Question)
	}

	return resp
}
