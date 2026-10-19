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

/**
 * Dominio del usuario en un tema. `isWeak` aplica la definición única del
 * backend (accuracy < 60% y al menos 5 respuestas).
 */
export interface TopicMasteryItem {
  topicId: string
  topicName: string
  attempts: number
  correctCount: number
  accuracy: number
  isWeak: boolean
}