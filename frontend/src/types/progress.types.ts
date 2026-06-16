export interface Progress {
  question_id: string
  repetitions: number
  ease_factor: number
  interval_days: number
  next_review_at: string
  last_reviewed_at: string
  is_saved: boolean
  is_mastered: boolean
}

import type { QuestionDifficulty, QuestionType } from './question.types'

export interface UpcomingQuestion {
  id: string
  content: string
  type: QuestionType
  difficulty: QuestionDifficulty
  topics: Array<{ id: string; slug: string; name: string }>
  next_review_at: string
}
