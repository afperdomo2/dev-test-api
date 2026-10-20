import type { Progress, ProgressItem, TopicMasteryItem } from '@devtest/shared'
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

/**
 * Dominio por tema (sin paginar: el número de temas es acotado).
 * El backend ya devuelve los temas ordenados (débiles primero).
 */
export async function getTopicMastery(): Promise<Array<TopicMasteryItem>> {
  return apiClient.get<Array<TopicMasteryItem>>('/api/v1/progress/topics-mastery')
}
