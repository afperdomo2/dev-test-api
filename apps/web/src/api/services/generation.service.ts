import apiClient from '@/api/client'
import type { PaginatedResponse } from '@/types/api.types'
import type {
  CreateGenerationJobRequest,
  GenerationJob,
  GenerationJobStatus,
  GenerationStats,
  Setting,
  UpdateSettingsRequest,
} from '@/types/generation.types'

export async function createGenerationJob(
  data: CreateGenerationJobRequest,
): Promise<GenerationJob> {
  const res = await apiClient.post<GenerationJob>('/api/v1/generation', data)
  return res.data
}

export async function listGenerationJobs(
  page: number,
  perPage: number,
  status?: GenerationJobStatus,
): Promise<PaginatedResponse<GenerationJob>> {
  const params: Record<string, string | number> = { page, perPage }
  if (status) params.status = status
  const res = await apiClient.get<PaginatedResponse<GenerationJob>>('/api/v1/generation', {
    params,
  })
  return res.data
}

export async function getGenerationStats(): Promise<GenerationStats> {
  const res = await apiClient.get<GenerationStats>('/api/v1/generation/stats')
  return res.data
}

export async function listSettings(): Promise<Array<Setting>> {
  const res = await apiClient.get<Array<Setting>>('/api/v1/settings')
  return res.data
}

export async function updateSettings(data: UpdateSettingsRequest): Promise<Array<Setting>> {
  const res = await apiClient.put<Array<Setting>>('/api/v1/settings', data)
  return res.data
}
