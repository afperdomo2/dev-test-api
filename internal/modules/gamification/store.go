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

	// Agregaciones para evaluar logros (D-16: sin tabla).
	CountCompletedSessions(userID uuid.UUID) (int, error)
	CountAnswers(userID uuid.UUID) (int, error)
	HasPerfectSession(userID uuid.UUID) (bool, error)
	HasFastSession(userID uuid.UUID) (bool, error)
	CountMasteredTopics(userID uuid.UUID) (int, error)
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
	err := s.db.Table("user_topic_xps").
		Select("user_topic_xps.topic_id AS topic_id, COALESCE(topics.name, '') AS name, user_topic_xps.xp AS xp").
		Joins("LEFT JOIN topics ON topics.id = user_topic_xps.topic_id").
		Where("user_topic_xps.user_id = ?", userID).
		Order("user_topic_xps.xp DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []TopicXpItem{}
	}
	return rows, nil
}

func (s *gormStore) CountCompletedSessions(userID uuid.UUID) (int, error) {
	var n int64
	err := s.db.Model(&models.Session{}).
		Where("user_id = ? AND status = ?", userID, models.SessionStatusCompleted).
		Count(&n).Error
	return int(n), err
}

func (s *gormStore) CountAnswers(userID uuid.UUID) (int, error) {
	var n int64
	err := s.db.Model(&models.SessionAnswer{}).
		Where("user_id = ?", userID).
		Count(&n).Error
	return int(n), err
}

// HasPerfectSession busca una sesión completada con al menos `minPerfectAnswers`
// respuestas y ninguna incorrecta.
func (s *gormStore) HasPerfectSession(userID uuid.UUID) (bool, error) {
	var exists bool
	err := s.db.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM sessions s
			WHERE s.user_id = ? AND s.status = ?
			  AND (SELECT COUNT(*) FROM session_answers a WHERE a.session_id = s.id) >= ?
			  AND NOT EXISTS (
				SELECT 1 FROM session_answers a
				WHERE a.session_id = s.id AND a.is_correct = false
			  )
		)`, userID, models.SessionStatusCompleted, minPerfectAnswers).Scan(&exists).Error
	return exists, err
}

// HasFastSession exige además que la media de `response_time_ms` sea <= 10 s.
// `NULLIF(..., 0)` ignora respuestas sin tiempo medido (clientes antiguos):
// si no hay ningún tiempo real, AVG es NULL y el logro no se concede.
func (s *gormStore) HasFastSession(userID uuid.UUID) (bool, error) {
	var exists bool
	err := s.db.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM sessions s
			WHERE s.user_id = ? AND s.status = ?
			  AND (SELECT COUNT(*) FROM session_answers a WHERE a.session_id = s.id) >= ?
			  AND NOT EXISTS (
				SELECT 1 FROM session_answers a
				WHERE a.session_id = s.id AND a.is_correct = false
			  )
			  AND (
				SELECT AVG(NULLIF(a.response_time_ms, 0))
				FROM session_answers a
				WHERE a.session_id = s.id
			  ) <= ?
		)`, userID, models.SessionStatusCompleted, minPerfectAnswers, fastSolverMaxAvgMs).Scan(&exists).Error
	return exists, err
}

// CountMasteredTopics cuenta temas donde el usuario tiene al menos
// `minMasteredQuestions` preguntas respondidas y todas están dominadas (SRS).
func (s *gormStore) CountMasteredTopics(userID uuid.UUID) (int, error) {
	var n int64
	err := s.db.Raw(`
		SELECT COUNT(*)
		FROM (
			SELECT qt.topic_id
			FROM user_question_progress p
			JOIN question_topics qt ON qt.question_id = p.question_id
			WHERE p.user_id = ?
			GROUP BY qt.topic_id
			HAVING COUNT(*) >= ? AND BOOL_AND(p.is_mastered)
		) t`, userID, minMasteredQuestions).Scan(&n).Error
	return int(n), err
}
