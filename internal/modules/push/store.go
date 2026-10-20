package push

import (
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Store interface {
	// Upsert registra el token; si ya existía (otro usuario u otro platform) lo
	// reasigna y refresca last_seen_at. Nunca duplica.
	Upsert(userID uuid.UUID, token, platform string) error
	// DeleteByTokenAndUser borra solo si el token pertenece al usuario.
	DeleteByTokenAndUser(userID uuid.UUID, token string) (bool, error)
	FindByUser(userID uuid.UUID) ([]models.DeviceToken, error)
	// DeleteByToken se usa cuando Expo responde DeviceNotRegistered (global).
	DeleteByToken(token string) error
}

type gormStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) Upsert(userID uuid.UUID, token, platform string) error {
	now := time.Now()
	row := models.DeviceToken{
		UserID:     userID,
		Token:      token,
		Platform:   platform,
		LastSeenAt: now,
	}

	return s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "token"}},
		DoUpdates: clause.Assignments(map[string]any{
			"user_id":      userID,
			"platform":     platform,
			"last_seen_at": now,
			"updated_at":   now,
		}),
	}).Create(&row).Error
}

func (s *gormStore) DeleteByTokenAndUser(userID uuid.UUID, token string) (bool, error) {
	result := s.db.
		Where("user_id = ? AND token = ?", userID, token).
		Delete(&models.DeviceToken{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (s *gormStore) FindByUser(userID uuid.UUID) ([]models.DeviceToken, error) {
	var tokens []models.DeviceToken
	err := s.db.Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}

func (s *gormStore) DeleteByToken(token string) error {
	return s.db.Where("token = ?", token).Delete(&models.DeviceToken{}).Error
}
