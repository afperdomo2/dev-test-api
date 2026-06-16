import apiClient from '@/api/client'
import type { Progress, UpcomingQuestion } from '@/types/progress.types'
import type { PaginatedResponse } from '@/types/api.types'

export async function submitProgressAnswer(questionId: string, isCorrect: boolean): Promise<void> {
  await apiClient.post(`/api/v1/progress/${questionId}/answer`, { is_correct: isCorrect })
}

export async function getUpcomingReviews(
  page: number,
  perPage: number,
): Promise<PaginatedResponse<UpcomingQuestion>> {
  const res = await apiClient.get<PaginatedResponse<UpcomingQuestion>>(
    '/api/v1/progress/upcoming',
    { params: { page, per_page: perPage } },
  )
  return res.data
}

export async function getSavedQuestions(
  page: number,
  perPage: number,
): Promise<PaginatedResponse<UpcomingQuestion>> {
  const res = await apiClient.get<PaginatedResponse<UpcomingQuestion>>('/api/v1/progress/saved', {
    params: { page, per_page: perPage },
  })
  return res.data
}

export async function toggleSaveQuestion(questionId: string): Promise<Progress> {
  const res = await apiClient.post<Progress>(`/api/v1/progress/${questionId}/toggle-save`)
  return res.data
}
