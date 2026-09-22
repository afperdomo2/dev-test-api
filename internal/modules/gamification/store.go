package gamification

import (
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store interface {
	GetUserDailyGoal(userID uuid.UUID) (int, error)
	SetUserDailyGoal(userID uuid.UUID, dailyGoal int) error
	GetStreak(userID uuid.UUID) (*models.UserStreak, error)
	UpsertStreak(streak *models.UserStreak) error
	GetXp(userID uuid.UUID) (*models.UserXp, error)
	UpsertXp(xp *models.UserXp) error
	GetTopicXp(userID, topicID uuid.UUID) (*models.UserTopicXp, error)
	UpsertTopicXp(xp *models.UserTopicXp) error
	GetQuestionMeta(questionID uuid.UUID) (difficulty string, topicIDs []uuid.UUID, err error)
	GetTopicXpOverview(userID uuid.UUID) ([]TopicXpItem, error)
}

type gormStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) GetUserDailyGoal(userID uuid.UUID) (int, error) {
	var user models.User
	err := s.db.Select("daily_goal_xp").First(&user, "id = ?", userID).Error
	if err != nil {
		return 0, err
	}
	return user.DailyGoalXp, nil
}

func (s *gormStore) SetUserDailyGoal(userID uuid.UUID, dailyGoal int) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Update("daily_goal_xp", dailyGoal).Error
}

func (s *gormStore) GetStreak(userID uuid.UUID) (*models.UserStreak, error) {
	var streak models.UserStreak
	err := s.db.First(&streak, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &streak, nil
}

func (s *gormStore) UpsertStreak(streak *models.UserStreak) error {
	return s.db.Save(streak).Error
}

func (s *gormStore) GetXp(userID uuid.UUID) (*models.UserXp, error) {
	var xp models.UserXp
	err := s.db.First(&xp, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &xp, nil
}

func (s *gormStore) UpsertXp(xp *models.UserXp) error {
	return s.db.Save(xp).Error
}

func (s *gormStore) GetTopicXp(userID, topicID uuid.UUID) (*models.UserTopicXp, error) {
	var xp models.UserTopicXp
	err := s.db.First(&xp, "user_id = ? AND topic_id = ?", userID, topicID).Error
	if err != nil {
		return nil, err
	}
	return &xp, nil
}

func (s *gormStore) UpsertTopicXp(xp *models.UserTopicXp) error {
	return s.db.Save(xp).Error
}

func (s *gormStore) GetQuestionMeta(questionID uuid.UUID) (string, []uuid.UUID, error) {
	var question models.Question
	err := s.db.Preload("QuestionTopics").First(&question, "id = ?", questionID).Error
	if err != nil {
		return "", nil, err
	}

	topicIDs := make([]uuid.UUID, 0, len(question.QuestionTopics))
	for _, qt := range question.QuestionTopics {
		topicIDs = append(topicIDs, qt.TopicID)
	}
	return question.Difficulty, topicIDs, nil
}

func (s *gormStore) GetTopicXpOverview(userID uuid.UUID) ([]TopicXpItem, error) {
	var rows []TopicXpItem
	err := s.db.Table("user_topic_xp").
		Select("user_topic_xp.topic_id AS topic_id, COALESCE(topics.name, '') AS name, user_topic_xp.xp AS xp").
		Joins("LEFT JOIN topics ON topics.id = user_topic_xp.topic_id").
		Where("user_topic_xp.user_id = ?", userID).
		Order("user_topic_xp.xp DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []TopicXpItem{}
	}
	return rows, nil
}
