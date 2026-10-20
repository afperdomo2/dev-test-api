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

export type BadgeId =
  | 'first_session'
  | 'streak_3'
  | 'streak_7'
  | 'streak_30'
  | 'perfect_session'
  | 'answers_100'
  | 'topics_mastered_5'
  | 'fast_solver'

export interface Badge {
  id: BadgeId
  name: string
  description: string
  /** Nombre de MaterialCommunityIcons con prefijo `mdi-`. */
  icon: string
  earned: boolean
  current: number
  target: number
}

export interface BadgesResponse {
  badges: Array<Badge>
  earned: number
  total: number
}