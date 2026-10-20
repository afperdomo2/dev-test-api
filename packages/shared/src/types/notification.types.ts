export type NotificationType = 'daily_review'

export interface Notification {
  id: string
  userId: string
  type: NotificationType
  title: string
  body: string
  isRead: boolean
  createdAt: string
}

export interface UnreadCount {
  count: number
}

export interface MarkAllReadResult {
  updated: number
}
