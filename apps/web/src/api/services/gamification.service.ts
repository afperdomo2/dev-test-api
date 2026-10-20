import apiClient from '@/api/client'
import type {
  OverviewResponse,
  Reward,
  DailyGoalRequest,
  BadgesResponse,
} from '@/types/gamification.types'

export async function getOverview(): Promise<OverviewResponse> {
  const res = await apiClient.get<OverviewResponse>('/api/v1/gamification/overview')
  return res.data
}

export async function setDailyGoal(data: DailyGoalRequest): Promise<Reward> {
  const res = await apiClient.put<Reward>('/api/v1/gamification/daily-goal', data)
  return res.data
}

export async function freezeStreak(): Promise<Reward> {
  const res = await apiClient.post<Reward>('/api/v1/gamification/streak/freeze')
  return res.data
}

/** Catálogo de logros con su estado, calculado al vuelo por el backend. */
export async function getBadges(): Promise<BadgesResponse> {
  const res = await apiClient.get<BadgesResponse>('/api/v1/gamification/badges')
  return res.data
}
