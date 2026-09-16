import apiClient from '@/api/client'
import type { OpenCodeUsage } from '@/types/opencode.types'

export async function getOpenCodeUsage(): Promise<OpenCodeUsage> {
  const res = await apiClient.get<OpenCodeUsage>('/api/v1/opencode/usage')
  return res.data
}
