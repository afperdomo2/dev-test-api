import apiClient from '@/api/client'
import type { Tip, CreateTipRequest, UpdateTipRequest, ImportResult } from '@/types/tip.types'
import type { PaginatedResponse } from '@/types/api.types'

export interface TipsFilters {
  topicId?: string
}

export async function listTips(
  page: number,
  perPage: number,
  filters?: TipsFilters,
): Promise<PaginatedResponse<Tip>> {
  const res = await apiClient.get<PaginatedResponse<Tip>>('/api/v1/tips', {
    params: { page, perPage, ...filters },
  })
  return res.data
}

export async function getTipById(id: string): Promise<Tip> {
  const res = await apiClient.get<Tip>(`/api/v1/tips/${id}`)
  return res.data
}

export async function createTip(data: CreateTipRequest): Promise<Tip> {
  const res = await apiClient.post<Tip>('/api/v1/tips', data)
  return res.data
}

export async function updateTip(id: string, data: UpdateTipRequest): Promise<Tip> {
  const res = await apiClient.put<Tip>(`/api/v1/tips/${id}`, data)
  return res.data
}

export async function deleteTip(id: string): Promise<void> {
  await apiClient.delete(`/api/v1/tips/${id}`)
}

export async function importTips(formData: FormData): Promise<ImportResult> {
  const res = await apiClient.post<ImportResult>('/api/v1/tips/import', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return res.data
}

export async function importTipsFromText(content: string): Promise<ImportResult> {
  const formData = new FormData()
  formData.append('content', content)
  const res = await apiClient.post<ImportResult>('/api/v1/tips/import', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return res.data
}
