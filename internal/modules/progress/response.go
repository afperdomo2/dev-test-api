package progress

import (
	"time"

	"github.com/felipe/dev-test-api/internal/modules/gamification"
	"github.com/felipe/dev-test-api/internal/modules/questions"
	"github.com/google/uuid"
)

type ProgressResponse struct {
	QuestionID     uuid.UUID            `json:"questionId"`
	Repetitions    int                  `json:"repetitions"`
	EaseFactor     float64              `json:"easeFactor"`
	IntervalDays   int                  `json:"intervalDays"`
	NextReviewAt   *time.Time           `json:"nextReviewAt,omitempty"`
	LastReviewedAt *time.Time           `json:"lastReviewedAt,omitempty"`
	IsSaved        bool                 `json:"isSaved"`
	IsMastered     bool                 `json:"isMastered"`
	Reward         *gamification.Reward `json:"reward,omitempty"`
}

type UpcomingItem struct {
	Question questions.QuestionResponse `json:"question"`
	Progress ProgressResponse           `json:"progress"`
}

// TopicMasteryItem es el dominio del usuario en un tema. `isWeak` se calcula
// con la definición única de D-08 (accuracy < 60% y >= 5 respuestas).
type TopicMasteryItem struct {
	TopicID      uuid.UUID `json:"topicId"`
	TopicName    string    `json:"topicName"`
	Attempts     int       `json:"attempts"`
	CorrectCount int       `json:"correctCount"`
	Accuracy     float64   `json:"accuracy"`
	IsWeak       bool      `json:"isWeak"`
}
