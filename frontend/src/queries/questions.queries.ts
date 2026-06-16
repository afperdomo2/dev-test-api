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
