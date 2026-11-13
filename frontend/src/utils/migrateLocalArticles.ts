import * as api from '@/api/articles'
import type { ArticlePatchPayload } from '@/api/articles'
import type { Article } from '@/types/article'
import { useArticlesStore } from '@/stores/articles'
import { useAuthStore } from '@/stores/auth'

const MIGRATED_KEY = 'abberant-articles-migrated'
const USER_ARTICLES_KEY = 'abberant-articles'
const OVERRIDES_KEY = 'abberant-articles-overrides'
const DELETED_KEY = 'abberant-articles-deleted'

function readStored<T>(key: string): T | null {
  try {
    const raw = localStorage.getItem(key)

    return raw ? (JSON.parse(raw) as T) : null
  } catch {
    return null
  }
}

/** True when the browser still has legacy local article state. */
function hasLegacyState(): boolean {
  return Boolean(
    localStorage.getItem(USER_ARTICLES_KEY) ||
      localStorage.getItem(OVERRIDES_KEY) ||
      localStorage.getItem(DELETED_KEY),
  )
}

function markMigrated() {
  try {
    localStorage.setItem(MIGRATED_KEY, '1')
    localStorage.removeItem(USER_ARTICLES_KEY)
    localStorage.removeItem(OVERRIDES_KEY)
    localStorage.removeItem(DELETED_KEY)
  } catch {
    // Ignore storage errors (e.g. private browsing).
  }
}

function buildCreatePayload(article: Article): api.ArticleCreatePayload {
  return {
    id: article.id,
    title: article.title,
    image: article.image,
    date: article.date,
    comments: article.comments,
    kind: article.kind,
    content: article.content,
    status: article.status,
    ...(article.rating !== undefined ? { rating: article.rating } : {}),
    ...(article.platform !== undefined ? { platform: article.platform } : {}),
    ...(article.verdict !== undefined ? { verdict: article.verdict } : {}),
    ...(article.section !== undefined ? { section: article.section } : {}),
  }
}

function buildPatchPayload(patch: Partial<Article>): ArticlePatchPayload {
  const payload: ArticlePatchPayload = {}

  if (patch.title !== undefined) payload.title = patch.title
  if (patch.image !== undefined) payload.image = patch.image
  if (patch.kind !== undefined) payload.kind = patch.kind
  if (patch.content !== undefined) payload.content = patch.content
  if (patch.status !== undefined) payload.status = patch.status
  if (patch.rating !== undefined) payload.rating = patch.rating
  if (patch.platform !== undefined) payload.platform = patch.platform
  if (patch.verdict !== undefined) payload.verdict = patch.verdict

  return payload
}

/**
 * Uploads legacy localStorage articles to the backend once, then clears the
 * local keys. Runs after the auth init; failures keep the keys so the upload
 * is retried on the next visit.
 */
export async function migrateLocalArticles(): Promise<void> {
  try {
    if (localStorage.getItem(MIGRATED_KEY)) return

    if (!hasLegacyState()) {
      markMigrated()

      return
    }

    const auth = useAuthStore()

    // Local articles were created through the editor, so only editors have them.
    if (!auth.canEditArticles) return

    const userArticles = readStored<Article[]>(USER_ARTICLES_KEY) ?? []
    const overrides = readStored<Record<string, Partial<Article>>>(OVERRIDES_KEY) ?? {}
    const deletedIds = readStored<string[]>(DELETED_KEY) ?? []

    const { articles: existing } = await api.fetchArticles()
    const existingIds = new Set(existing.map((article) => article.id))
    let changed = false

    for (const article of userArticles) {
      if (existingIds.has(article.id)) continue

      const { article: saved } = await api.createArticle(buildCreatePayload(article))

      existingIds.add(saved.id)
      changed = true
    }

    for (const [id, override] of Object.entries(overrides)) {
      if (deletedIds.includes(id) || !existingIds.has(id)) continue

      await api.updateArticle(id, buildPatchPayload(override))
      changed = true
    }

    for (const id of deletedIds) {
      if (!existingIds.has(id)) continue

      await api.deleteArticle(id)
      changed = true
    }

    markMigrated()

    if (changed) void useArticlesStore().refresh()
  } catch {
    // Keep the legacy keys; migration is retried on the next visit.
  }
}
