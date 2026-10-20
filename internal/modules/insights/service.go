package insights

import (
	"math"
	"sync"
	"time"

	"github.com/felipe/dev-test-api/internal/modules/gamification"
	"github.com/felipe/dev-test-api/internal/modules/progress"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
)

const (
	weekDays         = 7
	monthDays        = 30
	topTopicsLimit   = 5
	accuracyRounding = 10000
	dateLayout       = "2006-01-02"
)

type Service interface {
	Summary(userID uuid.UUID, period string) (*InsightSummaryResponse, error)
}

type insightsService struct {
	store        Store
	progress     progress.Service
	gamification gamification.Service
}

func NewService(store Store, progressSvc progress.Service, gamificationSvc gamification.Service) Service {
	return &insightsService{store: store, progress: progressSvc, gamification: gamificationSvc}
}

// Summary agrega la actividad de los últimos 7 o 30 días (ventana móvil).
//
// Las consultas son independientes y corren en paralelo con `sync.WaitGroup`
// (D-13: sin errgroup). Cada goroutine escribe SOLO su propia variable, así que
// no hay carrera; el mutex protege únicamente el slice de errores compartido.
func (s *insightsService) Summary(userID uuid.UUID, period string) (*InsightSummaryResponse, error) {
	if period == "" {
		period = "week"
	}
	days := weekDays
	if period == "month" {
		days = monthDays
	}

	now := time.Now().UTC()
	from := startOfUTCDay(now).AddDate(0, 0, -(days - 1))

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	sessions := 0
	answered, correct := 0, 0
	timeSpent := int64(0)
	var top []TopicAggregateRow
	var daily map[string]int
	var byDifficulty map[string]int
	var mastery []progress.TopicMasteryItem
	streakDays := 0

	collect := func(fn func() error) {
		defer wg.Done()
		if err := fn(); err != nil {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
		}
	}

	wg.Add(8)

	go collect(func() error {
		v, err := s.store.CountSessions(userID, from)
		sessions = v
		return err
	})
	go collect(func() error {
		total, ok, err := s.store.CountAnswers(userID, from)
		answered, correct = total, ok
		return err
	})
	go collect(func() error {
		v, err := s.store.SumResponseTimeMs(userID, from)
		timeSpent = v
		return err
	})
	go collect(func() error {
		v, err := s.store.TopTopics(userID, from, topTopicsLimit)
		top = v
		return err
	})
	go collect(func() error {
		v, err := s.store.DailyActivity(userID, from)
		daily = v
		return err
	})
	go collect(func() error {
		v, err := s.store.CountCorrectByDifficulty(userID, from)
		byDifficulty = v
		return err
	})
	go collect(func() error {
		v, err := s.progress.TopicMastery(userID)
		mastery = v
		return err
	})
	// La racha es un extra: si gamification falla, el informe sigue con 0.
	go func() {
		defer wg.Done()
		if overview, err := s.gamification.GetOverview(userID); err == nil && overview != nil {
			streakDays = overview.Streak.Current
		}
	}()

	wg.Wait()

	if len(errs) > 0 {
		return nil, apierr.ErrInternal("Error al calcular el informe", "")
	}

	accuracy := 0.0
	if answered > 0 {
		accuracy = round4(float64(correct) / float64(answered))
	}

	topItems := make([]TopicInsightItem, 0, len(top))
	for _, r := range top {
		acc := 0.0
		if r.Attempts > 0 {
			acc = round4(float64(r.CorrectCount) / float64(r.Attempts))
		}
		topItems = append(topItems, TopicInsightItem{
			TopicID:   r.TopicID,
			TopicName: r.TopicName,
			Attempts:  r.Attempts,
			Accuracy:  acc,
		})
	}

	weakItems := make([]progress.TopicMasteryItem, 0, len(mastery))
	for _, m := range mastery {
		if m.IsWeak {
			weakItems = append(weakItems, m)
		}
	}

	activity := make([]DailyActivityItem, 0, days)
	for i := 0; i < days; i++ {
		date := from.AddDate(0, 0, i).Format(dateLayout)
		activity = append(activity, DailyActivityItem{
			Date:              date,
			QuestionsAnswered: daily[date],
		})
	}

	xp := 0
	for difficulty, n := range byDifficulty {
		xp += n * gamification.XpForDifficulty(difficulty)
	}
	xp += sessions * gamification.SessionBonusXp()

	return &InsightSummaryResponse{
		Period:            period,
		From:              from,
		To:                now,
		Sessions:          sessions,
		QuestionsAnswered: answered,
		Accuracy:          accuracy,
		TimeSpentMs:       timeSpent,
		TopTopics:         topItems,
		WeakTopics:        weakItems,
		XpGained:          xp,
		StreakDays:        streakDays,
		DailyActivity:     activity,
	}, nil
}

func startOfUTCDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func round4(v float64) float64 {
	return math.Round(v*accuracyRounding) / accuracyRounding
}
