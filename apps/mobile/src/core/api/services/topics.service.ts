import type { PaginatedResponse, Topic } from '@devtest/shared'

import { apiClient } from '@/core/api/client'

export interface ListTopicsParams {
  page?: number
  perPage?: number
  search?: string
  hasQuestions?: boolean
}

/**
 * Lista temas visibles para el usuario.
 *
 * `hasQuestions` evita ofrecer temas sin preguntas (una sesión vacía es un
 * dead-end). El backend pagina; con búsqueda server-side 100 por página cubre
 * el catálogo típico sin scroll infinito.
 */
export async function listTopics(
  params: ListTopicsParams = {},
): Promise<PaginatedResponse<Topic>> {
  return apiClient.get<PaginatedResponse<Topic>>('/api/v1/topics', {
    page: params.page ?? 1,
    perPage: params.perPage ?? 100,
    search: params.search,
    hasQuestions: params.hasQuestions,
  })
}
