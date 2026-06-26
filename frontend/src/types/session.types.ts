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

export const SESSION_STATUS_FILTERS: Array<{ title: string; value: SessionStatus | undefined }> = [
  { title: 'Todas', value: undefined },
  { title: 'En progreso', value: 'in_progress' },
  { title: 'Completada', value: 'completed' },
  { title: 'Cancelada', value: 'cancelled' },
]

export interface Session {
  id: string
  userId: string
  name: string
  status: SessionStatus
  mode: SessionMode
  difficulty: SessionDifficulty
  questionLimit: number | null
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
  topicIds: Array<string>
  questionLimit?: number
}

export interface SubmitAnswerRequest {
  questionId: string
  answerText?: string
  selectedOptions?: Array<string>
  responseTimeMs?: number
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
