import { ref, watch } from 'vue'
import { defineStore } from 'pinia'
import type { Article } from '@/types/article'
import type { ArticlePayload } from '@/api/articles'
import * as api from '@/api/bookmarks'
import { useAuthStore } from '@/stores/auth'

export type BookmarksLoadState = 'idle' | 'loading' | 'ready' | 'error'

/** Adds the derived url to a backend article payload. */
function withUrl(payload: ArticlePayload): Article {
  return { ...payload, url: `/article/${payload.id}` }
}

export const useBookmarksStore = defineStore('bookmarks', () => {
  const auth = useAuthStore()

  /** Article ids bookmarked by the current account. */
  const ids = ref(new Set<string>())
  /** Whether the id set matches the current account (or its absence). */
  const loaded = ref(false)

  /** Full bookmarked articles, used by the /bookmarks page. */
  const list = ref<Article[]>([])
  const listState = ref<BookmarksLoadState>('idle')
  const listError = ref('')

  function clear() {
    ids.value = new Set()
    loaded.value = true
    list.value = []
    listState.value = 'idle'
    listError.value = ''
  }

  /** Loads the bookmark id set for the current account (or clears it). */
  async function loadIds() {
    if (!auth.account) {
      clear()
      return
    }

    try {
      const { ids: loadedIds } = await api.fetchBookmarkIds()

      ids.value = new Set(loadedIds)
      loaded.value = true
    } catch {
      // Keep whatever the user had; the buttons stay usable.
    }
  }

  /** Runs once per app start after the session is restored. */
  async function init() {
    if (loaded.value) return

    await loadIds()
  }

  function isBookmarked(id: string): boolean {
    return ids.value.has(id)
  }

  async function add(id: string) {
    ids.value.add(id)

    try {
      await api.addBookmark(id)
    } catch (error) {
      ids.value.delete(id)
      throw error
    }
  }

  async function remove(id: string) {
    ids.value.delete(id)

    try {
      await api.removeBookmark(id)
    } catch (error) {
      ids.value.add(id)
      throw error
    }
  }

  async function toggle(id: string) {
    return isBookmarked(id) ? remove(id) : add(id)
  }

  /** Loads the full bookmark list for the /bookmarks page. */
  async function loadList() {
    if (!auth.account) {
      clear()
      return
    }

    listState.value = 'loading'
    listError.value = ''

    try {
      const { articles: loaded } = await api.fetchBookmarks()

      list.value = loaded.map(withUrl)
      listState.value = 'ready'
    } catch (error) {
      listState.value = 'error'
      listError.value =
        error instanceof Error ? error.message : 'Не удалось загрузить закладки.'
    }
  }

  // Keep the id set in sync with logins and logouts.
  watch(
    () => auth.account,
    () => {
      void loadIds()
    },
  )

  return {
    ids,
    loaded,
    list,
    listState,
    listError,
    clear,
    loadIds,
    init,
    isBookmarked,
    add,
    remove,
    toggle,
    loadList,
  }
})
