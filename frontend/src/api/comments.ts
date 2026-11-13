import { api } from './client'

/** Public author info embedded in every comment. */
export interface CommentAuthor {
  username: string
  avatar: string
  role: string
}

/** Comment as stored by the backend; `parentId` is null for top-level comments. */
export interface CommentPayload {
  id: number
  articleId: string
  parentId: number | null
  body: string
  createdAt: string
  author: CommentAuthor
}

export interface CommentsResponse {
  comments: CommentPayload[]
}

export interface CommentResponse {
  comment: CommentPayload
}

export interface CommentCreatePayload {
  body: string
  /** Id of the top-level comment being replied to. */
  parentId?: number
}

export function fetchComments(articleId: string) {
  return api<CommentsResponse>(`/articles/${encodeURIComponent(articleId)}/comments`)
}

export function createComment(articleId: string, payload: CommentCreatePayload) {
  return api<CommentResponse>(`/articles/${encodeURIComponent(articleId)}/comments`, {
    method: 'POST',
    body: payload,
  })
}

export function deleteComment(id: number) {
  return api<{ status: string }>(`/comments/${id}`, { method: 'DELETE' })
}
