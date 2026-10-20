import { queryOptions } from '@tanstack/vue-query'
import * as gamificationService from '@/api/services/gamification.service'

export function overviewOptions() {
  return queryOptions({
    queryKey: ['gamification', 'overview'],
    queryFn: () => gamificationService.getOverview(),
    staleTime: 60 * 1000,
  })
}

export function setDailyGoalMutation() {
  return {
    mutationKey: ['gamification', 'daily-goal'],
    mutationFn: gamificationService.setDailyGoal,
  }
}

export function freezeStreakMutation() {
  return {
    mutationKey: ['gamification', 'streak', 'freeze'],
    mutationFn: gamificationService.freezeStreak,
  }
}

export function badgesOptions() {
  return queryOptions({
    queryKey: ['gamification', 'badges'],
    queryFn: () => gamificationService.getBadges(),
    // Los logros cambian poco: no hace falta refrescarlos en cada visita.
    staleTime: 5 * 60 * 1000,
  })
}
