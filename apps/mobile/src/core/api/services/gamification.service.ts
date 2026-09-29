import type { OverviewResponse, Reward, DailyGoalRequest } from '@devtest/shared'

import { apiClient } from '@/core/api/client'

export async function getOverview(): Promise<OverviewResponse> {
  return apiClient.get('/api/v1/gamification/overview')
}

export async function setDailyGoal(data: DailyGoalRequest): Promise<Reward> {
  return apiClient.put('/api/v1/gamification/daily-goal', data)
}

export async function freezeStreak(): Promise<Reward> {
  return apiClient.post('/api/v1/gamification/streak/freeze')
}
