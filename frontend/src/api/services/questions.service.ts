import apiClient from '@/api/client'
import type { Question, CreateQuestionRequest, UpdateQuestionRequest } from '@/types/question.types'
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
    params: { page, perPage, ...filters },
  })
  return res.data
}

export async function getQuestionById(id: string): Promise<Question> {
  const res = await apiClient.get<Question>(`/api/v1/questions/${id}`)
  return res.data
}

export async function createQuestion(data: CreateQuestionRequest): Promise<Question> {
  const res = await apiClient.post<Question>('/api/v1/questions', data)
  return res.data
}

export async function updateQuestion(id: string, data: UpdateQuestionRequest): Promise<Question> {
  const res = await apiClient.put<Question>(`/api/v1/questions/${id}`, data)
  return res.data
}

export async function deleteQuestion(id: string): Promise<void> {
  await apiClient.delete(`/api/v1/questions/${id}`)
}
