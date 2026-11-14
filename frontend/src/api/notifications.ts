import { api } from './client'

export interface NotificationActor {
  username: string
  avatar: string
}

export interface NotificationArticle {
  id: string
  title: string
}

export interface NotificationPayload {
  id: number
  /** `reply` — an answer to the user's comment; `article` — new material. */
  type: 'reply' | 'article'
  createdAt: string
  read: boolean
  actor?: NotificationActor
  article?: NotificationArticle
  commentId?: number
}

export interface NotificationsResponse {
  notifications: NotificationPayload[]
}

export interface UnreadCountResponse {
  count: number
}

export interface MarkAllReadResponse {
  status: string
  marked: number
}

export interface SubscriptionsResponse {
  sections: string[]
}

export function fetchNotifications() {
  return api<NotificationsResponse>('/notifications')
}

export function fetchUnreadCount() {
  return api<UnreadCountResponse>('/notifications/unread-count')
}

export function markAllNotificationsRead() {
  return api<MarkAllReadResponse>('/notifications/read', { method: 'POST' })
}

export function fetchSubscriptions() {
  return api<SubscriptionsResponse>('/subscriptions')
}

export function subscribe(section: string) {
  return api<{ status: string }>(`/subscriptions/${encodeURIComponent(section)}`, {
    method: 'PUT',
  })
}

export function unsubscribe(section: string) {
  return api<{ status: string }>(`/subscriptions/${encodeURIComponent(section)}`, {
    method: 'DELETE',
  })
}
