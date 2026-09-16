import apiClient from '@/api/client'
import type {
  Session,
  SessionDetail,
  CreateSessionRequest,
  SubmitAnswerRequest,
  SessionAnswer,
  SessionSummary,
  NextQuestionItem,
} from '@/types/session.types'
import type { PaginatedResponse } from '@/types/api.types'

export async function listSessions(
  page: number,
  perPage: number,
  status?: string,
): Promise<PaginatedResponse<Session>> {
  const params: Record<string, string | number> = { page, perPage }
  if (status) params.status = status
  const res = await apiClient.get<PaginatedResponse<Session>>('/api/v1/sessions', {
    params,
  })
  return res.data
}

export async function getSessionById(id: string): Promise<Session> {
  const res = await apiClient.get<SessionDetail>(`/api/v1/sessions/${id}`)
  return res.data.session
}

export async function getSessionDetail(id: string): Promise<SessionDetail> {
  const res = await apiClient.get<SessionDetail>(`/api/v1/sessions/${id}`)
  return res.data
}

export async function createSession(data: CreateSessionRequest): Promise<Session> {
  const res = await apiClient.post<Session>('/api/v1/sessions', data)
  return res.data
}

export async function finishSession(id: string): Promise<Session> {
  const res = await apiClient.put<Session>(`/api/v1/sessions/${id}/finish`)
  return res.data
}

export async function getNextQuestion(id: string): Promise<NextQuestionItem> {
  const res = await apiClient.get<{ question: NextQuestionItem }>(`/api/v1/sessions/${id}/next`)
  return res.data.question
}

export async function submitAnswer(
  sessionId: string,
  data: SubmitAnswerRequest,
): Promise<SessionAnswer> {
  const res = await apiClient.post<SessionAnswer>(`/api/v1/sessions/${sessionId}/answer`, data)
  return res.data
}

export async function getSessionSummary(id: string): Promise<SessionSummary> {
  const res = await apiClient.get<SessionSummary>(`/api/v1/sessions/${id}/summary`)
  return res.data
}

export async function deleteSession(id: string): Promise<void> {
  await apiClient.delete(`/api/v1/sessions/${id}`)
}
