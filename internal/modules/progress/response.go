package progress

import (
	"time"

	"github.com/felipe/dev-test-api/internal/modules/questions"
	"github.com/google/uuid"
)

type ProgressResponse struct {
	QuestionID     uuid.UUID  `json:"questionId"`
	Repetitions    int        `json:"repetitions"`
	EaseFactor     float64    `json:"easeFactor"`
	IntervalDays   int        `json:"intervalDays"`
	NextReviewAt   *time.Time `json:"nextReviewAt,omitempty"`
	LastReviewedAt *time.Time `json:"lastReviewedAt,omitempty"`
	IsSaved        bool       `json:"isSaved"`
	IsMastered     bool       `json:"isMastered"`
}

type UpcomingItem struct {
	Question questions.QuestionResponse `json:"question"`
	Progress ProgressResponse           `json:"progress"`
}
