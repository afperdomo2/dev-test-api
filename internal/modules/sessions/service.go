package sessions

import (
	"encoding/json"
	"log"
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/modules/progress"
	"github.com/felipe/dev-test-api/internal/services/ai"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	List(userID uuid.UUID, params ListSessionsParams) ([]SessionListResponse, int64, error)
	GetByID(sessionID uuid.UUID) (*SessionDetailResponse, error)
	Create(userID uuid.UUID, input CreateSessionRequest) (*SessionResponse, error)
	Finish(sessionID uuid.UUID) (*SessionResponse, error)
	NextQuestion(sessionID uuid.UUID) (*NextQuestionResponse, error)
	Answer(sessionID, userID uuid.UUID, input AnswerRequest) (*SessionAnswerResponse, error)
	Summary(sessionID uuid.UUID) (*SessionSummaryResponse, error)
}

type sessionService struct {
	store           Store
	progressService progress.Service
	aiGenerator     *ai.Generator
}

func NewService(store Store, progressService progress.Service, aiGenerator *ai.Generator) Service {
	return &sessionService{store: store, progressService: progressService, aiGenerator: aiGenerator}
}

func (s *sessionService) List(userID uuid.UUID, params ListSessionsParams) ([]SessionListResponse, int64, error) {
	sessions, total, err := s.store.FindPage(userID, params)
	if err != nil {
		return nil, 0, apierr.ErrInternal("Error al listar las sesiones", "")
	}

	result := make([]SessionListResponse, len(sessions))
	for i, sess := range sessions {
		result[i] = ToSessionListResponse(sess)
	}
	return result, total, nil
}

func (s *sessionService) GetByID(sessionID uuid.UUID) (*SessionDetailResponse, error) {
	sess, err := s.store.FindByID(sessionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Sesion", "")
		}
		return nil, apierr.ErrInternal("Error al obtener la sesion", "")
	}

	answers := make([]SessionAnswerDetailResponse, len(sess.Answers))
	for i, a := range sess.Answers {
		answers[i] = ToAnswerDetailResponse(a)
	}

	sessionResp := ToSessionResponse(*sess)
	return &SessionDetailResponse{
		Session: sessionResp,
		Answers: answers,
	}, nil
}

func (s *sessionService) Create(userID uuid.UUID, input CreateSessionRequest) (*SessionResponse, error) {
	session := &models.Session{
		UserID:        userID,
		Name:          input.Name,
		Status:        "in_progress",
		Mode:          input.Mode,
		Difficulty:    input.Difficulty,
		QuestionLimit: input.QuestionLimit,
	}

	if err := s.store.Create(session); err != nil {
		return nil, apierr.ErrInternal("Error al crear la sesion", "")
	}

	if err := s.store.AddSessionTopics(session.ID, input.TopicIDs); err != nil {
		return nil, apierr.ErrInternal("Error al asociar los temas", "")
	}

	session, err := s.store.FindByID(session.ID)
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener la sesion creada", "")
	}

	if session.Mode == "generate" {
		go s.generateBatch(session, 2)
	}

	resp := ToSessionResponse(*session)
	return &resp, nil
}

func (s *sessionService) Finish(sessionID uuid.UUID) (*SessionResponse, error) {
	sess, err := s.store.FindByID(sessionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Sesion", "")
		}
		return nil, apierr.ErrInternal("Error al obtener la sesion", "")
	}

	if sess.Status == "completed" {
		return nil, apierr.ErrConflict("Session Already Completed", "La sesion ya fue completada", "")
	}

	now := time.Now()
	sess.Status = "completed"
	sess.FinishedAt = &now

	if err := s.store.Update(sess); err != nil {
		return nil, apierr.ErrInternal("Error al finalizar la sesion", "")
	}

	sess, err = s.store.FindByID(sessionID)
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener la sesion actualizada", "")
	}

	resp := ToSessionResponse(*sess)
	return &resp, nil
}

func (s *sessionService) NextQuestion(sessionID uuid.UUID) (*NextQuestionResponse, error) {
	sess, err := s.store.FindByID(sessionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Sesion", "")
		}
		return nil, apierr.ErrInternal("Error al obtener la sesion", "")
	}

	if sess.Status == "completed" {
		return nil, apierr.ErrConflict("Session Completed", "La sesion ya fue completada", "")
	}

	answeredIDs, err := s.store.FindAnsweredQuestionIDs(sessionID)
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener las preguntas respondidas", "")
	}

	if sess.QuestionLimit != nil && len(answeredIDs) >= *sess.QuestionLimit {
		return nil, apierr.ErrNotFound("Pregunta", "Has alcanzado el limite de preguntas de esta sesion")
	}

	topicIDs := make([]uuid.UUID, len(sess.Topics))
	for i, t := range sess.Topics {
		topicIDs[i] = t.ID
	}

	question, err := s.store.FindNextQuestion(topicIDs, answeredIDs, sess.Difficulty, sess.Mode, sess.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if sess.Mode == "generate" {
				if err := s.aiGenerator.GenerateQuestion(sess); err != nil {
					log.Printf("⚠️ Error generando pregunta para sesión %s: %v", sess.ID, err)
					return nil, apierr.ErrNotFound("Pregunta", "No hay mas preguntas disponibles para esta sesion")
				}
				question, err = s.store.FindNextQuestion(topicIDs, answeredIDs, sess.Difficulty, sess.Mode, sess.UserID)
				if err != nil {
					return nil, apierr.ErrNotFound("Pregunta", "No hay mas preguntas disponibles para esta sesion")
				}
			} else {
				return nil, apierr.ErrNotFound("Pregunta", "No hay mas preguntas disponibles para esta sesion")
			}
		} else {
			return nil, apierr.ErrInternal("Error al obtener la siguiente pregunta", "")
		}
	}

	resp := NextQuestionResponse{
		Question: toNextQuestionItem(*question),
	}
	return &resp, nil
}

func (s *sessionService) Answer(sessionID, userID uuid.UUID, input AnswerRequest) (*SessionAnswerResponse, error) {
	sess, err := s.store.FindByID(sessionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Sesion", "")
		}
		return nil, apierr.ErrInternal("Error al obtener la sesion", "")
	}

	if sess.Status == "completed" {
		return nil, apierr.ErrConflict("Session Completed", "La sesion ya fue completada", "")
	}

	selectedJSON := "[]"
	if len(input.SelectedOptions) > 0 {
		b, _ := json.Marshal(input.SelectedOptions)
		selectedJSON = string(b)
	}

	isCorrect := false
	question, err := s.store.FindQuestionByID(input.QuestionID)
	if err == nil {
		isCorrect = evaluateCorrectness(question, input.SelectedOptions)
	}

	answer := &models.SessionAnswer{
		SessionID:       sessionID,
		UserID:          userID,
		QuestionID:      input.QuestionID,
		AnswerText:      input.AnswerText,
		SelectedOptions: selectedJSON,
		IsCorrect:       isCorrect,
		ResponseTimeMs:  input.ResponseTimeMs,
	}

	if err := s.store.CreateAnswer(answer); err != nil {
		return nil, apierr.ErrInternal("Error al guardar la respuesta", "")
	}

	s.progressService.Answer(userID, input.QuestionID, isCorrect)

	if question != nil {
		answer.Question = question
	}

	if sess.Mode == "generate" {
		answeredIDs, err := s.store.FindAnsweredQuestionIDs(sessionID)
		if err == nil {
			topicIDs := make([]uuid.UUID, len(sess.Topics))
			for i, t := range sess.Topics {
				topicIDs[i] = t.ID
			}
			available, _ := s.store.CountAvailableQuestions(topicIDs, answeredIDs, sess.Difficulty, sess.Mode, sess.UserID)
			if available < 2 {
				go s.generateBatch(sess, 1)
			}
		}
	}

	resp := ToAnswerResponse(*answer)
	return &resp, nil
}

func evaluateCorrectness(question *models.Question, selectedOptions []uuid.UUID) bool {
	switch question.Type {
	case "single_choice":
		for _, opt := range question.Options {
			if opt.IsCorrect {
				return len(selectedOptions) == 1 && selectedOptions[0] == opt.ID
			}
		}
		return false
	case "multiple_choice":
		correctIDs := make(map[uuid.UUID]bool)
		for _, opt := range question.Options {
			if opt.IsCorrect {
				correctIDs[opt.ID] = true
			}
		}
		if len(selectedOptions) != len(correctIDs) {
			return false
		}
		for _, id := range selectedOptions {
			if !correctIDs[id] {
				return false
			}
		}
		return len(correctIDs) > 0
	default:
		return false
	}
}

func (s *sessionService) generateBatch(sess *models.Session, count int) {
	for i := 0; i < count; i++ {
		if sess.QuestionLimit != nil && sess.QuestionsGenerated >= *sess.QuestionLimit {
			break
		}
		s.aiGenerator.GenerateQuestion(sess)
	}
}

func (s *sessionService) Summary(sessionID uuid.UUID) (*SessionSummaryResponse, error) {
	data, err := s.store.FindSummary(sessionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Sesion", "")
		}
		return nil, apierr.ErrInternal("Error al obtener el resumen de la sesion", "")
	}

	return &SessionSummaryResponse{
		AnswerCount:        int(data.AnswerCount),
		QuestionsGenerated: data.QuestionsGenerated,
		Status:             data.Status,
		QuestionLimit:      data.QuestionLimit,
	}, nil
}
