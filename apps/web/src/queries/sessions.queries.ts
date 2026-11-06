import { queryOptions } from '@tanstack/vue-query'
import * as sessionsService from '@/api/services/sessions.service'

type Getter<T> = () => T

export function sessionsListOptions(
  page: Getter<number>,
  perPage: Getter<number>,
  status: Getter<string | undefined>,
  sortBy: Getter<string | undefined>,
  sortOrder: Getter<string | undefined>,
) {
  return queryOptions({
    queryKey: ['sessions', 'list', page, perPage, status, sortBy, sortOrder],
    queryFn: () => sessionsService.listSessions(page(), perPage(), status(), sortBy(), sortOrder()),
    staleTime: 30 * 1000,
  })
}

/** Sesión en curso más reciente, para el banner "Continuar". */
export function sessionsResumeOptions() {
  return queryOptions({
    queryKey: ['sessions', 'list', 'resume'],
    queryFn: () => sessionsService.listSessions(1, 1, 'in_progress', 'started_at', 'desc'),
    staleTime: 30 * 1000,
  })
}

export function sessionDetailOptions(id: () => string) {
  return queryOptions({
    queryKey: ['sessions', 'detail', id],
    queryFn: () => sessionsService.getSessionById(id()),
  })
}

export function sessionDetailReviewOptions(id: () => string) {
  return queryOptions({
    queryKey: ['sessions', 'detail', 'review', id],
    queryFn: () => sessionsService.getSessionDetail(id()),
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

export function createDailySessionMutation() {
  return {
    mutationKey: ['sessions', 'daily'],
    mutationFn: sessionsService.createDailySession,
  }
}

export function finishSessionMutation() {
  return {
    mutationKey: ['sessions', 'finish'],
    mutationFn: sessionsService.finishSession,
  }
}

export function deleteSessionMutation() {
  return {
    mutationKey: ['sessions', 'delete'],
    mutationFn: sessionsService.deleteSession,
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
