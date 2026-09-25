package generation

import (
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/services/ai"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// dailyBudgetLockKey es la clave del lock consultivo de Postgres que serializa
// el encolado de trabajos. Así dos requests simultáneos no pueden sobrepasar el
// tope diario por una condición de carrera entre el SELECT y el INSERT.
const dailyBudgetLockKey int64 = 837421

// DailyBudget resume el consumo de la generación con IA del día en curso.
type DailyBudget struct {
	// Reserved es la cantidad comprometida (pending + running + completed).
	Reserved int
	// Generated es la cantidad realmente guardada en el banco.
	Generated int
}

type Store interface {
	CreateJob(job *models.GenerationJob, topicIDs []uuid.UUID) error
	FindPage(params ListJobsParams) ([]models.GenerationJob, int64, error)
	FindByID(id uuid.UUID) (*models.GenerationJob, error)
	TopicNames(jobID uuid.UUID) ([]string, error)
	ClaimNextPending() (*models.GenerationJob, error)
	// LoadJobTopics rellena JobTopics. Necesario tras un claim, que llega sin
	// asociaciones preloaded por venir de un RETURNING de SQL crudo.
	LoadJobTopics(job *models.GenerationJob) error
	CompleteJob(id uuid.UUID, generated int) error
	FailJob(id uuid.UUID, cause string) error
	RecoverRunning() (int64, error)
	DailyBudget(since time.Time) (*DailyBudget, error)
	StatusCounts() (map[string]int64, error)
	ExistingContentByTopics(topicIDs []uuid.UUID, limit int) ([]string, error)
	// EnqueueWithBudget hace el INSERT del job teniendo ya tomado el lock
	// consultivo del tope diario.
	EnqueueWithBudget(job *models.GenerationJob, topicIDs []uuid.UUID) (int, error)
	SaveQuestions(userID uuid.UUID, topicIDs []uuid.UUID, questions []*ai.GeneratedQuestion) (int, error)
	TopicsExist(topicIDs []uuid.UUID) ([]uuid.UUID, error)
	TopicsByIDs(topicIDs []uuid.UUID) ([]models.Topic, error)
	TopicNamesByJobs(jobIDs []uuid.UUID) (map[uuid.UUID][]string, error)
}

type gormStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) CreateJob(job *models.GenerationJob, topicIDs []uuid.UUID) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return createJobTx(tx, job, topicIDs)
	})
}

func createJobTx(tx *gorm.DB, job *models.GenerationJob, topicIDs []uuid.UUID) error {
	if err := tx.Omit("Topics", "JobTopics").Create(job).Error; err != nil {
		return err
	}
	return attachJobTopicsTx(tx, job.ID, topicIDs)
}

func attachJobTopicsTx(tx *gorm.DB, jobID uuid.UUID, topicIDs []uuid.UUID) error {
	if len(topicIDs) == 0 {
		return nil
	}
	rows := make([]models.GenerationJobTopic, 0, len(topicIDs))
	for _, id := range topicIDs {
		rows = append(rows, models.GenerationJobTopic{JobID: jobID, TopicID: id})
	}
	return tx.Create(&rows).Error
}

func (s *gormStore) EnqueueWithBudget(job *models.GenerationJob, topicIDs []uuid.UUID) (int, error) {
	reserved := 0
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", dailyBudgetLockKey).Error; err != nil {
			return err
		}

		since := startOfDay(time.Now().UTC())
		var current int
		if err := tx.Model(&models.GenerationJob{}).
			Where("created_at >= ?", since).
			Where("status <> ?", models.GenerationJobStatusFailed).
			Select("COALESCE(SUM(quantity), 0)").Scan(&current).Error; err != nil {
			return err
		}
		reserved = current

		return createJobTx(tx, job, topicIDs)
	})
	return reserved, err
}

func (s *gormStore) FindPage(params ListJobsParams) ([]models.GenerationJob, int64, error) {
	var jobs []models.GenerationJob
	var total int64

	query := s.db.Model(&models.GenerationJob{})
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset((params.Page - 1) * params.PerPage).Limit(params.PerPage).
		Order(sortConfig.OrderClause(params.SortBy, params.SortOrder)).
		Preload("Topics").
		Find(&jobs).Error
	return jobs, total, err
}

func (s *gormStore) FindByID(id uuid.UUID) (*models.GenerationJob, error) {
	var job models.GenerationJob
	if err := s.db.Preload("Topics").First(&job, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (s *gormStore) TopicNames(jobID uuid.UUID) ([]string, error) {
	var names []string
	err := s.db.Model(&models.Topic{}).
		Joins("JOIN generation_job_topics ON generation_job_topics.topic_id = topics.id").
		Where("generation_job_topics.job_id = ?", jobID).
		Order("topics.name ASC").
		Pluck("topics.name", &names).Error
	return names, err
}

// ClaimNextPending reclama atómicamente el trabajo pending más antiguo. El
// subquery con FOR UPDATE SKIP LOCKED hace que dos workers concurrentes nunca
// reclamen el mismo job.
func (s *gormStore) ClaimNextPending() (*models.GenerationJob, error) {
	var job models.GenerationJob
	err := s.db.Raw(`
		UPDATE generation_jobs
		SET status = ?, started_at = now(), updated_at = now()
		WHERE id = (
			SELECT id FROM generation_jobs
			WHERE status = ?
			ORDER BY created_at ASC
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		RETURNING *
	`, models.GenerationJobStatusRunning, models.GenerationJobStatusPending).
		Scan(&job).Error
	if err != nil {
		return nil, err
	}
	if job.ID == uuid.Nil {
		return nil, gorm.ErrRecordNotFound
	}
	if err := s.LoadJobTopics(&job); err != nil {
		return nil, err
	}
	return &job, nil
}

func (s *gormStore) LoadJobTopics(job *models.GenerationJob) error {
	return s.db.Where("job_id = ?", job.ID).Find(&job.JobTopics).Error
}

func (s *gormStore) CompleteJob(id uuid.UUID, generated int) error {
	now := time.Now().UTC()
	return s.db.Model(&models.GenerationJob{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":      models.GenerationJobStatusCompleted,
			"generated":   generated,
			"error":       nil,
			"finished_at": now,
			"updated_at":  now,
		}).Error
}

func (s *gormStore) FailJob(id uuid.UUID, cause string) error {
	now := time.Now().UTC()
	return s.db.Model(&models.GenerationJob{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":      models.GenerationJobStatusFailed,
			"error":       cause,
			"finished_at": now,
			"updated_at":  now,
		}).Error
}

// RecoverRunning devuelve a pending los jobs que quedaron en running por un
// crash o reinicio del servidor.
func (s *gormStore) RecoverRunning() (int64, error) {
	result := s.db.Model(&models.GenerationJob{}).
		Where("status = ?", models.GenerationJobStatusRunning).
		Updates(map[string]any{
			"status":     models.GenerationJobStatusPending,
			"started_at": nil,
			"updated_at": time.Now().UTC(),
		})
	return result.RowsAffected, result.Error
}

func (s *gormStore) DailyBudget(since time.Time) (*DailyBudget, error) {
	budget := &DailyBudget{}

	if err := s.db.Model(&models.GenerationJob{}).
		Where("created_at >= ?", since).
		Where("status <> ?", models.GenerationJobStatusFailed).
		Select("COALESCE(SUM(quantity), 0)").Scan(&budget.Reserved).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.GenerationJob{}).
		Where("created_at >= ?", since).
		Select("COALESCE(SUM(generated), 0)").Scan(&budget.Generated).Error; err != nil {
		return nil, err
	}

	return budget, nil
}

func (s *gormStore) StatusCounts() (map[string]int64, error) {
	type row struct {
		Status string
		Count  int64
	}
	var rows []row
	err := s.db.Model(&models.GenerationJob{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64, len(rows))
	for _, r := range rows {
		counts[r.Status] = r.Count
	}
	return counts, nil
}

// ExistingContentByTopics devuelve los enunciados ya existentes en los temas,
// para inyectarlos en el prompt y evitar duplicados.
func (s *gormStore) ExistingContentByTopics(topicIDs []uuid.UUID, limit int) ([]string, error) {
	if len(topicIDs) == 0 {
		return nil, nil
	}
	var contents []string
	err := s.db.Model(&models.Question{}).
		Joins("JOIN question_topics ON question_topics.question_id = questions.id").
		Where("question_topics.topic_id IN ?", topicIDs).
		Where("questions.deleted_at IS NULL").
		Order("questions.created_at DESC").
		Limit(limit).
		Pluck("questions.content", &contents).Error
	return contents, err
}

// SaveQuestions persiste un lote de preguntas generadas en una única
// transacción: o se guardan todas, o no se guarda ninguna. Devuelve cuántas
// quedaron guardadas.
func (s *gormStore) SaveQuestions(userID uuid.UUID, topicIDs []uuid.UUID, questions []*ai.GeneratedQuestion) (int, error) {
	if len(questions) == 0 {
		return 0, nil
	}

	saved := 0
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, gen := range questions {
			question := &models.Question{
				UserID:      userID,
				Type:        gen.Type,
				Content:     gen.Content,
				Explanation: gen.Explanation,
				Difficulty:  gen.Difficulty,
				IsPublic:    true,
				Source:      "ai_generated",
			}

			if gen.Type == "code_completion" {
				question.CodeChallenge = &models.CodeChallenge{
					StarterCode:    gen.StarterCode,
					ExpectedOutput: gen.ExpectedOutput,
					Language:       gen.Language,
				}
			} else {
				for _, opt := range gen.Options {
					question.Options = append(question.Options, models.QuestionOption{
						Content:   opt.Content,
						IsCorrect: opt.IsCorrect,
					})
				}
			}

			if err := tx.Omit("Topics", "QuestionTopics").Create(question).Error; err != nil {
				return err
			}
			for _, topicID := range topicIDs {
				if err := tx.Create(&models.QuestionTopic{
					QuestionID: question.ID,
					TopicID:    topicID,
				}).Error; err != nil {
					return err
				}
			}
			saved++
		}
		return nil
	})
	return saved, err
}

// TopicsExist filtra los IDs que no existen, para responder 422 con precision
// en vez de fallar con una violación de FK.
func (s *gormStore) TopicsExist(topicIDs []uuid.UUID) ([]uuid.UUID, error) {
	if len(topicIDs) == 0 {
		return nil, nil
	}
	var found []uuid.UUID
	err := s.db.Model(&models.Topic{}).
		Where("id IN ?", topicIDs).
		Pluck("id", &found).Error
	return found, err
}

func (s *gormStore) TopicsByIDs(topicIDs []uuid.UUID) ([]models.Topic, error) {
	if len(topicIDs) == 0 {
		return nil, nil
	}
	var topics []models.Topic
	err := s.db.Where("id IN ?", topicIDs).Find(&topics).Error
	return topics, err
}

// TopicNamesByJobs agrupa los nombres de temas por trabajo, para no disparar
// una consulta por trabajo al armar la respuesta paginada.
func (s *gormStore) TopicNamesByJobs(jobIDs []uuid.UUID) (map[uuid.UUID][]string, error) {
	out := make(map[uuid.UUID][]string, len(jobIDs))
	if len(jobIDs) == 0 {
		return out, nil
	}

	type row struct {
		JobID uuid.UUID
		Name  string
	}
	var rows []row
	err := s.db.Model(&models.Topic{}).
		Select("generation_job_topics.job_id AS job_id, topics.name AS name").
		Joins("JOIN generation_job_topics ON generation_job_topics.topic_id = topics.id").
		Where("generation_job_topics.job_id IN ?", jobIDs).
		Order("topics.name ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, r := range rows {
		out[r.JobID] = append(out[r.JobID], r.Name)
	}
	return out, nil
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
