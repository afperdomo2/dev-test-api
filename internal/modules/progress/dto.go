package progress

import (
	"time"

	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/modules/questions"
	"github.com/google/uuid"
)

var upcomingSortConfig = common.SortConfig{
	Allowed: []string{"next_review_at", "repetitions", "ease_factor"},
	Default: "next_review_at ASC",
}

var savedSortConfig = common.SortConfig{
	Allowed: []string{"updated_at", "repetitions", "ease_factor"},
	Default: "updated_at DESC",
}

type AnswerRequest struct {
	IsCorrect bool `json:"isCorrect" binding:"required"`
}

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
