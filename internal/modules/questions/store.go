package questions

import (
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store interface {
	FindPage(params ListQuestionsParams) ([]models.Question, int64, error)
	FindAll() ([]models.Question, error)
	FindByID(id uuid.UUID) (*models.Question, error)
	Create(question *models.Question) error
	Update(question *models.Question) error
	Delete(id uuid.UUID) error
	AddQuestionTopics(questionID uuid.UUID, topicIDs []uuid.UUID) error
	ReplaceQuestionTopics(questionID uuid.UUID, topicIDs []uuid.UUID) error
	ReplaceQuestionOptions(questionID uuid.UUID, options []models.QuestionOption) error
	BulkCreate(questions []*models.Question) error
	CountImportedSince(userID uuid.UUID, since time.Time) (int64, error)
	CountAiGeneratedSince(userID uuid.UUID, since time.Time) (int64, error)
	Stats(isAdmin bool, userID uuid.UUID) (*QuestionStats, error)
}

type gormStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) FindPage(params ListQuestionsParams) ([]models.Question, int64, error) {
	var questions []models.Question
	var total int64

	query := s.db.Model(&models.Question{})

	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.Difficulty != "" {
		query = query.Where("difficulty = ?", params.Difficulty)
	}
	if len(params.TopicIDs) > 0 {
		query = query.Where("id IN (SELECT question_id FROM question_topics WHERE topic_id IN ?)", params.TopicIDs)
	}
	if params.UserID != uuid.Nil {
		query = query.Where("(source = ? OR user_id = ? OR is_public = ?)", "ai_generated", params.UserID, true)
	}

	query.Count(&total)
	err := query.Offset((params.Page - 1) * params.PerPage).Limit(params.PerPage).
		Order(sortConfig.OrderClause(params.SortBy, params.SortOrder)).
		Preload("CodeChallenge").Preload("Topics").
		Find(&questions).Error
	return questions, total, err
}

func (s *gormStore) FindAll() ([]models.Question, error) {
	var questions []models.Question
	err := s.db.Preload("CodeChallenge").Preload("Topics").
		Order("created_at desc").Find(&questions).Error
	return questions, err
}

func (s *gormStore) FindByID(id uuid.UUID) (*models.Question, error) {
	var question models.Question
	err := s.db.Preload("Options").Preload("CodeChallenge").Preload("Topics").
		First(&question, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (s *gormStore) Create(question *models.Question) error {
	return s.db.Create(question).Error
}

func (s *gormStore) Update(question *models.Question) error {
	return s.db.Save(question).Error
}

func (s *gormStore) Delete(id uuid.UUID) error {
	return s.db.Delete(&models.Question{}, "id = ?", id).Error
}

func (s *gormStore) AddQuestionTopics(questionID uuid.UUID, topicIDs []uuid.UUID) error {
	for _, topicID := range topicIDs {
		if err := s.db.Create(&models.QuestionTopic{
			QuestionID: questionID,
			TopicID:    topicID,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *gormStore) ReplaceQuestionTopics(questionID uuid.UUID, topicIDs []uuid.UUID) error {
	if err := s.db.Where("question_id = ?", questionID).Delete(&models.QuestionTopic{}).Error; err != nil {
		return err
	}
	return s.AddQuestionTopics(questionID, topicIDs)
}

func (s *gormStore) ReplaceQuestionOptions(questionID uuid.UUID, options []models.QuestionOption) error {
	if err := s.db.Where("question_id = ?", questionID).Delete(&models.QuestionOption{}).Error; err != nil {
		return err
	}
	for i := range options {
		options[i].QuestionID = questionID
	}
	if len(options) > 0 {
		return s.db.Create(&options).Error
	}
	return nil
}

func (s *gormStore) BulkCreate(questions []*models.Question) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, q := range questions {
			if err := tx.Omit("Topics", "QuestionTopics").Create(q).Error; err != nil {
				return err
			}
			for _, tid := range topicIDsFromQuestion(q) {
				if err := tx.Create(&models.QuestionTopic{
					QuestionID: q.ID,
					TopicID:    tid,
				}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func topicIDsFromQuestion(q *models.Question) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(q.Topics))
	for _, t := range q.Topics {
		ids = append(ids, t.ID)
	}
	return ids
}

func (s *gormStore) CountImportedSince(userID uuid.UUID, since time.Time) (int64, error) {
	var count int64
	err := s.db.Model(&models.Question{}).
		Where("user_id = ?", userID).
		Where("source = ?", "imported").
		Where("created_at >= ?", since).
		Count(&count).Error
	return count, err
}

func (s *gormStore) CountAiGeneratedSince(userID uuid.UUID, since time.Time) (int64, error) {
	var count int64
	err := s.db.Model(&models.Question{}).
		Where("user_id = ?", userID).
		Where("source = ?", "ai_generated").
		Where("created_at >= ?", since).
		Count(&count).Error
	return count, err
}

func (s *gormStore) Stats(isAdmin bool, userID uuid.UUID) (*QuestionStats, error) {
	var visibility *gorm.DB
	if isAdmin {
		visibility = s.db.Model(&models.Question{}).Where("(questions.source = ? OR questions.is_public = ?)", "ai_generated", true)
	} else {
		visibility = s.db.Model(&models.Question{}).Where("(questions.source = ? OR questions.user_id = ? OR questions.is_public = ?)", "ai_generated", userID, true)
	}

	var total int64
	if err := visibility.Count(&total).Error; err != nil {
		return nil, err
	}

	// ByType
	var byType []TypeCount
	if err := visibility.Select("questions.type as type, COUNT(*) as count").Group("questions.type").Order("count DESC").Scan(&byType).Error; err != nil {
		return nil, err
	}
	if byType == nil {
		byType = []TypeCount{}
	}

	// ByDifficulty
	var byDifficulty []DifficultyCount
	if err := visibility.Select("questions.difficulty as difficulty, COUNT(*) as count").Group("questions.difficulty").Order("count DESC").Scan(&byDifficulty).Error; err != nil {
		return nil, err
	}
	if byDifficulty == nil {
		byDifficulty = []DifficultyCount{}
	}

	// ByTopic — join question_topics + topics
	topicWhere, topicArgs := func() (string, []any) {
		if isAdmin {
			return "(questions.source = ? OR questions.is_public = ?)", []any{"ai_generated", true}
		}
		return "(questions.source = ? OR questions.user_id = ? OR questions.is_public = ?)", []any{"ai_generated", userID, true}
	}()
	var byTopic []TopicCount
	if err := s.db.Model(&models.Question{}).
		Select("topics.id as topic_id, topics.name, topics.slug, topics.category, COUNT(*) as count").
		Joins("JOIN question_topics ON question_topics.question_id = questions.id").
		Joins("JOIN topics ON topics.id = question_topics.topic_id").
		Where(topicWhere, topicArgs...).
		Group("topics.id, topics.name, topics.slug, topics.category").
		Order("count DESC, topics.name ASC").
		Scan(&byTopic).Error; err != nil {
		return nil, err
	}
	if byTopic == nil {
		byTopic = []TopicCount{}
	}

	// ByCategory — via topics.category
	var byCategory []CategoryCount
	if err := s.db.Model(&models.Question{}).
		Select("topics.category as category, COUNT(*) as count").
		Joins("JOIN question_topics ON question_topics.question_id = questions.id").
		Joins("JOIN topics ON topics.id = question_topics.topic_id").
		Where(topicWhere, topicArgs...).
		Group("topics.category").
		Order("count DESC").
		Scan(&byCategory).Error; err != nil {
		return nil, err
	}
	if byCategory == nil {
		byCategory = []CategoryCount{}
	}

	return &QuestionStats{
		Total:        total,
		ByTopic:      byTopic,
		ByCategory:   byCategory,
		ByDifficulty: byDifficulty,
		ByType:       byType,
	}, nil
}
