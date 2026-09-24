export type TipSource = 'manual' | 'imported'

export interface Tip {
  id: string
  userId: string
  topicId: string
  topic: string
  content: string
  source: TipSource
  createdAt: string
  updatedAt: string
}

export type TipResponse = Tip

export interface CreateTipRequest {
  topicId: string
  content: string
}

export interface UpdateTipRequest {
  topicId?: string
  content?: string
}

export type { ImportRowError, ImportResult } from './question.types'