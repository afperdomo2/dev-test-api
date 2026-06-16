import apiClient from '@/api/client'
import type { Question } from '@/types/question.types'
import type { PaginatedResponse } from '@/types/api.types'

export interface QuestionsFilters {
  type?: string
  difficulty?: string
}

export async function listQuestions(
  page: number,
  perPage: number,
  filters?: QuestionsFilters,
): Promise<PaginatedResponse<Question>> {
  const res = await apiClient.get<PaginatedResponse<Question>>('/api/v1/questions', {
    params: { page, per_page: perPage, ...filters },
  })
  return res.data
}

export async function getQuestionById(id: string): Promise<Question> {
  const res = await apiClient.get<Question>(`/api/v1/questions/${id}`)
  return res.data
}
