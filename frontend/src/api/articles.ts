import { api } from './client'
import type { ArticleKind, ArticleStatus } from '@/types/article'

/** Article as stored by the backend; `url` is derived on the client. */
export interface ArticlePayload {
  id: string
  title: string
  image: string
  date: string
  comments: number
  kind: ArticleKind
  content: string
  author: string
  status: ArticleStatus
  rating?: number
  platform?: string
  verdict?: string
  section?: string
}

export interface ArticlesResponse {
  articles: ArticlePayload[]
}

export interface ArticleResponse {
  article: ArticlePayload
}

export type ArticleCreatePayload = Omit<ArticlePayload, 'id' | 'date' | 'url' | 'author' | 'comments'> & {
  /** Requested id; the server assigns `user-N` when it is missing or taken. */
  id?: string
  /** Publication date as YYYY-MM-DD; defaults to today on the server. */
  date?: string
  /** Ignored by the server: the author is the logged-in user. */
  author?: string
  comments?: number
}

export type ArticlePatchPayload = Partial<
  Pick<
    ArticlePayload,
    | 'title'
    | 'image'
    | 'kind'
    | 'content'
    | 'status'
    | 'rating'
    | 'platform'
    | 'verdict'
    | 'section'
    | 'date'
  >
>

export function fetchArticles() {
  return api<ArticlesResponse>('/articles')
}

export function createArticle(payload: ArticleCreatePayload) {
  return api<ArticleResponse>('/articles', { method: 'POST', body: payload })
}

export function updateArticle(id: string, patch: ArticlePatchPayload) {
  return api<ArticleResponse>(`/articles/${encodeURIComponent(id)}`, {
    method: 'PATCH',
    body: patch,
  })
}

export function deleteArticle(id: string) {
  return api<{ status: string }>(`/articles/${encodeURIComponent(id)}`, { method: 'DELETE' })
}
