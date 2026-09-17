import type { Question } from './question.types'

export interface Progress {
  questionId: string
  repetitions: number
  easeFactor: number
  intervalDays: number
  nextReviewAt?: string
  lastReviewedAt?: string
  isSaved: boolean
  isMastered: boolean
}

export interface ProgressItem {
  question: Question
  progress: Progress
}