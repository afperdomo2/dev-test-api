import apiClient from '@/api/client'
import type { Topic, CreateTopicRequest, UpdateTopicRequest } from '@/types/topic.types'
import type { PaginatedResponse } from '@/types/api.types'

export async function listTopics(
  page: number,
  perPage: number,
  sortBy?: string,
  sortOrder?: string,
  search?: string,
  myOnly?: boolean,
): Promise<PaginatedResponse<Topic>> {
  const params: Record<string, string | number | boolean> = { page, perPage }
  if (sortBy) params.sortBy = sortBy
  if (sortOrder) params.sortOrder = sortOrder
  if (search) params.search = search
  if (myOnly) params.myOnly = myOnly
  const res = await apiClient.get<PaginatedResponse<Topic>>('/api/v1/topics', { params })
  return res.data
}

export async function getTopicById(id: string): Promise<Topic> {
  const res = await apiClient.get<Topic>(`/api/v1/topics/${id}`)
  return res.data
}

export async function createTopic(data: CreateTopicRequest): Promise<Topic> {
  const res = await apiClient.post<Topic>('/api/v1/topics', data)
  return res.data
}

export async function updateTopic(id: string, data: UpdateTopicRequest): Promise<Topic> {
  const res = await apiClient.put<Topic>(`/api/v1/topics/${id}`, data)
  return res.data
}

export async function deleteTopic(id: string): Promise<void> {
  await apiClient.delete(`/api/v1/topics/${id}`)
}
