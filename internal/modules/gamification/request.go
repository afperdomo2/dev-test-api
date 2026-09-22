package gamification

type DailyGoalRequest struct {
	DailyGoalXp int `json:"dailyGoalXp" binding:"required,min=10,max=200"`
}
