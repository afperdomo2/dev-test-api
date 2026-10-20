package gamification

import (
	"math"
	"time"

	"github.com/felipe/dev-test-api/internal/models"
	"github.com/felipe/dev-test-api/pkg/apierr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	sessionBonusXp   = 20
	streakFreezeEarn = 7
	maxFreezes       = 3
	dateLayout       = "2006-01-02"
)

var xpByDifficulty = map[string]int{
	"beginner":     5,
	"intermediate": 10,
	"advanced":     15,
}

// XpForDifficulty expone el XP por respuesta correcta según la dificultad.
// Es la fuente única: `insights` la usa para derivar el XP de un periodo.
func XpForDifficulty(difficulty string) int {
	return xpByDifficulty[difficulty]
}

// SessionBonusXp expone el bonus fijo por sesión terminada.
func SessionBonusXp() int {
	return sessionBonusXp
}

type Service interface {
	RegisterAnswer(userID, questionID uuid.UUID, isCorrect bool) (*Reward, error)
	CompleteSession(userID uuid.UUID) (*Reward, error)
	GetOverview(userID uuid.UUID) (*OverviewResponse, error)
	GetDailyGoal(userID uuid.UUID) (int, error)
	SetDailyGoal(userID uuid.UUID, dailyGoal int) (*Reward, error)
	Freeze(userID uuid.UUID) (*Reward, error)
	GetBadges(userID uuid.UUID) (*BadgesResponse, error)
}

type gamificationService struct {
	store Store
}

func NewService(store Store) Service {
	return &gamificationService{store: store}
}

func (s *gamificationService) RegisterAnswer(userID, questionID uuid.UUID, isCorrect bool) (*Reward, error) {
	difficulty, topicIDs, err := s.store.GetQuestionMeta(questionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrValidation("La pregunta no existe", "")
		}
		return nil, apierr.ErrInternal("Error al obtener los datos de la pregunta", "")
	}

	xpGained := 0
	if isCorrect {
		xpGained = xpByDifficulty[difficulty]
	}

	streak, err := s.loadStreak(userID)
	if err != nil {
		return nil, err
	}
	dailyGoal, err := s.store.GetUserDailyGoal(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrValidation("El usuario no existe", "")
		}
		return nil, apierr.ErrInternal("Error al obtener la meta diaria", "")
	}

	applyActivity(streak, xpGained, dailyGoal)

	if err := s.store.UpsertStreak(streak); err != nil {
		return nil, apierr.ErrInternal("Error al guardar la racha", "")
	}

	userXp, err := s.loadXp(userID)
	if err != nil {
		return nil, err
	}
	userXp.Xp += xpGained
	if err := s.store.UpsertXp(userXp); err != nil {
		return nil, apierr.ErrInternal("Error al guardar el XP", "")
	}

	for _, topicID := range topicIDs {
		topicXp, err := s.loadTopicXp(userID, topicID)
		if err != nil {
			return nil, err
		}
		topicXp.Xp += xpGained
		if err := s.store.UpsertTopicXp(topicXp); err != nil {
			return nil, apierr.ErrInternal("Error al guardar el XP del tema", "")
		}
	}

	return toReward(userXp.Xp, xpGained, streak, dailyGoal), nil
}

func (s *gamificationService) CompleteSession(userID uuid.UUID) (*Reward, error) {
	streak, err := s.loadStreak(userID)
	if err != nil {
		return nil, err
	}
	dailyGoal, err := s.store.GetUserDailyGoal(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrValidation("El usuario no existe", "")
		}
		return nil, apierr.ErrInternal("Error al obtener la meta diaria", "")
	}

	applyActivity(streak, sessionBonusXp, dailyGoal)

	if err := s.store.UpsertStreak(streak); err != nil {
		return nil, apierr.ErrInternal("Error al guardar la racha", "")
	}

	userXp, err := s.loadXp(userID)
	if err != nil {
		return nil, err
	}
	userXp.Xp += sessionBonusXp
	if err := s.store.UpsertXp(userXp); err != nil {
		return nil, apierr.ErrInternal("Error al guardar el XP", "")
	}

	return toReward(userXp.Xp, sessionBonusXp, streak, dailyGoal), nil
}

func (s *gamificationService) GetOverview(userID uuid.UUID) (*OverviewResponse, error) {
	streak, err := s.loadStreak(userID)
	if err != nil {
		return nil, err
	}
	dailyGoal, err := s.store.GetUserDailyGoal(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrValidation("El usuario no existe", "")
		}
		return nil, apierr.ErrInternal("Error al obtener la meta diaria", "")
	}

	totalXp := 0
	userXp, err := s.store.GetXp(userID)
	if err == nil {
		totalXp = userXp.Xp
	} else if err != gorm.ErrRecordNotFound {
		return nil, apierr.ErrInternal("Error al obtener el XP", "")
	}

	topicXp, err := s.store.GetTopicXpOverview(userID)
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener el XP por tema", "")
	}

	reward := toReward(totalXp, 0, streak, dailyGoal)
	return &OverviewResponse{
		TotalXp:   totalXp,
		Level:     reward.Level,
		XpToLevel: reward.XpToLevel,
		Streak:    reward.Streak,
		DailyGoal: reward.DailyGoal,
		TopicXp:   topicXp,
	}, nil
}

func (s *gamificationService) GetDailyGoal(userID uuid.UUID) (int, error) {
	goal, err := s.store.GetUserDailyGoal(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, apierr.ErrValidation("El usuario no existe", "")
		}
		return 0, apierr.ErrInternal("Error al obtener la meta diaria", "")
	}
	return goal, nil
}

func (s *gamificationService) SetDailyGoal(userID uuid.UUID, dailyGoal int) (*Reward, error) {
	if dailyGoal < 10 || dailyGoal > 200 {
		return nil, apierr.ErrValidation("La meta diaria debe estar entre 10 y 200 XP", "")
	}

	if err := s.store.SetUserDailyGoal(userID, dailyGoal); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrValidation("El usuario no existe", "")
		}
		return nil, apierr.ErrInternal("Error al guardar la meta diaria", "")
	}

	streak, err := s.loadStreak(userID)
	if err != nil {
		return nil, err
	}
	totalXp := 0
	if userXp, err := s.store.GetXp(userID); err == nil {
		totalXp = userXp.Xp
	}
	return toReward(totalXp, 0, streak, dailyGoal), nil
}

func (s *gamificationService) Freeze(userID uuid.UUID) (*Reward, error) {
	streak, err := s.loadStreak(userID)
	if err != nil {
		return nil, err
	}
	dailyGoal, err := s.store.GetUserDailyGoal(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apierr.ErrValidation("El usuario no existe", "")
		}
		return nil, apierr.ErrInternal("Error al obtener la meta diaria", "")
	}

	if streak.FreezesAvailable <= 0 {
		return nil, apierr.ErrConflict("Sin freezes", "No tienes freezes disponibles", "")
	}

	today := time.Now().UTC().Format(dateLayout)
	streak.FreezesAvailable--
	streak.LastActivityDate = today
	now := time.Now()
	streak.LastActivityAt = &now

	if err := s.store.UpsertStreak(streak); err != nil {
		return nil, apierr.ErrInternal("Error al guardar la racha", "")
	}

	totalXp := 0
	if userXp, err := s.store.GetXp(userID); err == nil {
		totalXp = userXp.Xp
	}
	return toReward(totalXp, 0, streak, dailyGoal), nil
}

func (s *gamificationService) loadStreak(userID uuid.UUID) (*models.UserStreak, error) {
	streak, err := s.store.GetStreak(userID)
	if err == gorm.ErrRecordNotFound {
		return &models.UserStreak{
			UserID:        userID,
			CurrentStreak: 0,
			LongestStreak: 0,
		}, nil
	}
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener la racha", "")
	}
	return streak, nil
}

func (s *gamificationService) loadXp(userID uuid.UUID) (*models.UserXp, error) {
	xp, err := s.store.GetXp(userID)
	if err == gorm.ErrRecordNotFound {
		return &models.UserXp{UserID: userID, Xp: 0}, nil
	}
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener el XP", "")
	}
	return xp, nil
}

func (s *gamificationService) loadTopicXp(userID, topicID uuid.UUID) (*models.UserTopicXp, error) {
	xp, err := s.store.GetTopicXp(userID, topicID)
	if err == gorm.ErrRecordNotFound {
		return &models.UserTopicXp{UserID: userID, TopicID: topicID, Xp: 0}, nil
	}
	if err != nil {
		return nil, apierr.ErrInternal("Error al obtener el XP del tema", "")
	}
	return xp, nil
}

// applyActivity updates streak (continuity) and the daily XP counter in-place.
func applyActivity(streak *models.UserStreak, xpGained, dailyGoal int) {
	today := time.Now().UTC()
	todayStr := today.Format(dateLayout)

	// Daily XP counter (reset if a new day started).
	if streak.DailyXpDate != todayStr {
		streak.DailyXpToday = 0
		streak.DailyXpDate = todayStr
	}
	streak.DailyXpToday += xpGained

	// Streak continuity.
	lastDate := streak.LastActivityDate
	switch {
	case lastDate == "":
		streak.CurrentStreak = 1
	case lastDate == todayStr:
		// Already active today.
	default:
		last, err := time.Parse(dateLayout, lastDate)
		if err != nil {
			streak.CurrentStreak = 1
		} else {
			days := int(today.Sub(last.UTC()).Hours() / 24)
			if days <= 1 {
				streak.CurrentStreak++
			} else {
				// Gap > 1 day: try to cover missed days with freezes.
				missed := days - 1
				if streak.FreezesAvailable >= missed {
					streak.FreezesAvailable -= missed
					streak.CurrentStreak += days
				} else {
					streak.FreezesAvailable = 0
					streak.CurrentStreak = 1
				}
			}
		}
	}

	// Earn a freeze every 7 consecutive days (cap 3).
	if streak.CurrentStreak > 0 && streak.CurrentStreak%streakFreezeEarn == 0 && streak.FreezesAvailable < maxFreezes {
		streak.FreezesAvailable++
	}

	if streak.CurrentStreak > streak.LongestStreak {
		streak.LongestStreak = streak.CurrentStreak
	}

	now := time.Now()
	streak.LastActivityAt = &now
	streak.LastActivityDate = todayStr
}

func levelFromXp(xp int) int {
	return 1 + int(math.Sqrt(float64(xp)/100))
}

func xpToNextLevel(level int) int {
	return level * level * 100
}

func toReward(totalXp, xpGained int, streak *models.UserStreak, dailyGoal int) *Reward {
	level := levelFromXp(totalXp)
	return &Reward{
		XpGained:  xpGained,
		TotalXp:   totalXp,
		Level:     level,
		XpToLevel: xpToNextLevel(level) - totalXp,
		Streak: StreakInfo{
			Current:          streak.CurrentStreak,
			Longest:          streak.LongestStreak,
			FreezesAvailable: streak.FreezesAvailable,
		},
		DailyGoal: dailyGoalInfo(streak, dailyGoal),
	}
}

func dailyGoalInfo(streak *models.UserStreak, dailyGoal int) DailyGoalInfo {
	current := streak.DailyXpToday
	if streak.DailyXpDate != time.Now().UTC().Format(dateLayout) {
		current = 0
	}
	progress := 0
	if dailyGoal > 0 {
		progress = int(math.Round(float64(current) / float64(dailyGoal) * 100))
		if progress > 100 {
			progress = 100
		}
	}
	return DailyGoalInfo{
		Current:  current,
		Target:   dailyGoal,
		Progress: progress,
		Reached:  current >= dailyGoal,
	}
}
