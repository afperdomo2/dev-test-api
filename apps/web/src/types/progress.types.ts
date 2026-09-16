import type { QuestionDifficulty, QuestionType } from './question.types'

export interface Progress {
  questionId: string
  repetitions: number
  easeFactor: number
  intervalDays: number
  nextReviewAt: string
  lastReviewedAt: string
  isSaved: boolean
  isMastered: boolean
}

export interface UpcomingQuestion {
  id: string
  content: string
  type: QuestionType
  difficulty: QuestionDifficulty
  topics: Array<{ id: string; slug: string; name: string }>
  nextReviewAt: string
}
