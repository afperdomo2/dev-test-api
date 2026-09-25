export type GenerationJobStatus = 'pending' | 'running' | 'completed' | 'failed'

export type GenerationDifficulty = 'beginner' | 'intermediate' | 'advanced' | 'random'

export interface GenerationJob {
  id: string
  requestedBy: string
  topicIds: Array<string>
  topics: Array<string>
  difficulty: string
  quantity: number
  status: GenerationJobStatus
  generated: number
  error?: string
  createdAt: string
  startedAt?: string
  finishedAt?: string
}

export interface GenerationStats {
  /** Generadas y guardadas hoy. */
  generatedToday: number
  /** Comprometido hoy (pending + running + completed). Es lo que se cobra. */
  reservedToday: number
  dailyCap: number
  perJobCap: number
  defaultQuantity: number
  pending: number
  running: number
  failed: number
  aiConfigured: boolean
}

export interface CreateGenerationJobRequest {
  topicIds: Array<string>
  difficulty: GenerationDifficulty
  /** Opcional: si se omite, el backend usa el ajuste ai.default_quantity. */
  quantity?: number
}

export interface Setting {
  key: string
  /** Siempre string: el backend tipa según `type`. */
  value: string
  type: 'string' | 'int' | 'float'
  description: string
  min?: number
  max?: number
  updatedAt?: string
  updatedBy?: string
}

export interface UpdateSettingRequest {
  key: string
  value: string
}

export interface UpdateSettingsRequest {
  settings: Array<UpdateSettingRequest>
}
