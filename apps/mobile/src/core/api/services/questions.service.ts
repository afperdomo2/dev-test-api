import type { Question } from '@devtest/shared'
import type { PaginatedResponse } from '@devtest/shared'

import { apiClient } from '@/core/api/client'

export interface QuestionsFilters {
  type?: string
  difficulty?: string
  topicIds?: string
}

export async function listQuestions(
  page: number,
  perPage: number,
  filters?: QuestionsFilters,
): Promise<PaginatedResponse<Question>> {
  return apiClient.get('/api/v1/questions', { page, perPage, ...filters })
}

export async function getQuestionById(id: string): Promise<Question> {
  return apiClient.get(`/api/v1/questions/${id}`)
}
