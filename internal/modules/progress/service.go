package progress

import (
	"math"
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/internal/modules/questions"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Answer(userID, questionID uuid.UUID, isCorrect bool) (*ProgressResponse, error)
	Upcoming(userID uuid.UUID, page, perPage int) ([]UpcomingItem, int64, error)
	Saved(userID uuid.UUID, page, perPage int) ([]UpcomingItem, int64, error)
	ToggleSave(userID, questionID uuid.UUID) (*ProgressResponse, error)
}

type progressService struct {
	store Store
}

func NewService(store Store) Service {
	return &progressService{store: store}
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

	return toProgressResponse(p), nil
}

func (s *progressService) Upcoming(userID uuid.UUID, page, perPage int) ([]UpcomingItem, int64, error) {
	items, total, err := s.store.FindUpcoming(userID, page, perPage)
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

func (s *progressService) Saved(userID uuid.UUID, page, perPage int) ([]UpcomingItem, int64, error) {
	items, total, err := s.store.FindSaved(userID, page, perPage)
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
				IsSaved:    true,
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
