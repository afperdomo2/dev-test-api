import apiClient from '@/api/client'
import type { LoginRequest, SetupRequest, AuthResponse } from '@/types/auth.types'

export async function login(data: LoginRequest): Promise<AuthResponse> {
  const res = await apiClient.post<AuthResponse>('/api/v1/auth/login', data)
  return res.data
}

export async function setup(data: SetupRequest): Promise<AuthResponse> {
  const res = await apiClient.post<AuthResponse>('/api/v1/auth/setup', data)
  return res.data
}
