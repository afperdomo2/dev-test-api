import type {
  AuthResponse,
  LoginRequest,
  SetupRequest,
  StatusResponse,
  User,
} from '@devtest/shared'

import { apiClient } from '@/core/api/client'

export async function login(data: LoginRequest): Promise<AuthResponse> {
  return apiClient.post<AuthResponse>('/api/v1/auth/login', data)
}

export async function setup(data: SetupRequest): Promise<AuthResponse> {
  return apiClient.post<AuthResponse>('/api/v1/auth/setup', data)
}

export async function getStatus(): Promise<StatusResponse> {
  return apiClient.get<StatusResponse>('/api/v1/auth/status')
}

export async function getProfile(): Promise<User> {
  return apiClient.get<User>('/api/v1/profile')
}
