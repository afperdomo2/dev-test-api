export interface StreakInfo {
  current: number
  longest: number
  freezesAvailable: number
}

export interface DailyGoalInfo {
  current: number
  target: number
  progress: number
  reached: boolean
}

export interface Reward {
  xpGained: number
  totalXp: number
  level: number
  xpToLevel: number
  streak: StreakInfo
  dailyGoal: DailyGoalInfo
}

export interface TopicXpItem {
  topicId: string
  name: string
  xp: number
}

export interface OverviewResponse {
  totalXp: number
  level: number
  xpToLevel: number
  streak: StreakInfo
  dailyGoal: DailyGoalInfo
  topicXp: Array<TopicXpItem>
}

export interface DailyGoalRequest {
  dailyGoalXp: number
}