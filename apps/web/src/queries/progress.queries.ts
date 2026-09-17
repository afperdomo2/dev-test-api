import { queryOptions } from '@tanstack/vue-query'
import * as progressService from '@/api/services/progress.service'

export function progressOptions(id: () => string) {
  return queryOptions({
    queryKey: ['progress', 'detail', id],
    queryFn: () => progressService.getProgress(id()),
    staleTime: 30 * 1000,
  })
}

export function upcomingQuestionsOptions(page: () => number, perPage: () => number) {
  return queryOptions({
    queryKey: ['progress', 'upcoming', page, perPage],
    queryFn: () => progressService.getUpcomingReviews(page(), perPage()),
    staleTime: 0,
  })
}

export function savedQuestionsOptions(page: () => number, perPage: () => number) {
  return queryOptions({
    queryKey: ['progress', 'saved', page, perPage],
    queryFn: () => progressService.getSavedQuestions(page(), perPage()),
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
