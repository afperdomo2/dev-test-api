export const DEFAULT_DAILY_IMPORT_LIMIT = 200

export interface User {
  id: string
  email: string
  isAdmin: boolean
  dailyImportLimit: number
  createdAt: string
  updatedAt: string
}

export interface CreateUserRequest {
  email: string
  password: string
  isAdmin?: boolean
  dailyImportLimit?: number
}

export interface UpdateUserRequest {
  password?: string
  isAdmin?: boolean
  dailyImportLimit?: number
}
