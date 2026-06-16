package sessions

import (
	"encoding/json"
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/modules/progress"
	"github.com/felipe/dev-test-api/internal/modules/questions"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	List(userID uuid.UUID, page, perPage int, sortBy, sortOrder string) ([]SessionResponse, int64, error)
	GetByID(sessionID uuid.UUID) (*SessionDetailResponse, error)
	Create(userID uuid.UUID, input CreateSessionRequest) (*SessionResponse, error)
	Finish(sessionID uuid.UUID) (*SessionResponse, error)
	NextQuestion(sessionID uuid.UUID) (*NextQuestionResponse, error)
	Answer(sessionID, userID uuid.UUID, input AnswerRequest) (*SessionAnswerResponse, error)
}

type sessionService struct {
	store           Store
	progressService progress.Service
}

func NewService(store Store, progressService progress.Service) Service {
	return &sessionService{store: store, progressService: progressService}
}

func (s *sessionService) List(userID uuid.UUID, page, perPage int, sortBy, sortOrder string) ([]SessionResponse, int64, error) {
	sessions, total, err := s.store.FindPage(userID, page, perPage, sortBy, sortOrder)
	if err != nil {
		return nil, 0, apierr.ErrInternal("Error al listar las sesiones", "")
	}

	result := make([]SessionResponse, len(sessions))
	for i, sess := range sessions {
		result[i] = toSessionResponse(sess)
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

	answers := make([]SessionAnswerResponse, len(sess.Answers))
	for i, a := range sess.Answers {
		answers[i] = toAnswerResponse(a)
	}

	sessionResp := toSessionResponse(*sess)
	return &SessionDetailResponse{
		Session: sessionResp,
		Answers: answers,
	}, nil
}

func (s *sessionService) Create(userID uuid.UUID, input CreateSessionRequest) (*SessionResponse, error) {
	session := &models.Session{
		UserID:     userID,
		Name:       input.Name,
		Status:     "in_progress",
		Mode:       input.Mode,
		Difficulty: input.Difficulty,
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

	resp := toSessionResponse(*session)
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

	resp := toSessionResponse(*sess)
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

	answeredIDs, err := s.store.FindAnsweredQuestionIDs(sessionID)
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener las preguntas respondidas", "")
	}

	topicIDs := make([]uuid.UUID, len(sess.Topics))
	for i, t := range sess.Topics {
		topicIDs[i] = t.ID
	}

	question, err := s.store.FindNextQuestion(topicIDs, answeredIDs, sess.Difficulty, sess.Mode, sess.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Pregunta", "No hay mas preguntas disponibles para esta sesion")
		}
		return nil, apierr.ErrInternal("Error al obtener la siguiente pregunta", "")
	}

	resp := NextQuestionResponse{
		Question: questions.ToQuestionResponse(*question),
	}
	return &resp, nil
}

func (s *sessionService) Answer(sessionID, userID uuid.UUID, input AnswerRequest) (*SessionAnswerResponse, error) {
	selectedJSON := "[]"
	if len(input.SelectedOptions) > 0 {
		b, _ := json.Marshal(input.SelectedOptions)
		selectedJSON = string(b)
	}

	answer := &models.SessionAnswer{
		SessionID:       sessionID,
		UserID:          userID,
		QuestionID:      input.QuestionID,
		AnswerText:      input.AnswerText,
		SelectedOptions: selectedJSON,
		IsCorrect:       input.IsCorrect,
		ResponseTimeMs:  input.ResponseTimeMs,
	}

	if err := s.store.CreateAnswer(answer); err != nil {
		return nil, apierr.ErrInternal("Error al guardar la respuesta", "")
	}

	s.progressService.Answer(userID, input.QuestionID, input.IsCorrect)

	resp := toAnswerResponse(*answer)
	return &resp, nil
}
