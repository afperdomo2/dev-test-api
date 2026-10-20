package gamification

type StreakInfo struct {
	Current          int `json:"current"`
	Longest          int `json:"longest"`
	FreezesAvailable int `json:"freezesAvailable"`
}

type DailyGoalInfo struct {
	Current  int  `json:"current"`
	Target   int  `json:"target"`
	Progress int  `json:"progress"`
	Reached  bool `json:"reached"`
}

type Reward struct {
	XpGained  int           `json:"xpGained"`
	TotalXp   int           `json:"totalXp"`
	Level     int           `json:"level"`
	XpToLevel int           `json:"xpToLevel"`
	Streak    StreakInfo    `json:"streak"`
	DailyGoal DailyGoalInfo `json:"dailyGoal"`
}

type TopicXpItem struct {
	TopicID string `json:"topicId"`
	Name    string `json:"name"`
	Xp      int    `json:"xp"`
}

type OverviewResponse struct {
	TotalXp   int           `json:"totalXp"`
	Level     int           `json:"level"`
	XpToLevel int           `json:"xpToLevel"`
	Streak    StreakInfo    `json:"streak"`
	DailyGoal DailyGoalInfo `json:"dailyGoal"`
	TopicXp   []TopicXpItem `json:"topicXp"`
}

// BadgeItem es un logro con su estado actual.
type BadgeItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Earned      bool   `json:"earned"`
	Current     int    `json:"current"`
	Target      int    `json:"target"`
}

// BadgesResponse es el catálogo completo con el contador de conseguidos.
type BadgesResponse struct {
	Badges []BadgeItem `json:"badges"`
	Earned int         `json:"earned"`
	Total  int         `json:"total"`
}
