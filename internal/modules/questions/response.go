package questions

import (
	"time"

	"github.com/google/uuid"
)

type QuestionResponse struct {
	ID            uuid.UUID              `json:"id"`
	UserID        uuid.UUID              `json:"userId"`
	Type          string                 `json:"type"`
	Content       string                 `json:"content"`
	Explanation   string                 `json:"explanation,omitempty"`
	Difficulty    string                 `json:"difficulty"`
	IsPublic      bool                   `json:"isPublic"`
	Source        string                 `json:"source"`
	Options       []OptionResponse       `json:"options,omitempty"`
	CodeChallenge *CodeChallengeResponse `json:"codeChallenge,omitempty"`
	Topics        []TopicInfo            `json:"topics"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
}

type QuestionListResponse struct {
	ID            uuid.UUID              `json:"id"`
	UserID        uuid.UUID              `json:"userId"`
	Type          string                 `json:"type"`
	Content       string                 `json:"content"`
	Explanation   string                 `json:"explanation,omitempty"`
	Difficulty    string                 `json:"difficulty"`
	IsPublic      bool                   `json:"isPublic"`
	Source        string                 `json:"source"`
	Options       []OptionResponse       `json:"options,omitempty"`
	CodeChallenge *CodeChallengeResponse `json:"codeChallenge,omitempty"`
	Topics        []TopicInfo            `json:"topics"`
	CreatedAt     time.Time              `json:"createdAt"`
}

type OptionResponse struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	IsCorrect bool      `json:"isCorrect"`
}

type CodeChallengeResponse struct {
	ID             uuid.UUID `json:"id"`
	StarterCode    string    `json:"starterCode,omitempty"`
	ExpectedOutput string    `json:"expectedOutput,omitempty"`
	Language       string    `json:"language"`
	TestCasesJSON  string    `json:"testCases,omitempty"`
}

type TopicInfo struct {
	ID   uuid.UUID `json:"id"`
	Slug string    `json:"slug"`
	Name string    `json:"name"`
}
