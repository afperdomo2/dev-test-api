package gamification

import (
	"sync"

	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	// Criterios de los logros (OQ-1 y OQ-2 de SPEC-badges.md).
	minPerfectAnswers    = 5
	fastSolverMaxAvgMs   = 10000
	minMasteredQuestions = 3
	maxAnswersBadge      = 100
	minMasteredTopics    = 5
)

type badgeCriteria struct {
	sessions       int
	longestStreak  int
	totalAnswers   int
	hasPerfect     bool
	masteredTopics int
	hasFast        bool
}

// GetBadges evalúa el catálogo completo on-the-fly (D-16: sin tabla).
func (s *gamificationService) GetBadges(userID uuid.UUID) (*BadgesResponse, error) {
	criteria, errs := s.loadBadgeCriteria(userID)
	if len(errs) > 0 {
		return nil, apierr.ErrInternal("Error al calcular los logros", "")
	}

	items := make([]BadgeItem, 0, len(badgeCatalog))
	earnedCount := 0

	for _, def := range badgeCatalog {
		earned, current := evaluateBadge(def, criteria)
		if earned {
			earnedCount++
		}
		items = append(items, BadgeItem{
			ID:          string(def.ID),
			Name:        def.Name,
			Description: def.Description,
			Icon:        def.Icon,
			Earned:      earned,
			Current:     current,
			Target:      def.Target,
		})
	}

	return &BadgesResponse{
		Badges: items,
		Earned: earnedCount,
		Total:  len(badgeCatalog),
	}, nil
}

// loadBadgeCriteria consulta las 6 métricas en paralelo (D-13: sin errgroup).
// Cada goroutine escribe un campo distinto del struct; el mutex protege solo
// el slice de errores.
func (s *gamificationService) loadBadgeCriteria(userID uuid.UUID) (badgeCriteria, []error) {
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
		c    badgeCriteria
	)

	collect := func(fn func() error) {
		defer wg.Done()
		if err := fn(); err != nil {
			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
		}
	}

	wg.Add(6)

	go collect(func() error {
		v, err := s.store.CountCompletedSessions(userID)
		c.sessions = v
		return err
	})
	go collect(func() error {
		v, err := s.store.CountAnswers(userID)
		c.totalAnswers = v
		return err
	})
	go collect(func() error {
		v, err := s.store.HasPerfectSession(userID)
		c.hasPerfect = v
		return err
	})
	go collect(func() error {
		v, err := s.store.HasFastSession(userID)
		c.hasFast = v
		return err
	})
	go collect(func() error {
		v, err := s.store.CountMasteredTopics(userID)
		c.masteredTopics = v
		return err
	})
	go collect(func() error {
		streak, err := s.store.GetStreak(userID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			return err
		}
		c.longestStreak = streak.LongestStreak
		return nil
	})

	wg.Wait()
	return c, errs
}

// evaluateBadge devuelve si el logro está conseguido y el progreso actual
// (acotado al objetivo para que la UI muestre `min(raw, target)/target`).
func evaluateBadge(def BadgeDefinition, c badgeCriteria) (bool, int) {
	var raw int

	switch def.ID {
	case BadgeFirstSession:
		raw = c.sessions
	case BadgeStreak3, BadgeStreak7, BadgeStreak30:
		raw = c.longestStreak
	case BadgeAnswers100:
		raw = c.totalAnswers
	case BadgeTopicsMastered5:
		raw = c.masteredTopics
	case BadgePerfectSession:
		if c.hasPerfect {
			return true, def.Target
		}
		return false, 0
	case BadgeFastSolver:
		if c.hasFast {
			return true, def.Target
		}
		return false, 0
	default:
		return false, 0
	}

	if raw >= def.Target {
		return true, def.Target
	}
	return false, raw
}
