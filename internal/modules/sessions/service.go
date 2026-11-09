package sessions

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/modules/gamification"
	"github.com/felipe/dev-test-api/internal/modules/progress"
	"github.com/felipe/dev-test-api/internal/modules/tips"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	List(userID uuid.UUID, params ListSessionsParams) ([]SessionListResponse, int64, error)
	GetByID(sessionID, userID uuid.UUID) (*SessionDetailResponse, error)
	Create(userID uuid.UUID, input CreateSessionRequest) (*SessionResponse, error)
	CreateDaily(userID uuid.UUID) (*DailySessionResponse, error)
	Finish(sessionID, userID uuid.UUID) (*SessionResponse, error)
	NextQuestion(sessionID, userID uuid.UUID) (*NextQuestionResponse, error)
	Answer(sessionID, userID uuid.UUID, input AnswerRequest) (*SessionAnswerResponse, error)
	Summary(sessionID, userID uuid.UUID) (*SessionSummaryResponse, error)
	Delete(sessionID, userID uuid.UUID) error
}

type sessionService struct {
	store               Store
	progressService     progress.Service
	gamificationService gamification.Service
	tipService          tips.Service
}

func NewService(store Store, progressService progress.Service, gamificationService gamification.Service, tipService tips.Service) Service {
	return &sessionService{store: store, progressService: progressService, gamificationService: gamificationService, tipService: tipService}
}

// findOwnedSession carga la sesión y comprueba que pertenece al usuario.
//
// Si la sesión no existe o es de otro usuario devuelve el mismo ErrNotFound:
// así un tercero no puede distinguir "no existe" de "es de otro" (D-AZ-1).
// Toda operación por id pasa por aquí (D-AZ-2).
func findOwnedSession(store Store, sessionID, userID uuid.UUID) (*models.Session, error) {
	sess, err := store.FindByID(sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierr.ErrNotFound("Sesion", "")
		}
		return nil, apierr.ErrInternal("Error al obtener la sesion", "")
	}

	if sess.UserID != userID {
		return nil, apierr.ErrNotFound("Sesion", "")
	}

	return sess, nil
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

func (s *sessionService) GetByID(sessionID, userID uuid.UUID) (*SessionDetailResponse, error) {
	sess, err := findOwnedSession(s.store, sessionID, userID)
	if err != nil {
		return nil, err
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
	if input.Mode == "review" {
		available, err := s.store.CountAvailableQuestions(input.TopicIDs, nil, input.Difficulty, "review", userID)
		if err != nil {
			return nil, apierr.ErrInternal("Error al validar preguntas disponibles", "")
		}
		if available == 0 {
			return nil, apierr.ErrValidation("No hay preguntas disponibles para estos temas y dificultad.", "")
		}
		if input.QuestionLimit != nil && available < int64(*input.QuestionLimit) {
			msg := fmt.Sprintf("Solo hay %d preguntas disponibles para estos temas, pero pediste %d. Reduce el límite o elige más temas.", available, *input.QuestionLimit)
			return nil, apierr.ErrValidation(msg, "")
		}
	}

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

	resp := ToSessionResponse(*session)
	return &resp, nil
}

type dailySessionData struct {
	dueCount      int64
	warmupCount   int64
	practiceCount int64
	tip           *tips.TipResponse
	errs          []error
}

func (s *sessionService) gatherDailyData(userID uuid.UUID) dailySessionData {
	var (
		wg  sync.WaitGroup
		mu  sync.Mutex
		res dailySessionData
	)

	collect := func(fn func() error) {
		defer wg.Done()
		if err := fn(); err != nil {
			mu.Lock()
			res.errs = append(res.errs, err)
			mu.Unlock()
		}
	}

	wg.Add(4)
	go collect(func() error {
		n, err := s.store.CountDueReviews(userID)
		mu.Lock()
		res.dueCount = n
		mu.Unlock()
		return err
	})
	go collect(func() error {
		n, err := s.store.CountWarmupCandidates(userID)
		mu.Lock()
		res.warmupCount = n
		mu.Unlock()
		return err
	})
	go collect(func() error {
		n, err := s.store.CountPracticeCandidates(userID)
		mu.Lock()
		res.practiceCount = n
		mu.Unlock()
		return err
	})
	go collect(func() error {
		tip, err := s.tipService.GetTipOfTheDay(userID)
		if err != nil {
			var apiErr *apierr.APIError
			if errors.As(err, &apiErr) && apiErr.Status == http.StatusNotFound {
				return nil
			}
			return err
		}
		mu.Lock()
		res.tip = tip
		mu.Unlock()
		return nil
	})
	wg.Wait()

	return res
}

func (s *sessionService) CreateDaily(userID uuid.UUID) (*DailySessionResponse, error) {
	data := s.gatherDailyData(userID)
	if len(data.errs) > 0 {
		return nil, apierr.ErrInternal("Error al preparar la sesion diaria", "")
	}

	// Reutiliza la sesión diaria de hoy si sigue en progreso: así pulsar
	// "Repasar hoy" varias veces no acumula sesiones abiertas.
	if existing := s.findTodayDailySession(userID); existing != nil {
		return &DailySessionResponse{
			Session:     ToSessionResponse(*existing),
			Tip:         s.loadSessionTip(existing),
			DueCount:    data.dueCount,
			WarmupCount: data.warmupCount,
		}, nil
	}

	// Sin vencidas y sin banco visible no hay nada que repasar. No se crea una
	// sesión fantasma que se quedaría "en progreso" para siempre (ver D-52).
	if data.dueCount == 0 && data.practiceCount == 0 {
		return nil, apierr.ErrConflict(
			"No Questions For Daily Session",
			"No hay preguntas para repasar hoy. Crea una sesión de estudio.",
			"",
		)
	}

	session := &models.Session{
		UserID:     userID,
		Name:       "Repaso diario",
		Status:     "in_progress",
		Mode:       "daily",
		Difficulty: "intermediate",
	}
	if data.tip != nil {
		session.TipID = &data.tip.ID
	}

	if err := s.store.Create(session); err != nil {
		return nil, apierr.ErrInternal("Error al crear la sesion diaria", "")
	}

	session, err := s.store.FindByID(session.ID)
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener la sesion diaria", "")
	}

	sessResp := ToSessionResponse(*session)
	return &DailySessionResponse{
		Session:     sessResp,
		Tip:         data.tip,
		DueCount:    data.dueCount,
		WarmupCount: data.warmupCount,
	}, nil
}

// findTodayDailySession devuelve la sesión diaria de hoy en progreso, o nil si
// no existe (o si la consulta falla: el peor caso es crear una nueva).
func (s *sessionService) findTodayDailySession(userID uuid.UUID) *models.Session {
	startOfDay := time.Now().UTC().Truncate(24 * time.Hour)
	sess, err := s.store.FindTodayInProgressDaily(userID, startOfDay)
	if err != nil {
		return nil
	}
	return sess
}

// loadSessionTip recupera el tip asociado a una sesión reutilizada. Si no hay
// tip guardado o no se puede leer, la sesión funciona igual sin tarjeta.
func (s *sessionService) loadSessionTip(sess *models.Session) *tips.TipResponse {
	if sess.TipID == nil || s.tipService == nil {
		return nil
	}
	tip, err := s.tipService.GetByID(*sess.TipID)
	if err != nil {
		return nil
	}
	return tip
}

func (s *sessionService) Finish(sessionID, userID uuid.UUID) (*SessionResponse, error) {
	sess, err := findOwnedSession(s.store, sessionID, userID)
	if err != nil {
		return nil, err
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
	// Finalizar sin responder no otorga recompensa: una sesión vacía no debe
	// servir para granjear XP ni racha (ver D-52).
	if s.gamificationService != nil && len(sess.Answers) > 0 {
		reward, gErr := s.gamificationService.CompleteSession(sess.UserID)
		if gErr != nil {
			return nil, gErr
		}
		resp.Reward = reward
	}
	return &resp, nil
}

func (s *sessionService) NextQuestion(sessionID, userID uuid.UUID) (*NextQuestionResponse, error) {
	sess, err := findOwnedSession(s.store, sessionID, userID)
	if err != nil {
		return nil, err
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

	if sess.Mode == "daily" {
		question, err := s.store.FindNextDailyQuestion(sess.UserID, answeredIDs)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, apierr.ErrNotFound("Pregunta", "No hay mas preguntas disponibles para esta sesion")
			}
			return nil, apierr.ErrInternal("Error al obtener la siguiente pregunta", "")
		}

		resp := NextQuestionResponse{
			Question: toNextQuestionItem(*question),
		}
		return &resp, nil
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
		Question: toNextQuestionItem(*question),
	}
	return &resp, nil
}

func (s *sessionService) Answer(sessionID, userID uuid.UUID, input AnswerRequest) (*SessionAnswerResponse, error) {
	sess, err := findOwnedSession(s.store, sessionID, userID)
	if err != nil {
		return nil, err
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
		ClientAnswerID:  input.ClientAnswerID,
	}

	created, err := s.store.CreateAnswer(answer)
	if err != nil {
		return nil, apierr.ErrInternal("Error al guardar la respuesta", "")
	}

	if !created {
		// Reintento con el mismo clientAnswerId: devolver la respuesta ya
		// guardada y NO volver a aplicar progreso/gamificación (si no, doble
		// XP y doble SRS).
		return s.answerAlreadyStored(sessionID, input.ClientAnswerID)
	}

	var reward *gamification.Reward
	if s.progressService != nil {
		progressResp, pErr := s.progressService.Answer(userID, input.QuestionID, isCorrect)
		if pErr != nil {
			log.Printf("⚠️ Error al guardar el progreso/recompensa para la respuesta %s: %v", answer.ID, pErr)
			progressResp = nil
		}
		if progressResp != nil {
			reward = progressResp.Reward
		}
	}

	if question != nil {
		answer.Question = question
	}

	resp := ToAnswerResponse(*answer)
	resp.Reward = reward
	return &resp, nil
}

// answerAlreadyStored recupera la respuesta que ya existía para el mismo
// (session_id, client_answer_id) y la devuelve sin efectos secundarios.
func (s *sessionService) answerAlreadyStored(sessionID uuid.UUID, clientAnswerID *uuid.UUID) (*SessionAnswerResponse, error) {
	if clientAnswerID == nil {
		// `created=false` solo puede venir de un conflicto, y sin client id no
		// hay clave de idempotencia que recuperar.
		return nil, apierr.ErrInternal("Error al guardar la respuesta", "")
	}

	stored, err := s.store.FindByClientAnswerID(sessionID, *clientAnswerID)
	if err != nil {
		return nil, apierr.ErrInternal("Error al recuperar la respuesta guardada", "")
	}

	if question, qErr := s.store.FindQuestionByID(stored.QuestionID); qErr == nil {
		stored.Question = question
	}

	resp := ToAnswerResponse(*stored)
	return &resp, nil
}

func evaluateCorrectness(question *models.Question, selectedOptions []uuid.UUID) bool {
	switch question.Type {
	case "single_choice", "true_false":
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

func (s *sessionService) Delete(sessionID, userID uuid.UUID) error {
	sess, err := findOwnedSession(s.store, sessionID, userID)
	if err != nil {
		return err
	}

	if len(sess.Answers) > 0 {
		return apierr.ErrConflict("Session Has Answers", "No puedes eliminar una sesión que tiene respuestas. Finalízala en su lugar.", "")
	}

	if sess.CreatedAt.Before(time.Now().Add(-24 * time.Hour)) {
		return apierr.ErrForbidden("Solo se pueden eliminar sesiones creadas en las últimas 24 horas", "")
	}

	return s.store.Delete(sessionID)
}

func (s *sessionService) Summary(sessionID, userID uuid.UUID) (*SessionSummaryResponse, error) {
	// FindSummary solo agrega contadores y no filtra por usuario: el dueño se
	// comprueba antes con la misma regla que el resto (D-AZ-2).
	if _, err := findOwnedSession(s.store, sessionID, userID); err != nil {
		return nil, err
	}

	data, err := s.store.FindSummary(sessionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrNotFound("Sesion", "")
		}
		return nil, apierr.ErrInternal("Error al obtener el resumen de la sesion", "")
	}

	return &SessionSummaryResponse{
		AnswerCount:   int(data.AnswerCount),
		Status:        data.Status,
		QuestionLimit: data.QuestionLimit,
	}, nil
}
