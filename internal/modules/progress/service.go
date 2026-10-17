package progress

import (
	"math"
	"sort"
	"time"

	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/modules/gamification"
	"github.com/felipe/dev-test-api/internal/modules/questions"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	// D-08: un tema es "débil" por debajo del 60% de acierto y con historial
	// suficiente. La definición vive SOLO aquí; insights y badges la consumen.
	weakAccuracyThreshold = 0.6
	minAttempts           = 5

	// Factor de redondeo de accuracy (4 decimales): evita que el JSON arrastre
	// ruido de coma flotante y hace estables los tests.
	accuracyRounding = 10000
)

type Service interface {
	Answer(userID, questionID uuid.UUID, isCorrect bool) (*ProgressResponse, error)
	Get(userID, questionID uuid.UUID) (*ProgressResponse, error)
	Upcoming(userID uuid.UUID, params common.PaginationParams) ([]UpcomingItem, int64, error)
	Saved(userID uuid.UUID, params common.PaginationParams) ([]UpcomingItem, int64, error)
	ToggleSave(userID, questionID uuid.UUID) (*ProgressResponse, error)
	TopicMastery(userID uuid.UUID) ([]TopicMasteryItem, error)
}

type progressService struct {
	store        Store
	gamification gamification.Service
}

func NewService(store Store, gamification gamification.Service) Service {
	return &progressService{store: store, gamification: gamification}
}

func (s *progressService) Answer(userID, questionID uuid.UUID, isCorrect bool) (*ProgressResponse, error) {
	p, err := s.store.FindByUserAndQuestion(userID, questionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			p = &models.UserQuestionProgress{
				UserID:     userID,
				QuestionID: questionID,
			}
		} else {
			return nil, apierr.ErrInternal("Error al obtener el progreso", "")
		}
	}

	applySM2(p, isCorrect)

	if err := s.store.Upsert(p); err != nil {
		return nil, apierr.ErrInternal("Error al guardar el progreso", "")
	}

	var reward *gamification.Reward
	if s.gamification != nil {
		reward, err = s.gamification.RegisterAnswer(userID, questionID, isCorrect)
		if err != nil {
			return nil, err
		}
	}

	resp := toProgressResponse(p)
	resp.Reward = reward
	return resp, nil
}

func (s *progressService) Get(userID, questionID uuid.UUID) (*ProgressResponse, error) {
	p, err := s.store.FindByUserAndQuestion(userID, questionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &ProgressResponse{
				QuestionID: questionID,
				EaseFactor: 2.5,
				IsSaved:    false,
				IsMastered: false,
			}, nil
		}
		return nil, apierr.ErrInternal("Error al obtener el progreso", "")
	}
	return toProgressResponse(p), nil
}

func (s *progressService) Upcoming(userID uuid.UUID, params common.PaginationParams) ([]UpcomingItem, int64, error) {
	items, total, err := s.store.FindUpcoming(userID, params)
	if err != nil {
		return nil, 0, apierr.ErrInternal("Error al listar las preguntas pendientes", "")
	}

	result := make([]UpcomingItem, len(items))
	for i, p := range items {
		result[i] = UpcomingItem{
			Question: questions.ToQuestionResponse(*p.Question),
			Progress: *toProgressResponse(&p),
		}
	}
	return result, total, nil
}

func (s *progressService) Saved(userID uuid.UUID, params common.PaginationParams) ([]UpcomingItem, int64, error) {
	items, total, err := s.store.FindSaved(userID, params)
	if err != nil {
		return nil, 0, apierr.ErrInternal("Error al listar las preguntas guardadas", "")
	}

	result := make([]UpcomingItem, len(items))
	for i, p := range items {
		result[i] = UpcomingItem{
			Question: questions.ToQuestionResponse(*p.Question),
			Progress: *toProgressResponse(&p),
		}
	}
	return result, total, nil
}

func (s *progressService) ToggleSave(userID, questionID uuid.UUID) (*ProgressResponse, error) {
	p, err := s.store.FindByUserAndQuestion(userID, questionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			p = &models.UserQuestionProgress{
				UserID:     userID,
				QuestionID: questionID,
			}
		} else {
			return nil, apierr.ErrInternal("Error al obtener el progreso", "")
		}
	}

	p.IsSaved = !p.IsSaved

	if err := s.store.Upsert(p); err != nil {
		return nil, apierr.ErrInternal("Error al guardar el progreso", "")
	}

	return toProgressResponse(p), nil
}

// TopicMastery devuelve el dominio del usuario por tema, ordenado para la UI:
// temas débiles primero, luego menor accuracy y finalmente por nombre.
//
// Los temas sin respuestas no aparecen (el store los excluye). Devuelve un
// slice vacío —nunca nil— para que el envelope sea `{"data":[]}`.
func (s *progressService) TopicMastery(userID uuid.UUID) ([]TopicMasteryItem, error) {
	rows, err := s.store.FindTopicAccuracy(userID)
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener el dominio por tema", "")
	}

	items := make([]TopicMasteryItem, 0, len(rows))
	for _, r := range rows {
		accuracy := 0.0
		if r.Attempts > 0 {
			accuracy = math.Round(float64(r.CorrectCount)/float64(r.Attempts)*accuracyRounding) / accuracyRounding
		}
		items = append(items, TopicMasteryItem{
			TopicID:      r.TopicID,
			TopicName:    r.TopicName,
			Attempts:     r.Attempts,
			CorrectCount: r.CorrectCount,
			Accuracy:     accuracy,
			IsWeak:       accuracy < weakAccuracyThreshold && r.Attempts >= minAttempts,
		})
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].IsWeak != items[j].IsWeak {
			return items[i].IsWeak
		}
		if items[i].Accuracy != items[j].Accuracy {
			return items[i].Accuracy < items[j].Accuracy
		}
		return items[i].TopicName < items[j].TopicName
	})

	return items, nil
}

func applySM2(p *models.UserQuestionProgress, correct bool) {
	now := time.Now()

	if correct {
		if p.Repetitions == 0 {
			p.IntervalDays = 1
		} else if p.Repetitions == 1 {
			p.IntervalDays = 3
		} else {
			p.IntervalDays = int(math.Round(float64(p.IntervalDays) * p.EaseFactor))
		}
		p.Repetitions++
		p.EaseFactor = math.Max(1.3, p.EaseFactor+0.1)
		if p.Repetitions >= 5 {
			p.IsMastered = true
		}
	} else {
		p.Repetitions = 0
		p.IntervalDays = 1
		p.EaseFactor = math.Max(1.3, p.EaseFactor-0.2)
		p.IsMastered = false
	}

	next := now.Add(time.Duration(p.IntervalDays) * 24 * time.Hour)
	p.NextReviewAt = &next
	p.LastReviewedAt = &now
}

func toProgressResponse(p *models.UserQuestionProgress) *ProgressResponse {
	return &ProgressResponse{
		QuestionID:     p.QuestionID,
		Repetitions:    p.Repetitions,
		EaseFactor:     p.EaseFactor,
		IntervalDays:   p.IntervalDays,
		NextReviewAt:   p.NextReviewAt,
		LastReviewedAt: p.LastReviewedAt,
		IsSaved:        p.IsSaved,
		IsMastered:     p.IsMastered,
	}
}
