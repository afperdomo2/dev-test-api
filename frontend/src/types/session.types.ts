export type SessionStatus = 'in_progress' | 'completed' | 'cancelled'
export type SessionMode = 'generate' | 'review'
export type SessionDifficulty = 'beginner' | 'intermediate' | 'advanced'

export const SESSION_DIFFICULTY_LABELS: Record<SessionDifficulty, string> = {
  beginner: 'Principiante',
  intermediate: 'Intermedio',
  advanced: 'Avanzado',
}

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
  questionsGenerated: number
  score: number
  startedAt: string
  finishedAt: string | null
  topics: Array<string>
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
  explanation: string
  responseTimeMs: number
}

export interface SessionSummary {
  answerCount: number
  questionsGenerated: number
  status: SessionStatus
  questionLimit: number | null
}

export interface QuestionSnapshotOption {
  id: string
  content: string
  isCorrect: boolean
}

export interface QuestionSnapshot {
  content: string
  type: string
  difficulty: string
  options: Array<QuestionSnapshotOption>
  codeChallenge: CodeChallenge | null
  topics: Array<string>
}

export interface SessionAnswerDetail {
  id: string
  questionId: string
  answerText: string
  selectedOptions: Array<string>
  isCorrect: boolean
  aiFeedback: string
  explanation: string
  responseTimeMs: number
  question: QuestionSnapshot
  createdAt: string
}

export interface SessionDetail {
  session: Session
  answers: Array<SessionAnswerDetail>
}

export interface NextQuestionItem {
  id: string
  type: string
  content: string
  difficulty: string
  options: Array<{ id: string; content: string; isCorrect?: boolean }>
  codeChallenge: CodeChallenge | null
  topics: Array<string>
}

export interface CodeChallenge {
  id: string
  starterCode: string
  expectedOutput: string
  language: string
  testCases: string
}
