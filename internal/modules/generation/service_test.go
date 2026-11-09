package generation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/services/ai"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type fakeStore struct {
	existingTopicIDs []uuid.UUID
	topics           []models.Topic
	budget           *DailyBudget
	counts           map[string]int64

	enqueueErr    error
	enqueuedJob   *models.GenerationJob
	enqueuedTopic []uuid.UUID
	// reservedBefore es el valor que devuelve EnqueueWithBudget (ya reservado
	// por otros trabajos, sin contar el nuevo).
	reservedBefore int

	claimed   *models.GenerationJob
	claimErr  error
	saved     []*ai.GeneratedQuestion
	saveErr   error
	saveCount int

	savedTips      []*ai.GeneratedTip
	savedTipUser   uuid.UUID
	savedTipTopics []models.Topic
	saveTipsErr    error
	saveTipsN      int

	completedID uuid.UUID
	completedN  int
	failedID    uuid.UUID
	failedErr   string
}

// stubGenerator evita la red en los tests del módulo.
type stubGenerator struct {
	configured bool
	questions  []*ai.GeneratedQuestion
	err        error
	gotReq     ai.BatchRequest

	tips      []*ai.GeneratedTip
	tipsErr   error
	gotTipReq ai.TipBatchRequest
}

func (g *stubGenerator) IsConfigured() bool { return g.configured }

func (g *stubGenerator) GenerateQuestions(_ context.Context, req ai.BatchRequest) ([]*ai.GeneratedQuestion, error) {
	g.gotReq = req
	if g.err != nil {
		return nil, g.err
	}
	return g.questions, nil
}

func (g *stubGenerator) GenerateTips(_ context.Context, req ai.TipBatchRequest) ([]*ai.GeneratedTip, error) {
	g.gotTipReq = req
	if g.tipsErr != nil {
		return nil, g.tipsErr
	}
	return g.tips, nil
}

func (s *fakeStore) CreateJob(job *models.GenerationJob, topicIDs []uuid.UUID) error {
	s.enqueuedJob, s.enqueuedTopic = job, topicIDs
	return nil
}

func (s *fakeStore) EnqueueWithBudget(job *models.GenerationJob, topicIDs []uuid.UUID) (int, error) {
	if s.enqueueErr != nil {
		return 0, s.enqueueErr
	}
	s.enqueuedJob, s.enqueuedTopic = job, topicIDs
	return s.reservedBefore, nil
}

func (s *fakeStore) FindPage(ListJobsParams) ([]models.GenerationJob, int64, error) {
	return nil, 0, nil
}

func (s *fakeStore) FindByID(uuid.UUID) (*models.GenerationJob, error) {
	return nil, gorm.ErrRecordNotFound
}

func (s *fakeStore) TopicNames(uuid.UUID) ([]string, error) { return nil, nil }
func (s *fakeStore) TopicNamesByJobs([]uuid.UUID) (map[uuid.UUID][]string, error) {
	return map[uuid.UUID][]string{}, nil
}

func (s *fakeStore) ClaimNextPending() (*models.GenerationJob, error) { return s.claimed, s.claimErr }
func (s *fakeStore) LoadJobTopics(*models.GenerationJob) error        { return nil }
func (s *fakeStore) CompleteJob(id uuid.UUID, generated int) error {
	s.completedID, s.completedN = id, generated
	return nil
}

func (s *fakeStore) FailJob(id uuid.UUID, cause string) error {
	s.failedID, s.failedErr = id, cause
	return nil
}
func (s *fakeStore) RecoverRunning() (int64, error) { return 0, nil }
func (s *fakeStore) DailyBudget(time.Time) (*DailyBudget, error) {
	if s.budget == nil {
		return &DailyBudget{}, nil
	}
	return s.budget, nil
}
func (s *fakeStore) StatusCounts(string) (map[string]int64, error) {
	return s.counts, nil
}
func (s *fakeStore) ExistingContentByTopics([]uuid.UUID, int) ([]string, error) {
	return nil, nil
}
func (s *fakeStore) ExistingTipContentByTopics([]uuid.UUID, int) ([]string, error) {
	return nil, nil
}
func (s *fakeStore) TopicsExist([]uuid.UUID) ([]uuid.UUID, error) { return s.existingTopicIDs, nil }
func (s *fakeStore) TopicsByIDs([]uuid.UUID) ([]models.Topic, error) {
	return s.topics, nil
}
func (s *fakeStore) SaveQuestions(_ uuid.UUID, _ []uuid.UUID, questions []*ai.GeneratedQuestion) (int, error) {
	if s.saveErr != nil {
		return 0, s.saveErr
	}
	s.saved = questions
	return s.saveCount, nil
}

func (s *fakeStore) SaveTips(userID uuid.UUID, topics []models.Topic, tips []*ai.GeneratedTip) (int, error) {
	if s.saveTipsErr != nil {
		return 0, s.saveTipsErr
	}
	s.savedTipUser = userID
	s.savedTipTopics = topics
	s.savedTips = tips
	return s.saveTipsN, nil
}

type fakeSettings struct {
	maxPerJob  int
	maxPerDay  int
	defaultQty int
}

func (f fakeSettings) MaxPerJob() (int, error)       { return f.maxPerJob, nil }
func (f fakeSettings) MaxPerDay() (int, error)       { return f.maxPerDay, nil }
func (f fakeSettings) DefaultQuantity() (int, error) { return f.defaultQty, nil }

func newTestService(store *fakeStore, settings fakeSettings, gen AIGenerator) *service {
	return &service{store: store, generator: gen, settings: settings}
}

func okGen(questions int) *stubGenerator {
	list := make([]*ai.GeneratedQuestion, 0, questions)
	for range questions {
		list = append(list, &ai.GeneratedQuestion{
			Type:       "single_choice",
			Content:    "¿Qué es esto?",
			Difficulty: "beginner",
			Options:    []ai.GeneratedOption{{Content: "a", IsCorrect: true}},
		})
	}
	return &stubGenerator{configured: true, questions: list}
}

func TestCreateJobUsesDefaultQuantityWhenOmitted(t *testing.T) {
	topicID := uuid.New()
	store := &fakeStore{existingTopicIDs: []uuid.UUID{topicID}}
	svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, okGen(5))
	svc.generator = okGen(5)

	job, err := svc.CreateJob(context.Background(), uuid.New(), CreateJobRequest{
		TopicIDs:   []uuid.UUID{topicID},
		Difficulty: "beginner",
	})
	require.NoError(t, err)
	require.NotNil(t, store.enqueuedJob)
	assert.Equal(t, 5, store.enqueuedJob.Quantity, "sin quantity explícito se usa ai.default_quantity")
	assert.Equal(t, models.GenerationJobStatusPending, store.enqueuedJob.Status)
	require.NotNil(t, job)
}

func TestCreateJobRejectsQuantityAboveMaxPerJob(t *testing.T) {
	topicID := uuid.New()
	store := &fakeStore{existingTopicIDs: []uuid.UUID{topicID}}
	svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, okGen(5))

	_, err := svc.CreateJob(context.Background(), uuid.New(), CreateJobRequest{
		TopicIDs:   []uuid.UUID{topicID},
		Difficulty: "beginner",
		Quantity:   11,
	})
	requireAPIError(t, err, apierr.ErrValidation)
}

func TestCreateJobRejectsUnknownTopics(t *testing.T) {
	// El store devuelve menos topics de los pedidos: uno no existe.
	store := &fakeStore{existingTopicIDs: []uuid.UUID{uuid.New()}}
	svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, okGen(5))

	_, err := svc.CreateJob(context.Background(), uuid.New(), CreateJobRequest{
		TopicIDs:   []uuid.UUID{uuid.New(), uuid.New()},
		Difficulty: "beginner",
	})
	requireAPIError(t, err, apierr.ErrValidation)
}

func TestCreateJobRejectsWhenDailyCapWouldBeExceeded(t *testing.T) {
	topicID := uuid.New()
	store := &fakeStore{existingTopicIDs: []uuid.UUID{topicID}, reservedBefore: 198}
	svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, okGen(5))

	// 198 reservados + 5 = 203 > 200: debe rechazar con 429.
	_, err := svc.CreateJob(context.Background(), uuid.New(), CreateJobRequest{
		TopicIDs:   []uuid.UUID{topicID},
		Difficulty: "beginner",
		Quantity:   5,
	})
	requireAPIError(t, err, apierr.ErrTooManyRequests)
}

func TestCreateJobFitsExactlyAtDailyCap(t *testing.T) {
	topicID := uuid.New()
	store := &fakeStore{existingTopicIDs: []uuid.UUID{topicID}, reservedBefore: 195}
	svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, okGen(5))

	_, err := svc.CreateJob(context.Background(), uuid.New(), CreateJobRequest{
		TopicIDs:   []uuid.UUID{topicID},
		Difficulty: "beginner",
		Quantity:   5,
	})
	require.NoError(t, err, "195 + 5 = 200 cabe exactamente en el tope")
}

func TestListJobsRejectsUnknownStatus(t *testing.T) {
	svc := newTestService(&fakeStore{}, fakeSettings{}, okGen(1))

	_, _, err := svc.ListJobs(ListJobsParams{Status: "inventado"})
	requireAPIError(t, err, apierr.ErrValidation)
}

func TestStatsReportsReservedAndGeneratedSeparately(t *testing.T) {
	store := &fakeStore{
		budget: &DailyBudget{Reserved: 12, Generated: 7},
		counts: map[string]int64{
			models.GenerationJobStatusPending: 1,
			models.GenerationJobStatusRunning: 2,
			models.GenerationJobStatusFailed:  3,
		},
	}
	svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, okGen(1))

	stats, err := svc.Stats("")
	require.NoError(t, err)
	assert.Equal(t, 12, stats.ReservedToday)
	assert.Equal(t, 7, stats.GeneratedToday)
	assert.Equal(t, 200, stats.DailyCap)
	assert.Equal(t, 10, stats.PerJobCap)
	assert.Equal(t, 1, stats.Pending)
	assert.Equal(t, 2, stats.Running)
	assert.Equal(t, 3, stats.Failed)
}

func TestProcessJobMarksCompletedWithSavedCount(t *testing.T) {
	topicID := uuid.New()
	store := &fakeStore{
		topics:    []models.Topic{{ID: topicID, Name: "Go"}},
		saveCount: 3,
	}
	gen := okGen(5)
	svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, gen)

	jobID := uuid.New()
	job := models.GenerationJob{
		ID:         jobID,
		Quantity:   5,
		Difficulty: "beginner",
		JobTopics:  []models.GenerationJobTopic{{JobID: jobID, TopicID: topicID}},
	}

	require.NoError(t, svc.ProcessJob(context.Background(), job))
	assert.Equal(t, jobID, store.completedID)
	assert.Equal(t, 3, store.completedN, "se completa con cuántas se guardaron, no cuántas se pidieron")
	assert.Equal(t, 5, gen.gotReq.Count, "la IA recibe una sola llamada por la cantidad pedida")
	assert.Equal(t, []string{"Go"}, gen.gotReq.TopicNames)
}

func TestProcessJobDoesNotCompleteWhenGenerationFails(t *testing.T) {
	topicID := uuid.New()
	store := &fakeStore{topics: []models.Topic{{ID: topicID, Name: "Go"}}}
	gen := &stubGenerator{configured: true, err: errors.New("la IA no respondió")}
	svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, gen)

	jobID := uuid.New()
	job := models.GenerationJob{
		ID:         jobID,
		Quantity:   5,
		Difficulty: "beginner",
		JobTopics:  []models.GenerationJobTopic{{JobID: jobID, TopicID: topicID}},
	}

	err := svc.ProcessJob(context.Background(), job)
	require.Error(t, err)
	assert.Equal(t, uuid.Nil, store.completedID, "no debe marcar completado si falla la IA")
}

func TestProcessJobTipsUsesTipGeneratorAndSavesTips(t *testing.T) {
	topicID := uuid.New()
	store := &fakeStore{
		topics:    []models.Topic{{ID: topicID, Name: "Go", Slug: "go"}},
		saveTipsN: 4,
	}
	gen := &stubGenerator{
		configured: true,
		tips:       []*ai.GeneratedTip{{Content: "Usa defer", Topic: "go"}},
	}
	svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, gen)

	jobID := uuid.New()
	job := models.GenerationJob{
		ID:        jobID,
		Kind:      models.GenerationKindTip,
		Quantity:  5,
		JobTopics: []models.GenerationJobTopic{{JobID: jobID, TopicID: topicID}},
	}

	require.NoError(t, svc.ProcessJob(context.Background(), job))
	assert.Equal(t, jobID, store.completedID)
	assert.Equal(t, 4, store.completedN, "se completa con cuántos tips se guardaron")
	assert.Equal(t, 5, gen.gotTipReq.Count, "una sola llamada por la cantidad pedida")
	assert.Equal(t, []string{"Go"}, gen.gotTipReq.TopicNames)
	assert.Equal(t, []string{"go"}, gen.gotTipReq.TopicSlugs)
	assert.Len(t, store.savedTips, 1)
	assert.Empty(t, gen.gotReq.TopicNames, "un job de tips no llama al generador de preguntas")
}

func TestProcessJobTipsDoesNotCompleteWhenGenerationFails(t *testing.T) {
	topicID := uuid.New()
	store := &fakeStore{topics: []models.Topic{{ID: topicID, Name: "Go", Slug: "go"}}}
	gen := &stubGenerator{configured: true, tipsErr: errors.New("la IA no respondió")}
	svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, gen)

	jobID := uuid.New()
	job := models.GenerationJob{
		ID:        jobID,
		Kind:      models.GenerationKindTip,
		Quantity:  5,
		JobTopics: []models.GenerationJobTopic{{JobID: jobID, TopicID: topicID}},
	}

	require.Error(t, svc.ProcessJob(context.Background(), job))
	assert.Equal(t, uuid.Nil, store.completedID)
}

func TestCreateJobKindDefaultsToQuestionAndAcceptsTip(t *testing.T) {
	topicID := uuid.New()

	t.Run("sin kind se asume question", func(t *testing.T) {
		store := &fakeStore{existingTopicIDs: []uuid.UUID{topicID}}
		svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, okGen(5))

		_, err := svc.CreateJob(context.Background(), uuid.New(), CreateJobRequest{
			TopicIDs: []uuid.UUID{topicID},
		})
		require.NoError(t, err)
		assert.Equal(t, models.GenerationKindQuestion, store.enqueuedJob.Kind)
	})

	t.Run("kind tip ignora la dificultad", func(t *testing.T) {
		store := &fakeStore{existingTopicIDs: []uuid.UUID{topicID}}
		svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, okGen(5))

		_, err := svc.CreateJob(context.Background(), uuid.New(), CreateJobRequest{
			TopicIDs: []uuid.UUID{topicID},
			Kind:     models.GenerationKindTip,
		})
		require.NoError(t, err)
		assert.Equal(t, models.GenerationKindTip, store.enqueuedJob.Kind)
		assert.Equal(t, "intermediate", store.enqueuedJob.Difficulty)
	})

	t.Run("kind inválido se rechaza", func(t *testing.T) {
		topics := []uuid.UUID{uuid.New()}
		store := &fakeStore{existingTopicIDs: topics}
		svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, okGen(1))

		_, err := svc.CreateJob(context.Background(), uuid.New(), CreateJobRequest{
			TopicIDs: topics,
			Kind:     "inventado",
		})
		requireAPIError(t, err, apierr.ErrValidation)
	})
}

func TestCreateJobNormalizesQuestionType(t *testing.T) {
	topicID := uuid.New()

	t.Run("sin questionType se conserva single_choice", func(t *testing.T) {
		store := &fakeStore{existingTopicIDs: []uuid.UUID{topicID}}
		svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, okGen(1))

		_, err := svc.CreateJob(context.Background(), uuid.New(), CreateJobRequest{
			TopicIDs: []uuid.UUID{topicID},
		})
		require.NoError(t, err)
		assert.Equal(t, models.GenerationQuestionTypeSingleChoice, store.enqueuedJob.QuestionType)
	})

	t.Run("questionType válido se persiste", func(t *testing.T) {
		store := &fakeStore{existingTopicIDs: []uuid.UUID{topicID}}
		svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, okGen(1))

		_, err := svc.CreateJob(context.Background(), uuid.New(), CreateJobRequest{
			TopicIDs:     []uuid.UUID{topicID},
			QuestionType: models.GenerationQuestionTypeTrueFalse,
		})
		require.NoError(t, err)
		assert.Equal(t, models.GenerationQuestionTypeTrueFalse, store.enqueuedJob.QuestionType)
	})

	t.Run("questionType inválido se rechaza", func(t *testing.T) {
		store := &fakeStore{existingTopicIDs: []uuid.UUID{topicID}}
		svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, okGen(1))

		_, err := svc.CreateJob(context.Background(), uuid.New(), CreateJobRequest{
			TopicIDs:     []uuid.UUID{topicID},
			QuestionType: "code_completion",
		})
		requireAPIError(t, err, apierr.ErrValidation)
	})
}

func TestProcessJobPassesQuestionTypeToGenerator(t *testing.T) {
	topicID := uuid.New()
	store := &fakeStore{
		topics:    []models.Topic{{ID: topicID, Name: "Go"}},
		saveCount: 2,
	}
	gen := okGen(2)
	svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, gen)

	jobID := uuid.New()
	job := models.GenerationJob{
		ID:           jobID,
		Quantity:     2,
		Difficulty:   "beginner",
		QuestionType: models.GenerationQuestionTypeMultipleChoice,
		JobTopics:    []models.GenerationJobTopic{{JobID: jobID, TopicID: topicID}},
	}

	require.NoError(t, svc.ProcessJob(context.Background(), job))
	assert.Equal(t, models.GenerationQuestionTypeMultipleChoice, gen.gotReq.QuestionType)
}

func TestStatsRejectsUnknownKind(t *testing.T) {
	svc := newTestService(&fakeStore{}, fakeSettings{}, okGen(1))

	_, err := svc.Stats("inventado")
	requireAPIError(t, err, apierr.ErrValidation)
}

func TestCreateJobFailsClosedWhenStoreCannotVerifyBudget(t *testing.T) {
	topicID := uuid.New()
	store := &fakeStore{
		existingTopicIDs: []uuid.UUID{topicID},
		enqueueErr:       errors.New("db caída"),
	}
	svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, okGen(5))

	_, err := svc.CreateJob(context.Background(), uuid.New(), CreateJobRequest{
		TopicIDs:   []uuid.UUID{topicID},
		Difficulty: "beginner",
	})
	requireAPIError(t, err, apierr.ErrInternal)
}

func TestCreateJobRequiresConfiguredAI(t *testing.T) {
	topicID := uuid.New()
	store := &fakeStore{existingTopicIDs: []uuid.UUID{topicID}}
	// Generador sin configurar: no se debe encolar nada.
	gen := &stubGenerator{configured: false}
	svc := newTestService(store, fakeSettings{maxPerJob: 10, maxPerDay: 200, defaultQty: 5}, gen)

	_, err := svc.CreateJob(context.Background(), uuid.New(), CreateJobRequest{
		TopicIDs:   []uuid.UUID{topicID},
		Difficulty: "beginner",
	})
	requireAPIError(t, err, apierr.ErrValidation)
	assert.Nil(t, store.enqueuedJob, "no debe encolar nada sin IA configurada")
}

func TestTopicIDsAreDeduplicated(t *testing.T) {
	id := uuid.New()
	assert.Len(t, uniqueTopicIDs([]uuid.UUID{id, id, id}), 1)
	assert.Len(t, uniqueTopicIDs([]uuid.UUID{id, uuid.New()}), 2)
}

func requireAPIError(t *testing.T, err error, check func(string, string) *apierr.APIError) {
	t.Helper()
	require.Error(t, err)
	apiErr, ok := err.(*apierr.APIError)
	require.True(t, ok, "se esperaba un *apierr.APIError, se obtuvo %T", err)
	expected := check("", "")
	assert.Equal(t, expected.Status, apiErr.Status)
}
