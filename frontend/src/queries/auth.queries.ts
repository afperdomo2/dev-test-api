import * as authService from '@/api/services/auth.service'
import type { LoginRequest, SetupRequest } from '@/types/auth.types'

export function loginMutation() {
  return {
    mutationKey: ['auth', 'login'],
    mutationFn: (data: LoginRequest) => authService.login(data),
  }
}

export function setupMutation() {
  return {
    mutationKey: ['auth', 'setup'],
    mutationFn: (data: SetupRequest) => authService.setup(data),
  }
}
