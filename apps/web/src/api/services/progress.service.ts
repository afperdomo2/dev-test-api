import apiClient from '@/api/client'
import type { Progress, ProgressItem } from '@/types/progress.types'
import type { PaginatedResponse } from '@/types/api.types'

export async function submitProgressAnswer(questionId: string, isCorrect: boolean): Promise<void> {
  await apiClient.post(`/api/v1/progress/${questionId}/answer`, { isCorrect })
}

export async function getProgress(questionId: string): Promise<Progress> {
  const res = await apiClient.get<Progress>(`/api/v1/progress/${questionId}`)
  return res.data
}

export async function getUpcomingReviews(
  page: number,
  perPage: number,
): Promise<PaginatedResponse<ProgressItem>> {
  const res = await apiClient.get<PaginatedResponse<ProgressItem>>('/api/v1/progress/upcoming', {
    params: { page, perPage },
  })
  return res.data
}

export async function getSavedQuestions(
  page: number,
  perPage: number,
): Promise<PaginatedResponse<ProgressItem>> {
  const res = await apiClient.get<PaginatedResponse<ProgressItem>>('/api/v1/progress/saved', {
    params: { page, perPage },
  })
  return res.data
}

export async function toggleSaveQuestion(questionId: string): Promise<Progress> {
  const res = await apiClient.post<Progress>(`/api/v1/progress/${questionId}/toggle-save`)
  return res.data
}
