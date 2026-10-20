import type { TopicMasteryItem } from './progress.types'

export type InsightPeriod = 'week' | 'month'

export interface TopicInsightItem {
  topicId: string
  topicName: string
  attempts: number
  accuracy: number
}

export interface DailyActivityItem {
  /** YYYY-MM-DD (UTC) */
  date: string
  questionsAnswered: number
}

export interface InsightSummary {
  period: InsightPeriod
  from: string
  to: string
  sessions: number
  questionsAnswered: number
  accuracy: number
  timeSpentMs: number
  topTopics: Array<TopicInsightItem>
  /** Reutiliza el DTO de progress: "tema débil" tiene una sola definición. */
  weakTopics: Array<TopicMasteryItem>
  xpGained: number
  streakDays: number
  dailyActivity: Array<DailyActivityItem>
}

export const PERIOD_OPTIONS: Array<{ title: string; value: InsightPeriod }> = [
  { title: 'Semana', value: 'week' },
  { title: 'Mes', value: 'month' },
]
