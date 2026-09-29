import type {
  Session,
  SessionDetail,
  CreateSessionRequest,
  SubmitAnswerRequest,
  SessionAnswer,
  SessionSummary,
  NextQuestionItem,
  DailySessionResponse,
} from '@devtest/shared'
import type { PaginatedResponse } from '@devtest/shared'

import { apiClient } from '@/core/api/client'

export async function listSessions(
  page: number,
  perPage: number,
  status?: string,
): Promise<PaginatedResponse<Session>> {
  return apiClient.get('/api/v1/sessions', { page, perPage, status })
}

export async function getSessionById(id: string): Promise<Session> {
  const res = await apiClient.get<SessionDetail>(`/api/v1/sessions/${id}`)
  return res.session
}

export async function getSessionDetail(id: string): Promise<SessionDetail> {
  return apiClient.get(`/api/v1/sessions/${id}`)
}

export async function createSession(data: CreateSessionRequest): Promise<Session> {
  return apiClient.post('/api/v1/sessions', data)
}

export async function createDailySession(): Promise<DailySessionResponse> {
  return apiClient.post('/api/v1/sessions/daily')
}

export async function finishSession(id: string): Promise<Session> {
  return apiClient.put(`/api/v1/sessions/${id}/finish`)
}

export async function getNextQuestion(id: string): Promise<NextQuestionItem> {
  const res = await apiClient.get<{ question: NextQuestionItem }>(`/api/v1/sessions/${id}/next`)
  return res.question
}

export async function submitAnswer(
  sessionId: string,
  data: SubmitAnswerRequest,
): Promise<SessionAnswer> {
  return apiClient.post(`/api/v1/sessions/${sessionId}/answer`, data)
}

export async function getSessionSummary(id: string): Promise<SessionSummary> {
  return apiClient.get(`/api/v1/sessions/${id}/summary`)
}

export async function deleteSession(id: string): Promise<void> {
  await apiClient.delete(`/api/v1/sessions/${id}`)
}
