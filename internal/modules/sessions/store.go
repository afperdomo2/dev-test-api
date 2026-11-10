package sessions

import (
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Store interface {
	FindPage(userID uuid.UUID, params ListSessionsParams) ([]models.Session, int64, error)
	FindByID(id uuid.UUID) (*models.Session, error)
	Create(session *models.Session) error
	Update(session *models.Session) error
	AddSessionTopics(sessionID uuid.UUID, topicIDs []uuid.UUID) error
	// CreateAnswer inserta la respuesta. Si choca con
	// (session_id, client_answer_id) no inserta y devuelve created=false.
	CreateAnswer(answer *models.SessionAnswer) (created bool, err error)
	// FindByClientAnswerID recupera la respuesta ya guardada de un reintento.
	FindByClientAnswerID(sessionID, clientAnswerID uuid.UUID) (*models.SessionAnswer, error)
	FindAnsweredQuestionIDs(sessionID uuid.UUID) ([]uuid.UUID, error)
	FindNextQuestion(topicIDs []uuid.UUID, answeredIDs []uuid.UUID, difficulty string, mode string, userID uuid.UUID) (*models.Question, error)
	CountAvailableQuestions(topicIDs []uuid.UUID, answeredIDs []uuid.UUID, difficulty string, mode string, userID uuid.UUID) (int64, error)
	FindQuestionByID(id uuid.UUID) (*models.Question, error)
	FindNextDailyQuestion(userID uuid.UUID, answeredIDs []uuid.UUID, skipTrueFalse bool) (*models.Question, error)
	// CountSessionAnswersByType cuenta las respuestas de una sesión cuyo tipo de
	// pregunta coincide. Topa las preguntas de calentamiento por sesión.
	CountSessionAnswersByType(sessionID uuid.UUID, questionType string) (int64, error)
	CountDueReviews(userID uuid.UUID) (int64, error)
	CountWarmupCandidates(userID uuid.UUID) (int64, error)
	// CountPracticeCandidates cuenta el banco visible que puede usarse como
	// práctica normal (cualquier tipo, no dominadas). Es el último escalón del
	// repaso diario y el que decide si hay algo que repasar.
	CountPracticeCandidates(userID uuid.UUID) (int64, error)
	// FindTodayInProgressDaily recupera la sesión diaria de hoy que sigue en
	// progreso, para reutilizarla en vez de crear una nueva en cada toque.
	FindTodayInProgressDaily(userID uuid.UUID, since time.Time) (*models.Session, error)
	FindSummary(id uuid.UUID) (*SessionSummaryData, error)
	Delete(id uuid.UUID) error
}

type gormStore struct {
	db *gorm.DB
}

// excludedSessionQuestionType queda fuera de la selección de sesiones: ningún
// cliente puede responderla y el backend no puede corregirla todavía (D-56).
const excludedSessionQuestionType = "code_completion"

func NewStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) FindPage(userID uuid.UUID, params ListSessionsParams) ([]models.Session, int64, error) {
	var sessions []models.Session
	var total int64

	base := s.db.Where("user_id = ?", userID)
	if params.Status != "" {
		base = base.Where("status = ?", params.Status)
	}
	base.Model(&models.Session{}).Count(&total)

	err := base.Offset((params.Page - 1) * params.PerPage).Limit(params.PerPage).
		Preload("Topics").
		Preload("Answers").
		Order(sortConfig.OrderClause(params.SortBy, params.SortOrder)).
		Find(&sessions).Error
	return sessions, total, err
}

func (s *gormStore) FindByID(id uuid.UUID) (*models.Session, error) {
	var session models.Session
	err := s.db.Preload("Topics").
		Preload("Answers.Question.Options").
		Preload("Answers.Question.CodeChallenge").
		Preload("Answers.Question.Topics").
		First(&session, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (s *gormStore) Create(session *models.Session) error {
	return s.db.Create(session).Error
}

func (s *gormStore) Update(session *models.Session) error {
	return s.db.Save(session).Error
}

func (s *gormStore) AddSessionTopics(sessionID uuid.UUID, topicIDs []uuid.UUID) error {
	for _, topicID := range topicIDs {
		if err := s.db.Create(&models.SessionTopic{
			SessionID: sessionID,
			TopicID:   topicID,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *gormStore) CreateAnswer(answer *models.SessionAnswer) (bool, error) {
	tx := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "session_id"}, {Name: "client_answer_id"}},
		DoNothing: true,
	}).Create(answer)

	if tx.Error != nil {
		return false, tx.Error
	}
	return tx.RowsAffected > 0, nil
}

func (s *gormStore) FindByClientAnswerID(sessionID, clientAnswerID uuid.UUID) (*models.SessionAnswer, error) {
	var answer models.SessionAnswer
	err := s.db.
		Where("session_id = ? AND client_answer_id = ?", sessionID, clientAnswerID).
		First(&answer).Error
	if err != nil {
		return nil, err
	}
	return &answer, nil
}

func (s *gormStore) FindAnsweredQuestionIDs(sessionID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := s.db.Model(&models.SessionAnswer{}).
		Where("session_id = ?", sessionID).
		Pluck("question_id", &ids).Error
	return ids, err
}

func (s *gormStore) FindNextQuestion(topicIDs []uuid.UUID, answeredIDs []uuid.UUID, difficulty string, mode string, userID uuid.UUID) (*models.Question, error) {
	query := s.db.Model(&models.Question{}).
		Joins("JOIN question_topics ON question_topics.question_id = questions.id").
		Where("question_topics.topic_id IN ?", topicIDs).
		Where("questions.deleted_at IS NULL").
		Where("questions.type <> ?", excludedSessionQuestionType)

	if difficulty != "" && difficulty != "random" {
		query = query.Where("questions.difficulty = ?", difficulty)
	}
	if len(answeredIDs) > 0 {
		query = query.Where("questions.id NOT IN ?", answeredIDs)
	}

	query = query.Where("(questions.source = ? OR questions.user_id = ? OR questions.is_public = ?)", "ai_generated", userID, true)

	var question models.Question
	err := query.Order("RANDOM()").
		Preload("Options").
		Preload("CodeChallenge").
		Preload("Topics").
		First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (s *gormStore) CountAvailableQuestions(topicIDs []uuid.UUID, answeredIDs []uuid.UUID, difficulty string, mode string, userID uuid.UUID) (int64, error) {
	query := s.db.Model(&models.Question{}).
		Joins("JOIN question_topics ON question_topics.question_id = questions.id").
		Where("question_topics.topic_id IN ?", topicIDs).
		Where("questions.deleted_at IS NULL").
		Where("questions.type <> ?", excludedSessionQuestionType)

	if difficulty != "" && difficulty != "random" {
		query = query.Where("questions.difficulty = ?", difficulty)
	}
	if len(answeredIDs) > 0 {
		query = query.Where("questions.id NOT IN ?", answeredIDs)
	}

	query = query.Where("(questions.source = ? OR questions.user_id = ? OR questions.is_public = ?)", "ai_generated", userID, true)

	var count int64
	err := query.Count(&count).Error
	return count, err
}

func (s *gormStore) FindQuestionByID(id uuid.UUID) (*models.Question, error) {
	var question models.Question
	err := s.db.Preload("Options").First(&question, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (s *gormStore) FindNextDailyQuestion(userID uuid.UUID, answeredIDs []uuid.UUID, skipTrueFalse bool) (*models.Question, error) {
	// Escalón 1: repaso espaciado (vencidas). Si ya se agotó el cupo de
	// true_false de la sesión, también se saltan las vencidas de ese tipo.
	due, err := s.findNextDueReview(userID, answeredIDs, skipTrueFalse)
	if err == nil {
		return due, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Escalón 2: calentamiento (verdaderas/falsas), solo mientras quede cupo.
	if !skipTrueFalse {
		warmup, err := s.findNextWarmup(userID, answeredIDs)
		if err == nil {
			return warmup, nil
		}
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	// Escalón 3: práctica normal. Sin este fallback una sesión diaria sin
	// vencidas ni verdaderas/falsas se queda vacía para siempre.
	return s.findNextPractice(userID, answeredIDs)
}

func (s *gormStore) findNextDueReview(userID uuid.UUID, answeredIDs []uuid.UUID, skipTrueFalse bool) (*models.Question, error) {
	query := s.db.Model(&models.UserQuestionProgress{}).
		Joins("JOIN questions ON questions.id = user_question_progress.question_id AND questions.deleted_at IS NULL").
		Where("user_question_progress.user_id = ?", userID).
		Where("user_question_progress.next_review_at IS NOT NULL").
		Where("user_question_progress.next_review_at <= now()").
		Where("user_question_progress.is_mastered = false").
		Where("questions.type <> ?", excludedSessionQuestionType).
		Where("(questions.source = ? OR questions.user_id = ? OR questions.is_public = ?)", "ai_generated", userID, true)
	if skipTrueFalse {
		query = query.Where("questions.type <> ?", "true_false")
	}
	if len(answeredIDs) > 0 {
		query = query.Where("user_question_progress.question_id NOT IN ?", answeredIDs)
	}

	var due models.UserQuestionProgress
	err := query.Order("user_question_progress.next_review_at ASC").
		Preload("Question.Options").
		Preload("Question.CodeChallenge").
		Preload("Question.Topics").
		First(&due).Error
	if err != nil {
		return nil, err
	}
	return due.Question, nil
}

func (s *gormStore) findNextWarmup(userID uuid.UUID, answeredIDs []uuid.UUID) (*models.Question, error) {
	query := s.db.Model(&models.Question{}).
		Where("deleted_at IS NULL").
		Where("type = ?", "true_false").
		Where("(questions.source = ? OR questions.user_id = ? OR questions.is_public = ?)", "ai_generated", userID, true)
	if len(answeredIDs) > 0 {
		query = query.Where("questions.id NOT IN ?", answeredIDs)
	}

	var question models.Question
	err := query.Order("RANDOM()").
		Preload("Options").
		Preload("CodeChallenge").
		Preload("Topics").
		First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

// practiceCandidatesQuery es el banco visible que aún sirve como práctica
// normal: cualquier tipo, excluyendo las preguntas ya dominadas.
func (s *gormStore) practiceCandidatesQuery(userID uuid.UUID) *gorm.DB {
	return s.db.Model(&models.Question{}).
		Joins("LEFT JOIN user_question_progress ON user_question_progress.question_id = questions.id AND user_question_progress.user_id = ?", userID).
		Where("questions.deleted_at IS NULL").
		Where("questions.type <> ?", excludedSessionQuestionType).
		Where("(questions.source = ? OR questions.user_id = ? OR questions.is_public = ?)", "ai_generated", userID, true).
		Where("(user_question_progress.question_id IS NULL OR user_question_progress.is_mastered = false)")
}

func (s *gormStore) findNextPractice(userID uuid.UUID, answeredIDs []uuid.UUID) (*models.Question, error) {
	query := s.practiceCandidatesQuery(userID)
	if len(answeredIDs) > 0 {
		query = query.Where("questions.id NOT IN ?", answeredIDs)
	}

	var question models.Question
	err := query.Order("RANDOM()").
		Preload("Options").
		Preload("CodeChallenge").
		Preload("Topics").
		First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (s *gormStore) CountDueReviews(userID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.Model(&models.UserQuestionProgress{}).
		Joins("JOIN questions ON questions.id = user_question_progress.question_id AND questions.deleted_at IS NULL").
		Where("user_question_progress.user_id = ?", userID).
		Where("user_question_progress.next_review_at IS NOT NULL").
		Where("user_question_progress.next_review_at <= now()").
		Where("user_question_progress.is_mastered = false").
		Where("questions.type <> ?", excludedSessionQuestionType).
		Where("(questions.source = ? OR questions.user_id = ? OR questions.is_public = ?)", "ai_generated", userID, true).
		Count(&count).Error
	return count, err
}

func (s *gormStore) CountWarmupCandidates(userID uuid.UUID) (int64, error) {
	var count int64
	err := s.db.Model(&models.Question{}).
		Where("deleted_at IS NULL").
		Where("type = ?", "true_false").
		Where("(source = ? OR user_id = ? OR is_public = ?)", "ai_generated", userID, true).
		Count(&count).Error
	return count, err
}

func (s *gormStore) CountPracticeCandidates(userID uuid.UUID) (int64, error) {
	var count int64
	err := s.practiceCandidatesQuery(userID).Count(&count).Error
	return count, err
}

func (s *gormStore) CountSessionAnswersByType(sessionID uuid.UUID, questionType string) (int64, error) {
	var count int64
	err := s.db.Model(&models.SessionAnswer{}).
		Joins("JOIN questions ON questions.id = session_answers.question_id").
		Where("session_answers.session_id = ?", sessionID).
		Where("questions.type = ?", questionType).
		Count(&count).Error
	return count, err
}

func (s *gormStore) FindTodayInProgressDaily(userID uuid.UUID, since time.Time) (*models.Session, error) {
	var session models.Session
	err := s.db.
		Where("user_id = ? AND mode = ? AND status = ? AND created_at >= ?", userID, "daily", "in_progress", since).
		Order("created_at DESC").
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

type SessionSummaryData struct {
	Status        string
	QuestionLimit *int
	AnswerCount   int64
}

func (s *gormStore) Delete(id uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("session_id = ?", id).Delete(&models.SessionAnswer{}).Error; err != nil {
			return err
		}
		if err := tx.Where("session_id = ?", id).Delete(&models.SessionTopic{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Session{ID: id}).Error
	})
}

func (s *gormStore) FindSummary(id uuid.UUID) (*SessionSummaryData, error) {
	var data SessionSummaryData
	err := s.db.Raw(`
		SELECT s.status, s.question_limit,
			(SELECT COUNT(*) FROM session_answers WHERE session_id = ?) as answer_count
		FROM sessions s WHERE s.id = ?
	`, id, id).Scan(&data).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}
