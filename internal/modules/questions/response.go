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
	Topics        []string               `json:"topics"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
}

type QuestionListResponse struct {
	ID            uuid.UUID              `json:"id"`
	UserID        uuid.UUID              `json:"userId"`
	Type          string                 `json:"type"`
	Content       string                 `json:"content"`
	Difficulty    string                 `json:"difficulty"`
	IsPublic      bool                   `json:"isPublic"`
	Source        string                 `json:"source"`
	CodeChallenge *CodeChallengeResponse `json:"codeChallenge,omitempty"`
	Topics        []string               `json:"topics"`
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

type ImportRowError struct {
	Row    int    `json:"row"`
	Reason string `json:"reason"`
}

type ImportResult struct {
	Imported int              `json:"imported"`
	Failed   int              `json:"failed"`
	Total    int              `json:"total"`
	Errors   []ImportRowError `json:"errors"`
}

type ImportQuota struct {
	DailyLimit int `json:"dailyLimit"`
	UsedToday  int `json:"usedToday"`
	Remaining  int `json:"remaining"`
}

type AiQuota struct {
	DailyLimit int `json:"dailyLimit"`
	UsedToday  int `json:"usedToday"`
	Remaining  int `json:"remaining"`
}

type TopicCount struct {
	TopicID  uuid.UUID `json:"topicId"`
	Name     string    `json:"name"`
	Slug     string    `json:"slug"`
	Category string    `json:"category"`
	Count    int64     `json:"count"`
}

type CategoryCount struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

type DifficultyCount struct {
	Difficulty string `json:"difficulty"`
	Count      int64  `json:"count"`
}

type TypeCount struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

type QuestionStats struct {
	Total        int64             `json:"total"`
	ByTopic      []TopicCount      `json:"byTopic"`
	ByCategory   []CategoryCount   `json:"byCategory"`
	ByDifficulty []DifficultyCount `json:"byDifficulty"`
	ByType       []TypeCount       `json:"byType"`
}
