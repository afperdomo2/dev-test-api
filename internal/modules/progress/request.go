package progress

import "github.com/felipe/dev-test-api/internal/common"

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
