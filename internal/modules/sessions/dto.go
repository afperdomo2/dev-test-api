package sessions

import (
	"encoding/json"
	"time"

	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/modules/questions"
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
	TopicIDs      []uuid.UUID `json:"topic_ids" binding:"required,min=1"`
	QuestionLimit *int        `json:"question_limit" binding:"omitempty,min=1,max=50"`
}

type AnswerRequest struct {
	QuestionID      uuid.UUID   `json:"question_id" binding:"required"`
	IsCorrect       bool        `json:"is_correct"`
	AnswerText      string      `json:"answer_text,omitempty"`
	SelectedOptions []uuid.UUID `json:"selected_options,omitempty"`
	ResponseTimeMs  int         `json:"response_time_ms"`
}

type SessionResponse struct {
	ID            uuid.UUID   `json:"id"`
	UserID        uuid.UUID   `json:"user_id"`
	Name          string      `json:"name"`
	Status        string      `json:"status"`
	Mode          string      `json:"mode"`
	Difficulty    string      `json:"difficulty"`
	QuestionLimit *int        `json:"question_limit,omitempty"`
	Score         *float64    `json:"score,omitempty"`
	StartedAt     time.Time   `json:"started_at"`
	FinishedAt    *time.Time  `json:"finished_at,omitempty"`
	Topics        []TopicInfo `json:"topics"`
	AnswerCount   int         `json:"answer_count"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type SessionListResponse struct {
	ID            uuid.UUID   `json:"id"`
	UserID        uuid.UUID   `json:"user_id"`
	Name          string      `json:"name"`
	Status        string      `json:"status"`
	Mode          string      `json:"mode"`
	Difficulty    string      `json:"difficulty"`
	QuestionLimit *int        `json:"question_limit,omitempty"`
	Score         *float64    `json:"score,omitempty"`
	StartedAt     time.Time   `json:"started_at"`
	FinishedAt    *time.Time  `json:"finished_at,omitempty"`
	Topics        []TopicInfo `json:"topics"`
	AnswerCount   int         `json:"answer_count"`
	CreatedAt     time.Time   `json:"created_at"`
}

type SessionDetailResponse struct {
	Session SessionResponse         `json:"session"`
	Answers []SessionAnswerResponse `json:"answers"`
}

type SessionAnswerResponse struct {
	ID              uuid.UUID                  `json:"id"`
	QuestionID      uuid.UUID                  `json:"question_id"`
	AnswerText      string                     `json:"answer_text,omitempty"`
	SelectedOptions []uuid.UUID                `json:"selected_options,omitempty"`
	IsCorrect       bool                       `json:"is_correct"`
	AiFeedback      string                     `json:"ai_feedback,omitempty"`
	ResponseTimeMs  int                        `json:"response_time_ms"`
	Question        questions.QuestionResponse `json:"question"`
	CreatedAt       time.Time                  `json:"created_at"`
}

type NextQuestionResponse struct {
	Question questions.QuestionResponse `json:"question"`
}

type TopicInfo struct {
	ID   uuid.UUID `json:"id"`
	Slug string    `json:"slug"`
	Name string    `json:"name"`
}

func toSessionResponse(s models.Session) SessionResponse {
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

func toSessionListResponse(s models.Session) SessionListResponse {
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

func toAnswerResponse(a models.SessionAnswer) SessionAnswerResponse {
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
