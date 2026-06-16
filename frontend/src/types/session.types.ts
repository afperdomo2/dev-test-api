export type SessionStatus = 'in_progress' | 'completed' | 'cancelled'
export type SessionMode = 'generate' | 'review'
export type SessionDifficulty = 'beginner' | 'intermediate' | 'advanced'

export const SESSION_STATUS_LABELS: Record<SessionStatus, string> = {
  in_progress: 'En progreso',
  completed: 'Completada',
  cancelled: 'Cancelada',
}

export const SESSION_STATUS_COLORS: Record<SessionStatus, string> = {
  in_progress: 'warning',
  completed: 'success',
  cancelled: 'grey',
}

export const SESSION_MODE_LABELS: Record<SessionMode, string> = {
  generate: 'Generar',
  review: 'Repasar',
}

export const SESSION_MODES: Array<{ title: string; value: SessionMode }> = [
  { title: 'Generar preguntas', value: 'generate' },
  { title: 'Repasar guardadas', value: 'review' },
]

export const SESSION_DIFFICULTIES: Array<{ title: string; value: SessionDifficulty }> = [
  { title: 'Principiante', value: 'beginner' },
  { title: 'Intermedio', value: 'intermediate' },
  { title: 'Avanzado', value: 'advanced' },
]

export interface Session {
  id: string
  userId: string
  name: string
  status: SessionStatus
  mode: SessionMode
  difficulty: SessionDifficulty
  score: number
  startedAt: string
  finishedAt: string | null
  topics: Array<{ id: string; slug: string; name: string }>
  answerCount: number
  createdAt: string
  updatedAt: string
}

export interface CreateSessionRequest {
  name: string
  mode: SessionMode
  difficulty: SessionDifficulty
  topic_ids: Array<string>
}

export interface SubmitAnswerRequest {
  question_id: string
  answer_text?: string
  selected_options?: Array<string>
  response_time_ms?: number
}

export interface SessionAnswer {
  id: string
  sessionId: string
  userId: string
  questionId: string
  answerText: string
  selectedOptions: Array<string>
  isCorrect: boolean
  aiFeedback: string
  responseTimeMs: number
}
