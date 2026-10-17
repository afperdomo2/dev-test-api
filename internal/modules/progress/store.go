package progress

import (
	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TopicAccuracyRow es una proyección de agregación por tema, NO un DTO de API.
// La consume el service para calcular accuracy e isWeak.
type TopicAccuracyRow struct {
	TopicID      uuid.UUID `gorm:"column:topic_id"`
	TopicName    string    `gorm:"column:topic_name"`
	Attempts     int       `gorm:"column:attempts"`
	CorrectCount int       `gorm:"column:correct_count"`
}

type Store interface {
	FindByUserAndQuestion(userID, questionID uuid.UUID) (*models.UserQuestionProgress, error)
	Upsert(progress *models.UserQuestionProgress) error
	FindUpcoming(userID uuid.UUID, params common.PaginationParams) ([]models.UserQuestionProgress, int64, error)
	FindSaved(userID uuid.UUID, params common.PaginationParams) ([]models.UserQuestionProgress, int64, error)
	FindTopicAccuracy(userID uuid.UUID) ([]TopicAccuracyRow, error)
}

type gormStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) FindByUserAndQuestion(userID, questionID uuid.UUID) (*models.UserQuestionProgress, error) {
	var p models.UserQuestionProgress
	err := s.db.First(&p, "user_id = ? AND question_id = ?", userID, questionID).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *gormStore) Upsert(p *models.UserQuestionProgress) error {
	return s.db.Save(p).Error
}

func (s *gormStore) FindUpcoming(userID uuid.UUID, params common.PaginationParams) ([]models.UserQuestionProgress, int64, error) {
	var items []models.UserQuestionProgress
	var total int64

	base := s.db.Where("user_id = ? AND next_review_at IS NOT NULL AND next_review_at <= now() AND is_mastered = false", userID)
	base.Model(&models.UserQuestionProgress{}).Count(&total)

	err := base.Offset((params.Page - 1) * params.PerPage).Limit(params.PerPage).
		Preload("Question.Options").
		Preload("Question.CodeChallenge").
		Preload("Question.Topics").
		Order(upcomingSortConfig.OrderClause(params.SortBy, params.SortOrder)).
		Find(&items).Error
	return items, total, err
}

func (s *gormStore) FindSaved(userID uuid.UUID, params common.PaginationParams) ([]models.UserQuestionProgress, int64, error) {
	var items []models.UserQuestionProgress
	var total int64

	base := s.db.Where("user_id = ? AND is_saved = true", userID)
	base.Model(&models.UserQuestionProgress{}).Count(&total)

	err := base.Offset((params.Page - 1) * params.PerPage).Limit(params.PerPage).
		Preload("Question.Options").
		Preload("Question.CodeChallenge").
		Preload("Question.Topics").
		Order(savedSortConfig.OrderClause(params.SortBy, params.SortOrder)).
		Find(&items).Error
	return items, total, err
}

// FindTopicAccuracy agrega las respuestas del usuario por tema.
//
// Solo cuenta respuestas de sesiones completadas (D-10) y los temas sin
// respuestas no aparecen: el JOIN los excluye. Una pregunta asociada a varios
// temas cuenta para cada uno de ellos.
func (s *gormStore) FindTopicAccuracy(userID uuid.UUID) ([]TopicAccuracyRow, error) {
	var rows []TopicAccuracyRow
	err := s.db.
		Table("session_answers AS sa").
		Select(`t.id AS topic_id, t.name AS topic_name,
		        COUNT(sa.id) AS attempts,
		        COUNT(sa.id) FILTER (WHERE sa.is_correct) AS correct_count`).
		Joins("JOIN sessions AS s ON s.id = sa.session_id AND s.status = ?", models.SessionStatusCompleted).
		Joins("JOIN question_topics AS qt ON qt.question_id = sa.question_id").
		Joins("JOIN topics AS t ON t.id = qt.topic_id").
		Where("sa.user_id = ?", userID).
		Group("t.id, t.name").
		Order("t.name ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}
