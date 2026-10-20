import { queryOptions } from '@tanstack/vue-query'
import * as notificationsService from '@/api/services/notifications.service'

export function notificationsOptions(page: () => number, perPage: () => number) {
  return queryOptions({
    queryKey: ['notifications', 'list', page, perPage],
    queryFn: () => notificationsService.getNotifications(page(), perPage()),
    staleTime: 30 * 1000,
  })
}

export function unreadCountOptions() {
  return queryOptions({
    queryKey: ['notifications', 'unread-count'],
    queryFn: () => notificationsService.getUnreadCount(),
    // Sondeo ligero + refresco al volver a la pestaña (OQ-3).
    refetchInterval: 60 * 1000,
    refetchOnWindowFocus: true,
    staleTime: 30 * 1000,
  })
}

export function markReadMutation() {
  return {
    mutationKey: ['notifications', 'mark-read'],
    mutationFn: ({ id }: { id: string }) => notificationsService.markNotificationRead(id),
  }
}

export function markAllReadMutation() {
  return {
    mutationKey: ['notifications', 'mark-all-read'],
    mutationFn: notificationsService.markAllNotificationsRead,
  }
}
