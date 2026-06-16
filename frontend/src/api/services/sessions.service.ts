import apiClient from '@/api/client'
import type {
  Session,
  CreateSessionRequest,
  SubmitAnswerRequest,
  SessionAnswer,
} from '@/types/session.types'
import type { Question } from '@/types/question.types'
import type { PaginatedResponse } from '@/types/api.types'

export async function listSessions(
  page: number,
  perPage: number,
): Promise<PaginatedResponse<Session>> {
  const res = await apiClient.get<PaginatedResponse<Session>>('/api/v1/sessions', {
    params: { page, per_page: perPage },
  })
  return res.data
}

export async function getSessionById(id: string): Promise<Session> {
  const res = await apiClient.get<Session>(`/api/v1/sessions/${id}`)
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

export async function getNextQuestion(id: string): Promise<Question> {
  const res = await apiClient.get<Question>(`/api/v1/sessions/${id}/next`)
  return res.data
}

export async function submitAnswer(
  sessionId: string,
  data: SubmitAnswerRequest,
): Promise<SessionAnswer> {
  const res = await apiClient.post<SessionAnswer>(`/api/v1/sessions/${sessionId}/answer`, data)
  return res.data
}
