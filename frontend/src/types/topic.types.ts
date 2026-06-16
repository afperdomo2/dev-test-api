export const TOPIC_CATEGORIES = [
  'lenguajes',
  'frontend',
  'backend',
  'devops',
  'arquitectura',
  'base-datos',
  'conceptos',
] as const

export type TopicCategory = (typeof TOPIC_CATEGORIES)[number]

export interface Topic {
  id: string
  slug: string
  name: string
  category: string
  is_system: boolean
  created_at: string
  updated_at: string
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
