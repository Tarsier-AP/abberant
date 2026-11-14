import { api } from './client'
import type { ArticlePayload } from './articles'

/** Ids of the current user's bookmarks, newest first. */
export interface BookmarkIdsResponse {
  ids: string[]
}

export interface BookmarksListResponse {
  articles: ArticlePayload[]
}

export function fetchBookmarkIds() {
  return api<BookmarkIdsResponse>('/bookmarks/ids')
}

export function fetchBookmarks() {
  return api<BookmarksListResponse>('/bookmarks')
}

export function addBookmark(articleId: string) {
  return api<{ status: string }>(`/bookmarks/${encodeURIComponent(articleId)}`, {
    method: 'PUT',
  })
}

export function removeBookmark(articleId: string) {
  return api<{ status: string }>(`/bookmarks/${encodeURIComponent(articleId)}`, {
    method: 'DELETE',
  })
}
