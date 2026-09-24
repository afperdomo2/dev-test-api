import apiClient from '@/api/client'
import type {
  Report,
  CreateReportRequest,
  ReviewReportRequest,
  ReportStatus,
} from '@/types/moderation.types'
import type { PaginatedResponse } from '@/types/api.types'

export async function createReport(data: CreateReportRequest): Promise<Report> {
  const res = await apiClient.post<Report>('/api/v1/reports', data)
  return res.data
}

export async function listReports(
  page: number,
  perPage: number,
  status?: ReportStatus,
): Promise<PaginatedResponse<Report>> {
  const params: Record<string, string | number> = { page, perPage }
  if (status) params.status = status
  const res = await apiClient.get<PaginatedResponse<Report>>('/api/v1/reports', { params })
  return res.data
}

export async function reviewReport(id: string, data: ReviewReportRequest): Promise<Report> {
  const res = await apiClient.put<Report>(`/api/v1/reports/${id}/review`, data)
  return res.data
}
