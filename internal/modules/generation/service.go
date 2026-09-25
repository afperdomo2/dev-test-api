package generation

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/services/ai"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// maxExistingContent acota cuántas preguntas previas se cargan para el prompt.
const maxExistingContent = 60

// QuestionGenerator es la parte del servicio de IA que usa este módulo. Es una
// interfaz (no *ai.Generator) para poder testear el módulo sin red.
type QuestionGenerator interface {
	IsConfigured() bool
	GenerateQuestions(ctx context.Context, req ai.BatchRequest) ([]*ai.GeneratedQuestion, error)
}

type SettingsProvider interface {
	MaxPerJob() (int, error)
	MaxPerDay() (int, error)
	DefaultQuantity() (int, error)
}

type Service interface {
	CreateJob(ctx context.Context, requestedBy uuid.UUID, req CreateJobRequest) (*JobResponse, error)
	ListJobs(params ListJobsParams) ([]JobResponse, int64, error)
	GetJob(id uuid.UUID) (*JobResponse, error)
	Stats() (*StatsResponse, error)

	// ProcessJob ejecuta un trabajo ya reclamado. Lo usa el worker, no los handlers.
	ProcessJob(ctx context.Context, job models.GenerationJob) error
	// RecoverRunning devuelve a pending los jobs interrumpidos por un reinicio.
	RecoverRunning() (int64, error)
	NextPendingJob() (*models.GenerationJob, error)
	CompleteJob(id uuid.UUID, generated int) error
	FailJob(id uuid.UUID, cause string) error
}

type service struct {
	store     Store
	generator QuestionGenerator
	settings  SettingsProvider
}

func NewService(store Store, generator QuestionGenerator, settings SettingsProvider) Service {
	return &service{store: store, generator: generator, settings: settings}
}

func (s *service) CreateJob(_ context.Context, requestedBy uuid.UUID, req CreateJobRequest) (*JobResponse, error) {
	if !s.generator.IsConfigured() {
		return nil, apierr.ErrValidation(
			"La generación con IA no está configurada (faltan AI_API_URL o AI_API_KEY)", "",
		)
	}

	topicIDs := uniqueTopicIDs(req.TopicIDs)
	if len(topicIDs) == 0 {
		return nil, apierr.ErrValidation("Debes seleccionar al menos un tema", "")
	}

	found, err := s.store.TopicsExist(topicIDs)
	if err != nil {
		return nil, apierr.ErrInternal("Error al validar los temas", "")
	}
	if len(found) != len(topicIDs) {
		return nil, apierr.ErrValidation("Uno o más temas no existen", "")
	}

	quantity := req.Quantity
	if quantity == 0 {
		quantity, err = s.settings.DefaultQuantity()
		if err != nil {
			return nil, apierr.ErrInternal("Error al obtener la cantidad por defecto", "")
		}
	}

	maxPerJob, err := s.settings.MaxPerJob()
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener el máximo por trabajo", "")
	}
	if quantity > maxPerJob {
		return nil, apierr.ErrValidation(
			fmt.Sprintf("La cantidad no puede superar el máximo por trabajo (%d)", maxPerJob), "",
		)
	}

	maxPerDay, err := s.settings.MaxPerDay()
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener el tope diario", "")
	}

	job := &models.GenerationJob{
		RequestedBy: requestedBy,
		Difficulty:  req.Difficulty,
		Quantity:    quantity,
		Status:      models.GenerationJobStatusPending,
	}

	// El lock consultivo + el INSERT ocurren en la misma transacción, así que
	// dos requests concurrentes no pueden sobrepasar el tope diario.
	reserved, err := s.store.EnqueueWithBudget(job, topicIDs)
	if err != nil {
		return nil, apierr.ErrInternal("Error al encolar el trabajo", "")
	}
	if reserved+quantity > maxPerDay {
		return nil, apierr.ErrTooManyRequests(
			fmt.Sprintf(
				"Tope diario de generación alcanzado (%d/%d). Se han reservado %d preguntas hoy",
				reserved, maxPerDay, reserved,
			),
			"",
		)
	}

	topics, err := s.store.TopicsByIDs(topicIDs)
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener los temas", "")
	}
	names := make([]string, 0, len(topics))
	for _, t := range topics {
		names = append(names, t.Name)
	}

	log.Printf("📥 Trabajo de generación encolado: %d preguntas, dificultad %s, temas [%s]",
		quantity, job.Difficulty, strings.Join(names, ", "))

	out := toJobResponse(*job, names)
	return &out, nil
}

func (s *service) ListJobs(params ListJobsParams) ([]JobResponse, int64, error) {
	if params.Status != "" && !validJobStatuses[params.Status] {
		return nil, 0, apierr.ErrValidation(
			fmt.Sprintf("status '%s' no es válido. Valores permitidos: %s",
				params.Status, strings.Join(jobStatusList(), ", ")),
			"",
		)
	}
	if params.Page < common.DefaultPage {
		params.Page = common.DefaultPage
	}
	if params.PerPage <= 0 {
		params.PerPage = common.DefaultPerPage
	}

	jobs, total, err := s.store.FindPage(params)
	if err != nil {
		return nil, 0, apierr.ErrInternal("Error al obtener los trabajos", "")
	}

	names, err := s.topicNamesMap(jobs)
	if err != nil {
		return nil, 0, apierr.ErrInternal("Error al obtener los temas de los trabajos", "")
	}

	out := make([]JobResponse, 0, len(jobs))
	for _, job := range jobs {
		out = append(out, toJobResponse(job, names[job.ID]))
	}
	return out, total, nil
}

func (s *service) GetJob(id uuid.UUID) (*JobResponse, error) {
	job, err := s.store.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierr.ErrNotFound("Trabajo de generación", "")
		}
		return nil, apierr.ErrInternal("Error al obtener el trabajo", "")
	}

	names, err := s.store.TopicNames(id)
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener los temas del trabajo", "")
	}

	out := toJobResponse(*job, names)
	return &out, nil
}

func (s *service) Stats() (*StatsResponse, error) {
	budget, err := s.store.DailyBudget(startOfDay(time.Now().UTC()))
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener el consumo diario", "")
	}
	counts, err := s.store.StatusCounts()
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener el estado de los trabajos", "")
	}

	perJob, err := s.settings.MaxPerJob()
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener el máximo por trabajo", "")
	}
	perDay, err := s.settings.MaxPerDay()
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener el tope diario", "")
	}
	defaultQuantity, err := s.settings.DefaultQuantity()
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener la cantidad por defecto", "")
	}

	return &StatsResponse{
		GeneratedToday:  budget.Generated,
		ReservedToday:   budget.Reserved,
		DailyCap:        perDay,
		PerJobCap:       perJob,
		DefaultQuantity: defaultQuantity,
		Pending:         int(counts[models.GenerationJobStatusPending]),
		Running:         int(counts[models.GenerationJobStatusRunning]),
		Failed:          int(counts[models.GenerationJobStatusFailed]),
		AiConfigured:    s.generator.IsConfigured(),
	}, nil
}

func (s *service) ProcessJob(ctx context.Context, job models.GenerationJob) error {
	start := time.Now()

	topics, err := s.store.TopicsByIDs(topicIDsOf(job))
	if err != nil {
		return fmt.Errorf("error al obtener los temas: %w", err)
	}
	if len(topics) == 0 {
		return errors.New("el trabajo no tiene temas asociados")
	}
	names := make([]string, 0, len(topics))
	for _, t := range topics {
		names = append(names, t.Name)
	}

	existing, err := s.store.ExistingContentByTopics(topicIDsOf(job), maxExistingContent)
	if err != nil {
		// No es fatal: sin el histórico la IA puede repetir algún concepto.
		log.Printf("⚠️ No se pudieron obtener las preguntas existentes del trabajo %s: %v", job.ID, err)
	}

	generated, err := s.generator.GenerateQuestions(ctx, ai.BatchRequest{
		TopicNames:      names,
		Difficulty:      job.Difficulty,
		Count:           job.Quantity,
		ExistingContent: existing,
	})
	if err != nil {
		return err
	}

	saved, err := s.store.SaveQuestions(job.RequestedBy, topicIDsOf(job), generated)
	if err != nil {
		return fmt.Errorf("error al guardar las preguntas generadas: %w", err)
	}

	if err := s.store.CompleteJob(job.ID, saved); err != nil {
		return fmt.Errorf("error al marcar el trabajo como completado: %w", err)
	}

	if saved < job.Quantity {
		log.Printf("⚠️ Trabajo %s completado parcialmente: %d de %d preguntas", job.ID, saved, job.Quantity)
	} else {
		log.Printf("✅ Trabajo %s completado en %v: %d preguntas sobre [%s]",
			job.ID, time.Since(start).Round(time.Millisecond), saved, strings.Join(names, ", "))
	}

	return nil
}

func (s *service) RecoverRunning() (int64, error) {
	affected, err := s.store.RecoverRunning()
	if err != nil {
		return 0, err
	}
	if affected > 0 {
		log.Printf("🔄 %d trabajo(s) devueltos a pending tras un reinicio previo", affected)
	}
	return affected, nil
}

func (s *service) NextPendingJob() (*models.GenerationJob, error) {
	return s.store.ClaimNextPending()
}

func (s *service) CompleteJob(id uuid.UUID, generated int) error {
	return s.store.CompleteJob(id, generated)
}

func (s *service) FailJob(id uuid.UUID, cause string) error {
	return s.store.FailJob(id, cause)
}

var validJobStatuses = map[string]bool{
	models.GenerationJobStatusPending:   true,
	models.GenerationJobStatusRunning:   true,
	models.GenerationJobStatusCompleted: true,
	models.GenerationJobStatusFailed:    true,
}

func jobStatusList() []string {
	return []string{
		models.GenerationJobStatusPending,
		models.GenerationJobStatusRunning,
		models.GenerationJobStatusCompleted,
		models.GenerationJobStatusFailed,
	}
}

// topicNamesMap agrupa los nombres de temas por job para la respuesta paginada,
// evitando una consulta por trabajo.
func (s *service) topicNamesMap(jobs []models.GenerationJob) (map[uuid.UUID][]string, error) {
	ids := make([]uuid.UUID, 0, len(jobs))
	for _, job := range jobs {
		ids = append(ids, job.ID)
	}
	return s.store.TopicNamesByJobs(ids)
}

func topicIDsOf(job models.GenerationJob) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(job.JobTopics)+len(job.Topics))
	for _, t := range job.JobTopics {
		ids = append(ids, t.TopicID)
	}
	for _, t := range job.Topics {
		ids = append(ids, t.ID)
	}
	return ids
}
