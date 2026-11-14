import { api } from './client'

export type ActivityType = 'comment' | 'vote' | 'bookmark'

export interface ActivityItem {
  type: ActivityType
  createdAt: string
  article: {
    id: string
    title: string
  }
  /** Comment events: the comment itself. */
  commentId?: number
  body?: string
  /** Vote events: stars the user gave. */
  stars?: number
}

export interface ActivityResponse {
  activity: ActivityItem[]
}

export function fetchActivity() {
  return api<ActivityResponse>('/activity')
}
