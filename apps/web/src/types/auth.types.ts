import type { User } from './user.types'

export interface LoginRequest {
  email: string
  password: string
}

export interface SetupRequest {
  email: string
  password: string
}

export interface StatusResponse {
  initialized: boolean
}

export interface AuthResponse {
  token: string
  user: User
}
