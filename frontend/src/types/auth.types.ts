export interface LoginRequest {
  email: string
  password: string
}

export interface SetupRequest {
  email: string
  password: string
}

import type { User } from './user.types'

export interface AuthResponse {
  token: string
  user: User
}
