package sessions

import (
	"time"

	"github.com/felipe/dev-test-api/internal/modules/questions"
	"github.com/google/uuid"
)

type SessionResponse struct {
	ID                 uuid.UUID  `json:"id"`
	UserID             uuid.UUID  `json:"userId"`
	Name               string     `json:"name"`
	Status             string     `json:"status"`
	Mode               string     `json:"mode"`
	Difficulty         string     `json:"difficulty"`
	QuestionLimit      *int       `json:"questionLimit,omitempty"`
	QuestionsGenerated int        `json:"questionsGenerated"`
	Score              *float64   `json:"score,omitempty"`
	StartedAt          time.Time  `json:"startedAt"`
	FinishedAt         *time.Time `json:"finishedAt,omitempty"`
	Topics             []string   `json:"topics"`
	AnswerCount        int        `json:"answerCount"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`
}

type SessionListResponse struct {
	ID                 uuid.UUID  `json:"id"`
	UserID             uuid.UUID  `json:"userId"`
	Name               string     `json:"name"`
	Status             string     `json:"status"`
	Mode               string     `json:"mode"`
	Difficulty         string     `json:"difficulty"`
	QuestionLimit      *int       `json:"questionLimit,omitempty"`
	QuestionsGenerated int        `json:"questionsGenerated"`
	Score              *float64   `json:"score,omitempty"`
	StartedAt          time.Time  `json:"startedAt"`
	FinishedAt         *time.Time `json:"finishedAt,omitempty"`
	Topics             []string   `json:"topics"`
	AnswerCount        int        `json:"answerCount"`
	CreatedAt          time.Time  `json:"createdAt"`
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

type NextQuestionOption struct {
	ID      uuid.UUID `json:"id"`
	Content string    `json:"content"`
}

type NextQuestionItem struct {
	ID            uuid.UUID                        `json:"id"`
	Type          string                           `json:"type"`
	Content       string                           `json:"content"`
	Difficulty    string                           `json:"difficulty"`
	Options       []NextQuestionOption             `json:"options,omitempty"`
	CodeChallenge *questions.CodeChallengeResponse `json:"codeChallenge,omitempty"`
	Topics        []string                         `json:"topics"`
}

type NextQuestionResponse struct {
	Question NextQuestionItem `json:"question"`
}

type SessionSummaryResponse struct {
	AnswerCount        int    `json:"answerCount"`
	QuestionsGenerated int    `json:"questionsGenerated"`
	Status             string `json:"status"`
	QuestionLimit      *int   `json:"questionLimit,omitempty"`
}
