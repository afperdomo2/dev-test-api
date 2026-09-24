import apiClient from '@/api/client'
import type { OverviewResponse, Reward, DailyGoalRequest } from '@/types/gamification.types'

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
