import { apiClient } from '@/core/api/client'

export type PushPlatform = 'ios' | 'android'

export async function registerToken(token: string, platform: PushPlatform): Promise<void> {
  await apiClient.post('/api/v1/push/tokens', { token, platform })
}

export async function unregisterToken(token: string): Promise<void> {
  await apiClient.post('/api/v1/push/tokens/unregister', { token })
}
