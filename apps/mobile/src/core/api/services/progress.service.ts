import type { Progress, ProgressItem } from '@devtest/shared'
import type { PaginatedResponse } from '@devtest/shared'

import { apiClient } from '@/core/api/client'

export async function getUpcomingReviews(
  page: number,
  perPage: number,
): Promise<PaginatedResponse<ProgressItem>> {
  return apiClient.get('/api/v1/progress/upcoming', { page, perPage })
}

export async function getSavedQuestions(
  page: number,
  perPage: number,
): Promise<PaginatedResponse<ProgressItem>> {
  return apiClient.get('/api/v1/progress/saved', { page, perPage })
}

export async function toggleSaveQuestion(questionId: string): Promise<Progress> {
  return apiClient.post(`/api/v1/progress/${questionId}/toggle-save`)
}
