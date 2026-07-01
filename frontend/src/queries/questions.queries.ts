import { queryOptions } from '@tanstack/vue-query'
import * as questionsService from '@/api/services/questions.service'
import type { QuestionsFilters } from '@/api/services/questions.service'

export function questionsListOptions(
  page: () => number,
  perPage: () => number,
  filters: () => QuestionsFilters,
) {
  return queryOptions({
    queryKey: ['questions', 'list', page, perPage, filters],
    queryFn: () => questionsService.listQuestions(page(), perPage(), filters()),
    staleTime: 30 * 1000,
  })
}

export function questionDetailOptions(id: () => string) {
  return queryOptions({
    queryKey: ['questions', 'detail', id],
    queryFn: () => questionsService.getQuestionById(id()),
  })
}

export function createQuestionMutation() {
  return {
    mutationKey: ['questions', 'create'],
    mutationFn: questionsService.createQuestion,
  }
}

export function updateQuestionMutation() {
  return {
    mutationKey: ['questions', 'update'],
    mutationFn: ({
      id,
      data,
    }: {
      id: string
      data: Parameters<typeof questionsService.updateQuestion>[1]
    }) => questionsService.updateQuestion(id, data),
  }
}

export function deleteQuestionMutation() {
  return {
    mutationKey: ['questions', 'delete'],
    mutationFn: questionsService.deleteQuestion,
  }
}
