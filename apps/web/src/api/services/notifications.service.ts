import apiClient from '@/api/client'
import type { MarkAllReadResult, Notification, UnreadCount } from '@/types/notification.types'
import type { PaginatedResponse } from '@/types/api.types'

export async function getNotifications(
  page: number,
  perPage: number,
): Promise<PaginatedResponse<Notification>> {
  const res = await apiClient.get<PaginatedResponse<Notification>>('/api/v1/notifications', {
    params: { page, perPage },
  })
  return res.data
}

export async function getUnreadCount(): Promise<UnreadCount> {
  const res = await apiClient.get<UnreadCount>('/api/v1/notifications/unread-count')
  return res.data
}

export async function markNotificationRead(id: string): Promise<Notification> {
  const res = await apiClient.post<Notification>(`/api/v1/notifications/${id}/read`)
  return res.data
}

export async function markAllNotificationsRead(): Promise<MarkAllReadResult> {
  const res = await apiClient.post<MarkAllReadResult>('/api/v1/notifications/read-all')
  return res.data
}
