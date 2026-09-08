package questions

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type topicStore interface {
	FindBySlugAndUser(slug string, createdBy *uuid.UUID) (*models.Topic, error)
}

type userStore interface {
	FindByID(id uuid.UUID) (*models.User, error)
}

type Service interface {
	List(params ListQuestionsParams) ([]QuestionListResponse, int64, error)
	GetByID(id uuid.UUID) (*QuestionResponse, error)
	Create(userID uuid.UUID, input CreateQuestionRequest) (*QuestionResponse, error)
	Update(id uuid.UUID, userID uuid.UUID, input UpdateQuestionRequest) (*QuestionResponse, error)
	Delete(id uuid.UUID, userID uuid.UUID) error
	Import(userID uuid.UUID, r io.Reader) (*ImportResult, error)
	GetImportQuota(userID uuid.UUID) (*ImportQuota, error)
	GetAiQuota(userID uuid.UUID) (*AiQuota, error)
	Stats(isAdmin bool, userID uuid.UUID) (*QuestionStats, error)
}

type questionService struct {
	store      Store
	topicStore topicStore
	userStore  userStore
}

func NewService(store Store, topicStore topicStore, userStore userStore) Service {
	return &questionService{store: store, topicStore: topicStore, userStore: userStore}
}

func (s *questionService) List(params ListQuestionsParams) ([]QuestionListResponse, int64, error) {
	questions, total, err := s.store.FindPage(params)
	if err != nil {
		return nil, 0, apierr.ErrInternal("Error al listar las preguntas", "")
	}

	result := make([]QuestionListResponse, len(questions))
	for i, q := range questions {
		result[i] = ToQuestionListResponse(q)
	}
	return result, total, nil
}

func (s *questionService) GetByID(id uuid.UUID) (*QuestionResponse, error) {
	question, err := s.store.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Pregunta", "")
		}
		return nil, apierr.ErrInternal("Error al obtener la pregunta", "")
	}

	resp := ToQuestionResponse(*question)
	return &resp, nil
}

func (s *questionService) Create(userID uuid.UUID, input CreateQuestionRequest) (*QuestionResponse, error) {
	if input.Type == "single_choice" || input.Type == "multiple_choice" {
		if len(input.Options) == 0 {
			return nil, apierr.ErrValidation("Se requieren opciones para preguntas de tipo choice", "")
		}
	}
	if input.Type == "code_completion" && input.Language == "" {
		return nil, apierr.ErrValidation("Se requiere el lenguaje para preguntas code_completion", "")
	}

	question := &models.Question{
		UserID:      userID,
		Type:        input.Type,
		Content:     input.Content,
		Explanation: input.Explanation,
		Difficulty:  input.Difficulty,
		Source:      input.Source,
	}

	if question.Source == "" {
		question.Source = "manual"
	}

	if input.Type == "code_completion" {
		question.CodeChallenge = &models.CodeChallenge{
			StarterCode:    input.StarterCode,
			ExpectedOutput: input.ExpectedOutput,
			Language:       input.Language,
			TestCasesJSON:  input.TestCasesJSON,
		}
	} else {
		for _, opt := range input.Options {
			question.Options = append(question.Options, models.QuestionOption{
				Content:   opt.Content,
				IsCorrect: opt.IsCorrect,
			})
		}
	}

	if err := s.store.Create(question); err != nil {
		return nil, apierr.ErrInternal("Error al crear la pregunta", "")
	}

	if err := s.store.AddQuestionTopics(question.ID, input.TopicIDs); err != nil {
		return nil, apierr.ErrInternal("Error al asociar los temas", "")
	}

	created, err := s.store.FindByID(question.ID)
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener la pregunta creada", "")
	}
	question = created

	resp := ToQuestionResponse(*question)
	return &resp, nil
}

func (s *questionService) Update(id uuid.UUID, userID uuid.UUID, input UpdateQuestionRequest) (*QuestionResponse, error) {
	question, err := s.store.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Pregunta", "")
		}
		return nil, apierr.ErrInternal("Error al obtener la pregunta", "")
	}

	if !canModify(question, userID) {
		return nil, apierr.ErrForbidden("No tienes permiso para modificar esta pregunta", "")
	}

	if input.Content != "" {
		question.Content = input.Content
	}
	if input.Explanation != "" {
		question.Explanation = input.Explanation
	}
	if input.Difficulty != "" {
		question.Difficulty = input.Difficulty
	}

	if input.TopicIDs != nil {
		if err := s.store.ReplaceQuestionTopics(question.ID, input.TopicIDs); err != nil {
			return nil, apierr.ErrInternal("Error al actualizar los temas", "")
		}
		question.Topics = nil
	}

	if input.Options != nil {
		options := make([]models.QuestionOption, len(input.Options))
		for i, opt := range input.Options {
			options[i] = models.QuestionOption{
				Content:   opt.Content,
				IsCorrect: opt.IsCorrect,
			}
		}
		if err := s.store.ReplaceQuestionOptions(question.ID, options); err != nil {
			return nil, apierr.ErrInternal("Error al actualizar las opciones", "")
		}
		question.Options = nil
	}

	if input.StarterCode != "" || input.Language != "" {
		if question.CodeChallenge == nil {
			question.CodeChallenge = &models.CodeChallenge{
				QuestionID: question.ID,
			}
		}
		if input.StarterCode != "" {
			question.CodeChallenge.StarterCode = input.StarterCode
		}
		if input.ExpectedOutput != "" {
			question.CodeChallenge.ExpectedOutput = input.ExpectedOutput
		}
		if input.Language != "" {
			question.CodeChallenge.Language = input.Language
		}
		if input.TestCasesJSON != "" {
			question.CodeChallenge.TestCasesJSON = input.TestCasesJSON
		}
	}

	if err := s.store.Update(question); err != nil {
		return nil, apierr.ErrInternal("Error al actualizar la pregunta", "")
	}

	question, err = s.store.FindByID(question.ID)
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener la pregunta actualizada", "")
	}

	resp := ToQuestionResponse(*question)
	return &resp, nil
}

func canModify(question *models.Question, userID uuid.UUID) bool {
	if question.Source != "manual" && question.Source != "imported" {
		return false
	}
	return question.UserID == userID
}

func (s *questionService) Delete(id uuid.UUID, userID uuid.UUID) error {
	question, err := s.store.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return apierr.ErrNotFound("Pregunta", "")
		}
		return apierr.ErrInternal("Error al obtener la pregunta", "")
	}

	if !canModify(question, userID) {
		return apierr.ErrForbidden("No tienes permiso para eliminar esta pregunta", "")
	}

	if err := s.store.Delete(id); err != nil {
		return apierr.ErrInternal("Error al eliminar la pregunta", "")
	}

	return nil
}

func startOfDayUTC(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func (s *questionService) Stats(isAdmin bool, userID uuid.UUID) (*QuestionStats, error) {
	stats, err := s.store.Stats(isAdmin, userID)
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener las estadísticas", "")
	}
	return stats, nil
}

func (s *questionService) GetImportQuota(userID uuid.UUID) (*ImportQuota, error) {
	user, err := s.userStore.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Usuario", "")
		}
		return nil, apierr.ErrInternal("Error al obtener el usuario", "")
	}

	since := startOfDayUTC(time.Now())
	used, err := s.store.CountImportedSince(userID, since)
	if err != nil {
		return nil, apierr.ErrInternal("Error al calcular el cupo diario", "")
	}

	remaining := user.DailyImportLimit - int(used)
	if remaining < 0 {
		remaining = 0
	}

	return &ImportQuota{
		DailyLimit: user.DailyImportLimit,
		UsedToday:  int(used),
		Remaining:  remaining,
	}, nil
}

func (s *questionService) GetAiQuota(userID uuid.UUID) (*AiQuota, error) {
	user, err := s.userStore.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Usuario", "")
		}
		return nil, apierr.ErrInternal("Error al obtener el usuario", "")
	}

	since := startOfDayUTC(time.Now())
	used, err := s.store.CountAiGeneratedSince(userID, since)
	if err != nil {
		return nil, apierr.ErrInternal("Error al calcular el cupo diario de IA", "")
	}

	remaining := user.DailyAiLimit - int(used)
	if remaining < 0 {
		remaining = 0
	}

	return &AiQuota{
		DailyLimit: user.DailyAiLimit,
		UsedToday:  int(used),
		Remaining:  remaining,
	}, nil
}

func (s *questionService) Import(userID uuid.UUID, r io.Reader) (*ImportResult, error) {
	user, err := s.userStore.FindByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Usuario", "")
		}
		return nil, apierr.ErrInternal("Error al obtener el usuario", "")
	}

	isAdmin := user.IsAdmin
	remaining := maxQuestionsPerFile
	if !isAdmin {
		since := startOfDayUTC(time.Now())
		used, err := s.store.CountImportedSince(userID, since)
		if err != nil {
			return nil, apierr.ErrInternal("Error al calcular el cupo diario", "")
		}
		remaining = user.DailyImportLimit - int(used)
		if remaining <= 0 {
			return nil, apierr.ErrTooManyRequests(
				fmt.Sprintf("Has alcanzado tu límite diario de %d preguntas importadas", user.DailyImportLimit), "")
		}
	}

	rows, err := parseCSV(r)
	if err != nil {
		return nil, apierr.ErrValidation(err.Error(), "")
	}

	if len(rows) == 0 {
		return nil, apierr.ErrValidation("El archivo no contiene preguntas", "")
	}

	totalRows := len(rows)
	var overflowErrors []ImportRowError

	if len(rows) > maxQuestionsPerFile {
		for _, row := range rows[maxQuestionsPerFile:] {
			overflowErrors = append(overflowErrors, ImportRowError{
				Row:    row.RowNum,
				Reason: fmt.Sprintf("Supera el límite de %d preguntas por archivo", maxQuestionsPerFile),
			})
		}
		rows = rows[:maxQuestionsPerFile]
	}

	if len(rows) > remaining {
		for _, row := range rows[remaining:] {
			overflowErrors = append(overflowErrors, ImportRowError{
				Row:    row.RowNum,
				Reason: "Supera el cupo diario restante",
			})
		}
		rows = rows[:remaining]
	}

	topicCache := make(map[string]uuid.UUID)
	validQuestions := make([]*models.Question, 0, len(rows))
	rowErrors := overflowErrors

	for _, row := range rows {
		validated, vErr := validateImportRow(row)
		if vErr != nil {
			rowErrors = append(rowErrors, ImportRowError{Row: row.RowNum, Reason: vErr.Error()})
			continue
		}

		topicIDs, tErr := s.resolveTopicIDs(validated.TopicSlugs, userID, isAdmin, topicCache)
		if tErr != nil {
			rowErrors = append(rowErrors, ImportRowError{Row: row.RowNum, Reason: tErr.Error()})
			continue
		}

		q := &models.Question{
			UserID:      userID,
			Type:        validated.Type,
			Content:     strings.TrimSpace(row.Content),
			Explanation: strings.TrimSpace(row.Explanation),
			Difficulty:  validated.Difficulty,
			IsPublic:    isAdmin,
			Source:      "imported",
		}

		for _, o := range validated.Opts {
			q.Options = append(q.Options, models.QuestionOption{
				Content:   o.Content,
				IsCorrect: o.IsCorrect,
			})
		}

		topics := make([]models.Topic, len(topicIDs))
		for i, tid := range topicIDs {
			topics[i] = models.Topic{ID: tid}
		}
		q.Topics = topics

		validQuestions = append(validQuestions, q)
	}

	result := &ImportResult{
		Total:  totalRows,
		Failed: len(rowErrors),
		Errors: rowErrors,
	}

	if len(validQuestions) > 0 {
		if err := s.store.BulkCreate(validQuestions); err != nil {
			return nil, apierr.ErrInternal("Error al importar las preguntas", "")
		}
		result.Imported = len(validQuestions)
	}

	if result.Errors == nil {
		result.Errors = []ImportRowError{}
	}

	return result, nil
}

func (s *questionService) resolveTopicIDs(slugs []string, userID uuid.UUID, isAdmin bool, cache map[string]uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	for _, slug := range slugs {
		if cached, ok := cache[slug]; ok {
			ids = append(ids, cached)
			continue
		}
		tid, err := s.findTopic(slug, userID, isAdmin, cache)
		if err != nil {
			return nil, err
		}
		ids = append(ids, tid)
	}
	return ids, nil
}

func (s *questionService) findTopic(slug string, userID uuid.UUID, isAdmin bool, cache map[string]uuid.UUID) (uuid.UUID, error) {
	if tid, ok := cache[slug]; ok {
		return tid, nil
	}

	if isAdmin {
		t, err := s.topicStore.FindBySlugAndUser(slug, nil)
		if err == nil {
			cache[slug] = t.ID
			return t.ID, nil
		}
		if err != gorm.ErrRecordNotFound {
			return uuid.Nil, err
		}
		return uuid.Nil, fmt.Errorf("tema no encontrado: %s", slug)
	}

	if t, err := s.topicStore.FindBySlugAndUser(slug, &userID); err == nil {
		cache[slug] = t.ID
		return t.ID, nil
	} else if err != gorm.ErrRecordNotFound {
		return uuid.Nil, err
	}

	if t, err := s.topicStore.FindBySlugAndUser(slug, nil); err == nil {
		cache[slug] = t.ID
		return t.ID, nil
	} else if err != gorm.ErrRecordNotFound {
		return uuid.Nil, err
	}

	return uuid.Nil, fmt.Errorf("tema no encontrado: %s", slug)
}
