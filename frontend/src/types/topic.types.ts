export const TOPIC_CATEGORIES = [
  'lenguajes',
  'frontend',
  'backend',
  'devops',
  'arquitectura',
  'base-datos',
  'conceptos',
  'movil',
  'cloud',
  'testing',
  'seguridad',
  'ia',
] as const

export type TopicCategory = string

export interface Topic {
  id: string
  slug: string
  name: string
  category: string
  isSystem: boolean
  createdAt: string
  updatedAt: string
}

export interface CreateTopicRequest {
  slug: string
  name: string
  category: string
}

export interface UpdateTopicRequest {
  slug?: string
  name?: string
  category?: string
}
