package progress

import (
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store interface {
	FindByUserAndQuestion(userID, questionID uuid.UUID) (*models.UserQuestionProgress, error)
	Upsert(progress *models.UserQuestionProgress) error
	FindUpcoming(userID uuid.UUID, page, perPage int, sortBy, sortOrder string) ([]models.UserQuestionProgress, int64, error)
	FindSaved(userID uuid.UUID, page, perPage int, sortBy, sortOrder string) ([]models.UserQuestionProgress, int64, error)
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

func (s *gormStore) FindUpcoming(userID uuid.UUID, page, perPage int, sortBy, sortOrder string) ([]models.UserQuestionProgress, int64, error) {
	var items []models.UserQuestionProgress
	var total int64

	base := s.db.Where("user_id = ? AND is_saved = true AND next_review_at <= now()", userID)
	base.Model(&models.UserQuestionProgress{}).Count(&total)

	err := base.Offset((page - 1) * perPage).Limit(perPage).
		Preload("Question.Options").
		Preload("Question.CodeChallenge").
		Preload("Question.Topics").
		Order(upcomingSortConfig.OrderClause(sortBy, sortOrder)).
		Find(&items).Error
	return items, total, err
}

func (s *gormStore) FindSaved(userID uuid.UUID, page, perPage int, sortBy, sortOrder string) ([]models.UserQuestionProgress, int64, error) {
	var items []models.UserQuestionProgress
	var total int64

	base := s.db.Where("user_id = ? AND is_saved = true", userID)
	base.Model(&models.UserQuestionProgress{}).Count(&total)

	err := base.Offset((page - 1) * perPage).Limit(perPage).
		Preload("Question.Options").
		Preload("Question.CodeChallenge").
		Preload("Question.Topics").
		Order(savedSortConfig.OrderClause(sortBy, sortOrder)).
		Find(&items).Error
	return items, total, err
}
