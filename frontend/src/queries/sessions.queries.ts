import { queryOptions } from '@tanstack/vue-query'
import * as sessionsService from '@/api/services/sessions.service'

export function sessionsListOptions(page: () => number, perPage: () => number) {
  return queryOptions({
    queryKey: ['sessions', 'list', page, perPage],
    queryFn: () => sessionsService.listSessions(page(), perPage()),
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
