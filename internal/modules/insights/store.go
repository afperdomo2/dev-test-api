package insights

import (
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TopicAggregateRow es una proyección de agregación, NO un DTO de API.
type TopicAggregateRow struct {
	TopicID      uuid.UUID `gorm:"column:topic_id"`
	TopicName    string    `gorm:"column:topic_name"`
	Attempts     int       `gorm:"column:attempts"`
	CorrectCount int       `gorm:"column:correct_count"`
}

type Store interface {
	CountSessions(userID uuid.UUID, from time.Time) (int, error)
	CountAnswers(userID uuid.UUID, from time.Time) (total int, correct int, err error)
	SumResponseTimeMs(userID uuid.UUID, from time.Time) (int64, error)
	TopTopics(userID uuid.UUID, from time.Time, limit int) ([]TopicAggregateRow, error)
	DailyActivity(userID uuid.UUID, from time.Time) (map[string]int, error)
	CountCorrectByDifficulty(userID uuid.UUID, from time.Time) (map[string]int, error)
}

type gormStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) CountSessions(userID uuid.UUID, from time.Time) (int, error) {
	var n int64
	err := s.db.Model(&models.Session{}).
		Where("user_id = ? AND status = ? AND created_at >= ?", userID, models.SessionStatusCompleted, from).
		Count(&n).Error
	return int(n), err
}

func (s *gormStore) CountAnswers(userID uuid.UUID, from time.Time) (int, int, error) {
	var row struct {
		Total   int64 `gorm:"column:total"`
		Correct int64 `gorm:"column:correct"`
	}
	err := s.db.
		Table("session_answers AS sa").
		Select("COUNT(*) AS total, COUNT(*) FILTER (WHERE sa.is_correct) AS correct").
		Joins("JOIN sessions AS s ON s.id = sa.session_id AND s.status = ?", models.SessionStatusCompleted).
		Where("sa.user_id = ? AND sa.created_at >= ?", userID, from).
		Scan(&row).Error
	return int(row.Total), int(row.Correct), err
}

func (s *gormStore) SumResponseTimeMs(userID uuid.UUID, from time.Time) (int64, error) {
	var row struct {
		Total int64 `gorm:"column:total"`
	}
	err := s.db.
		Table("session_answers AS sa").
		Select("COALESCE(SUM(sa.response_time_ms), 0) AS total").
		Joins("JOIN sessions AS s ON s.id = sa.session_id AND s.status = ?", models.SessionStatusCompleted).
		Where("sa.user_id = ? AND sa.created_at >= ?", userID, from).
		Scan(&row).Error
	return row.Total, err
}

func (s *gormStore) TopTopics(userID uuid.UUID, from time.Time, limit int) ([]TopicAggregateRow, error) {
	var rows []TopicAggregateRow
	err := s.db.
		Table("session_answers AS sa").
		Select(`t.id AS topic_id, t.name AS topic_name,
		        COUNT(sa.id) AS attempts,
		        COUNT(sa.id) FILTER (WHERE sa.is_correct) AS correct_count`).
		Joins("JOIN sessions AS s ON s.id = sa.session_id AND s.status = ?", models.SessionStatusCompleted).
		Joins("JOIN question_topics AS qt ON qt.question_id = sa.question_id").
		Joins("JOIN topics AS t ON t.id = qt.topic_id").
		Where("sa.user_id = ? AND sa.created_at >= ?", userID, from).
		Group("t.id, t.name").
		Order("attempts DESC, (COUNT(*) FILTER (WHERE sa.is_correct))::numeric / COUNT(*) DESC, t.name ASC").
		Limit(limit).
		Scan(&rows).Error
	return rows, err
}

func (s *gormStore) DailyActivity(userID uuid.UUID, from time.Time) (map[string]int, error) {
	var rows []struct {
		Day string `gorm:"column:day"`
		N   int    `gorm:"column:n"`
	}
	err := s.db.
		Table("session_answers AS sa").
		Select("to_char(sa.created_at AT TIME ZONE 'UTC', 'YYYY-MM-DD') AS day, COUNT(*) AS n").
		Joins("JOIN sessions AS s ON s.id = sa.session_id AND s.status = ?", models.SessionStatusCompleted).
		Where("sa.user_id = ? AND sa.created_at >= ?", userID, from).
		Group("day").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[string]int, len(rows))
	for _, r := range rows {
		out[r.Day] = r.N
	}
	return out, nil
}

func (s *gormStore) CountCorrectByDifficulty(userID uuid.UUID, from time.Time) (map[string]int, error) {
	var rows []struct {
		Difficulty string `gorm:"column:difficulty"`
		N          int    `gorm:"column:n"`
	}
	err := s.db.
		Table("session_answers AS sa").
		Select("q.difficulty AS difficulty, COUNT(*) AS n").
		Joins("JOIN sessions AS s ON s.id = sa.session_id AND s.status = ?", models.SessionStatusCompleted).
		Joins("JOIN questions AS q ON q.id = sa.question_id").
		Where("sa.user_id = ? AND sa.created_at >= ? AND sa.is_correct", userID, from).
		Group("q.difficulty").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[string]int, len(rows))
	for _, r := range rows {
		out[r.Difficulty] = r.N
	}
	return out, nil
}
