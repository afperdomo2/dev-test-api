package notifications

import (
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DueReviewUser struct {
	UserID uuid.UUID `gorm:"column:user_id"`
	N      int64     `gorm:"column:n"`
}

type Store interface {
	FindPage(userID uuid.UUID, params ListNotificationsParams) ([]models.Notification, int64, error)
	FindByIDAndUser(id, userID uuid.UUID) (*models.Notification, error)
	MarkRead(id, userID uuid.UUID) error
	MarkAllRead(userID uuid.UUID) (int64, error)
	CountUnread(userID uuid.UUID) (int64, error)
	// InsertBatch inserta ignorando conflictos y devuelve solo las filas nuevas
	// (las que el job debe entregar por push).
	InsertBatch(items []models.Notification) ([]models.Notification, error)
	DueReviewCounts() ([]DueReviewUser, error)
}

type gormStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) FindPage(userID uuid.UUID, params ListNotificationsParams) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := s.db.Model(&models.Notification{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset((params.Page - 1) * params.PerPage).
		Limit(params.PerPage).
		Find(&notifications).Error
	return notifications, total, err
}

func (s *gormStore) FindByIDAndUser(id, userID uuid.UUID) (*models.Notification, error) {
	var notification models.Notification
	err := s.db.
		Where("id = ? AND user_id = ?", id, userID).
		First(&notification).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (s *gormStore) MarkRead(id, userID uuid.UUID) error {
	return s.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func (s *gormStore) MarkAllRead(userID uuid.UUID) (int64, error) {
	result := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Update("is_read", true)
	return result.RowsAffected, result.Error
}

func (s *gormStore) CountUnread(userID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}

func (s *gormStore) InsertBatch(items []models.Notification) ([]models.Notification, error) {
	if len(items) == 0 {
		return []models.Notification{}, nil
	}

	// El índice único (user_id, type, creation_day) garantiza la idempotencia.
	// RETURNING devuelve solo las filas que realmente se insertaron.
	err := s.db.
		Clauses(clause.OnConflict{DoNothing: true}, clause.Returning{}).
		Create(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *gormStore) DueReviewCounts() ([]DueReviewUser, error) {
	var due []DueReviewUser
	// Misma definición de "repaso vencido" que sessions.CountDueReviews (D-09).
	err := s.db.Raw(`
		SELECT user_id, COUNT(*) AS n
		FROM user_question_progress
		WHERE next_review_at IS NOT NULL
		  AND next_review_at <= now()
		  AND is_mastered = false
		GROUP BY user_id
		HAVING COUNT(*) > 0
	`).Scan(&due).Error
	return due, err
}
