import { queryOptions, infiniteQueryOptions } from '@tanstack/vue-query'
import * as sessionsService from '@/api/services/sessions.service'
import type { PaginatedResponse } from '@/types/api.types'
import type { Session } from '@/types/session.types'

export function sessionsListOptions(page: () => number, perPage: () => number) {
  return queryOptions({
    queryKey: ['sessions', 'list', page, perPage],
    queryFn: () => sessionsService.listSessions(page(), perPage()),
    staleTime: 30 * 1000,
  })
}

const SESSIONS_PER_PAGE = 20

export function sessionsInfiniteOptions(status: () => string | undefined) {
  return infiniteQueryOptions({
    queryKey: ['sessions', 'list', 'infinite', status],
    queryFn: ({ pageParam = 1 }) =>
      sessionsService.listSessions(pageParam as number, SESSIONS_PER_PAGE, status()),
    getNextPageParam: (
      lastPage: PaginatedResponse<Session>,
      allPages: Array<PaginatedResponse<Session>>,
    ) => {
      const totalFetched = allPages.reduce((sum, p) => sum + p.data.length, 0)
      return totalFetched < lastPage.meta.total ? allPages.length + 1 : undefined
    },
    initialPageParam: 1,
    staleTime: 30 * 1000,
  })
}

export function sessionDetailOptions(id: () => string) {
  return queryOptions({
    queryKey: ['sessions', 'detail', id],
    queryFn: () => sessionsService.getSessionById(id()),
  })
}

export function nextQuestionOptions(sessionId: () => string) {
  return queryOptions({
    queryKey: ['sessions', 'next', sessionId],
    queryFn: () => sessionsService.getNextQuestion(sessionId()),
    staleTime: 0,
  })
}

export function sessionSummaryOptions(sessionId: () => string) {
  return queryOptions({
    queryKey: ['sessions', 'summary', sessionId],
    queryFn: () => sessionsService.getSessionSummary(sessionId()),
    staleTime: 0,
  })
}

export function createSessionMutation() {
  return {
    mutationKey: ['sessions', 'create'],
    mutationFn: sessionsService.createSession,
  }
}

export function finishSessionMutation() {
  return {
    mutationKey: ['sessions', 'finish'],
    mutationFn: sessionsService.finishSession,
  }
}

export function submitAnswerMutation() {
  return {
    mutationKey: ['sessions', 'answer'],
    mutationFn: ({
      sessionId,
      data,
    }: {
      sessionId: string
      data: Parameters<typeof sessionsService.submitAnswer>[1]
    }) => sessionsService.submitAnswer(sessionId, data),
  }
}
