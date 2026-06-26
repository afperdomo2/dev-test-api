package sessions

import (
	"time"

	"github.com/felipe/dev-test-api/internal/modules/questions"
	"github.com/google/uuid"
)

type SessionResponse struct {
	ID            uuid.UUID   `json:"id"`
	UserID        uuid.UUID   `json:"userId"`
	Name          string      `json:"name"`
	Status        string      `json:"status"`
	Mode          string      `json:"mode"`
	Difficulty    string      `json:"difficulty"`
	QuestionLimit *int        `json:"questionLimit,omitempty"`
	Score         *float64    `json:"score,omitempty"`
	StartedAt     time.Time   `json:"startedAt"`
	FinishedAt    *time.Time  `json:"finishedAt,omitempty"`
	Topics        []TopicInfo `json:"topics"`
	AnswerCount   int         `json:"answerCount"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt"`
}

type SessionListResponse struct {
	ID            uuid.UUID   `json:"id"`
	UserID        uuid.UUID   `json:"userId"`
	Name          string      `json:"name"`
	Status        string      `json:"status"`
	Mode          string      `json:"mode"`
	Difficulty    string      `json:"difficulty"`
	QuestionLimit *int        `json:"questionLimit,omitempty"`
	Score         *float64    `json:"score,omitempty"`
	StartedAt     time.Time   `json:"startedAt"`
	FinishedAt    *time.Time  `json:"finishedAt,omitempty"`
	Topics        []TopicInfo `json:"topics"`
	AnswerCount   int         `json:"answerCount"`
	CreatedAt     time.Time   `json:"createdAt"`
}

type SessionDetailResponse struct {
	Session SessionResponse         `json:"session"`
	Answers []SessionAnswerResponse `json:"answers"`
}

type SessionAnswerResponse struct {
	ID              uuid.UUID                  `json:"id"`
	QuestionID      uuid.UUID                  `json:"questionId"`
	AnswerText      string                     `json:"answerText,omitempty"`
	SelectedOptions []uuid.UUID                `json:"selectedOptions,omitempty"`
	IsCorrect       bool                       `json:"isCorrect"`
	AiFeedback      string                     `json:"aiFeedback,omitempty"`
	ResponseTimeMs  int                        `json:"responseTimeMs"`
	Question        questions.QuestionResponse `json:"question"`
	CreatedAt       time.Time                  `json:"createdAt"`
}

type NextQuestionResponse struct {
	Question questions.QuestionResponse `json:"question"`
}

type TopicInfo struct {
	ID   uuid.UUID `json:"id"`
	Slug string    `json:"slug"`
	Name string    `json:"name"`
}
