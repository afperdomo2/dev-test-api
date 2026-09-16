export const DEFAULT_DAILY_IMPORT_LIMIT = 200
export const DEFAULT_DAILY_AI_LIMIT = 20

export interface User {
  id: string
  email: string
  isAdmin: boolean
  dailyImportLimit: number
  dailyAiLimit: number
  createdAt: string
  updatedAt: string
}

export interface CreateUserRequest {
  email: string
  password: string
  isAdmin?: boolean
  dailyImportLimit?: number
  dailyAiLimit?: number
}

export interface UpdateUserRequest {
  password?: string
  isAdmin?: boolean
  dailyImportLimit?: number
  dailyAiLimit?: number
}
