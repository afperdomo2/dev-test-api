export type SessionStatus = 'in_progress' | 'completed' | 'cancelled'
export type SessionMode = 'generate' | 'review'
export type SessionDifficulty = 'beginner' | 'intermediate' | 'advanced'

export interface Session {
  id: string
  user_id: string
  name: string
  status: SessionStatus
  mode: SessionMode
  difficulty: SessionDifficulty
  score: number
  started_at: string
  finished_at: string | null
  topics: Array<{ id: string; slug: string; name: string }>
  answer_count: number
  created_at: string
  updated_at: string
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
  session_id: string
  user_id: string
  question_id: string
  answer_text: string
  selected_options: Array<string>
  is_correct: boolean
  ai_feedback: string
  response_time_ms: number
}
