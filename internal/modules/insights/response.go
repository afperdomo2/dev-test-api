package insights

import (
	"time"

	"github.com/felipe/dev-test-api/internal/modules/progress"
	"github.com/google/uuid"
)

// TopicInsightItem es un tema con su actividad dentro del periodo.
type TopicInsightItem struct {
	TopicID   uuid.UUID `json:"topicId"`
	TopicName string    `json:"topicName"`
	Attempts  int       `json:"attempts"`
	Accuracy  float64   `json:"accuracy"`
}

// DailyActivityItem es un día de la ventana (siempre presente, aunque sea 0).
type DailyActivityItem struct {
	Date              string `json:"date"` // YYYY-MM-DD (UTC)
	QuestionsAnswered int    `json:"questionsAnswered"`
}

// InsightSummaryResponse es el informe completo de un periodo.
//
// `WeakTopics` reutiliza el DTO de `progress` a propósito: la definición de
// "tema débil" es única (D-08) y no se duplica aquí.
type InsightSummaryResponse struct {
	Period            string                      `json:"period"`
	From              time.Time                   `json:"from"`
	To                time.Time                   `json:"to"`
	Sessions          int                         `json:"sessions"`
	QuestionsAnswered int                         `json:"questionsAnswered"`
	Accuracy          float64                     `json:"accuracy"`
	TimeSpentMs       int64                       `json:"timeSpentMs"`
	TopTopics         []TopicInsightItem          `json:"topTopics"`
	WeakTopics        []progress.TopicMasteryItem `json:"weakTopics"`
	XpGained          int                         `json:"xpGained"`
	StreakDays        int                         `json:"streakDays"`
	DailyActivity     []DailyActivityItem         `json:"dailyActivity"`
}
