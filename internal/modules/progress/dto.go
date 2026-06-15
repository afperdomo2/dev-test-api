package progress

import (
	"time"

	"github.com/felipe/dev-test-api/internal/modules/questions"
	"github.com/google/uuid"
)

type AnswerRequest struct {
	IsCorrect bool `json:"is_correct" binding:"required"`
}

type ProgressResponse struct {
	QuestionID     uuid.UUID  `json:"question_id"`
	Repetitions    int        `json:"repetitions"`
	EaseFactor     float64    `json:"ease_factor"`
	IntervalDays   int        `json:"interval_days"`
	NextReviewAt   *time.Time `json:"next_review_at,omitempty"`
	LastReviewedAt *time.Time `json:"last_reviewed_at,omitempty"`
	IsSaved        bool       `json:"is_saved"`
	IsMastered     bool       `json:"is_mastered"`
}

type UpcomingItem struct {
	Question questions.QuestionResponse `json:"question"`
	Progress ProgressResponse           `json:"progress"`
}
