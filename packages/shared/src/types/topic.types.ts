export const TOPIC_CATEGORIES = [
  'lenguajes',
  'frontend',
  'backend',
  'base-datos',
  'datos',
  'devops',
  'cloud',
  'arquitectura',
  'conceptos',
  'testing',
  'seguridad',
  'ia',
  'movil',
] as const

export type TopicCategory = string

export interface Topic {
  id: string
  slug: string
  name: string
  category: string
  isSystem: boolean
  createdBy?: string
  createdAt: string
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
