import { queryOptions } from '@tanstack/vue-query'
import * as progressService from '@/api/services/progress.service'

export function upcomingQuestionsOptions(params: { page: number; perPage: number }) {
  return queryOptions({
    queryKey: ['progress', 'upcoming', params],
    queryFn: () => progressService.getUpcomingReviews(params.page, params.perPage),
    staleTime: 0,
  })
}

export function savedQuestionsOptions(params: { page: number; perPage: number }) {
  return queryOptions({
    queryKey: ['progress', 'saved', params],
    queryFn: () => progressService.getSavedQuestions(params.page, params.perPage),
    staleTime: 30 * 1000,
  })
}

export function submitProgressAnswerMutation() {
  return {
    mutationKey: ['progress', 'answer'],
    mutationFn: ({ questionId, isCorrect }: { questionId: string; isCorrect: boolean }) =>
      progressService.submitProgressAnswer(questionId, isCorrect),
  }
}

export function toggleSaveMutation() {
  return {
    mutationKey: ['progress', 'toggle-save'],
    mutationFn: progressService.toggleSaveQuestion,
  }
}
