import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import type { Article, ArticleKind, ArticleSection } from '@/types/article'
import {
  mockArticles,
  newsPlatforms,
  walkthroughHighlights,
} from '@/data/articles'

const USER_ARTICLES_KEY = 'abberant-articles'
const OVERRIDES_KEY = 'abberant-articles-overrides'
const DELETED_KEY = 'abberant-articles-deleted'

type ArticleOverrides = Record<string, Partial<Article>>

function readStored<T>(key: string): T | null {
  try {
    const raw = localStorage.getItem(key)
    return raw ? (JSON.parse(raw) as T) : null
  } catch {
    return null
  }
}

function writeStored(key: string, value: unknown) {
  try {
    if (value === null) {
      localStorage.removeItem(key)
    } else {
      localStorage.setItem(key, JSON.stringify(value))
    }
  } catch {
    // Ignore storage errors (e.g. private browsing, quota exceeded).
  }
}

function byDateDesc(a: Article, b: Article): number {
  return b.date.localeCompare(a.date)
}

export const useArticlesStore = defineStore('articles', () => {
  /** Articles created in the editor, persisted in localStorage. */
  const userArticles = ref<Article[]>(readStored<Article[]>(USER_ARTICLES_KEY) ?? [])
  /** Edits of mock articles, keyed by id, persisted in localStorage. */
  const overrides = ref<ArticleOverrides>(readStored<ArticleOverrides>(OVERRIDES_KEY) ?? {})
  /** Ids of mock articles removed by the owner. */
  const deletedIds = ref<string[]>(readStored<string[]>(DELETED_KEY) ?? [])

  function persist() {
    writeStored(USER_ARTICLES_KEY, userArticles.value)
    writeStored(OVERRIDES_KEY, overrides.value)
    writeStored(DELETED_KEY, deletedIds.value)
  }

  const mockWithOverrides = computed(() =>
    mockArticles
      .filter((article) => !deletedIds.value.includes(article.id))
      .map((article) => {
        const override = overrides.value[article.id]

        return override ? { ...article, ...override } : article
      }),
  )

  /** Every article visible to the editor: mock seed plus user-created ones. */
  const all = computed(() => [...mockWithOverrides.value, ...userArticles.value])

  const publishedNews = computed(() =>
    all.value.filter((article) => article.kind === 'news' && article.status === 'published'),
  )

  const publishedReviews = computed(() =>
    all.value.filter((article) => article.kind === 'review' && article.status === 'published'),
  )

  /** Applies overrides and hides drafts/deleted entries for a static mock list. */
  function feed(articles: Article[]): Article[] {
    const overridden = new Map(mockWithOverrides.value.map((article) => [article.id, article]))

    return articles
      .filter((article) => !deletedIds.value.includes(article.id))
      .map((article) => overridden.get(article.id) ?? article)
      .filter((article) => article.status === 'published')
  }

  const newsSections = computed<ArticleSection[]>(() => [
    {
      id: 'all',
      title: 'Все',
      href: '/news/all',
      articles: [...publishedNews.value].sort(byDateDesc),
    },
    ...newsPlatforms.map((section) => ({
      ...section,
      articles: feed(section.articles),
    })),
  ])

  const reviewSections = computed<ArticleSection[]>(() => {
    const sorted = [...publishedReviews.value]

    return [
      { id: 'recent', title: 'Недавние', href: '/reviews/recent', articles: sorted.sort(byDateDesc) },
      {
        id: 'popular',
        title: 'Популярные',
        href: '/reviews/popular',
        articles: [...publishedReviews.value].sort((a, b) => b.comments - a.comments),
      },
    ]
  })

  const homeSections = computed<ArticleSection[]>(() => [
    {
      id: 'news',
      title: 'Новости',
      href: '/news',
      articles: [...publishedNews.value].sort(byDateDesc),
    },
    {
      id: 'reviews',
      title: 'Обзоры',
      href: '/reviews',
      articles: [...publishedReviews.value].sort(byDateDesc),
    },
    {
      id: 'walkthroughs',
      title: 'Прохождения',
      href: '/walkthroughs',
      articles: feed(walkthroughHighlights),
    },
  ])

  function getArticle(id: string): Article | undefined {
    return all.value.find((article) => article.id === id)
  }

  /** Insertion point for freshly created articles, matching mock id patterns. */
  function nextUserArticleId(): string {
    const used = new Set(all.value.map((article) => article.id))
    let index = 1

    while (used.has(`user-${index}`)) index += 1

    return `user-${index}`
  }

  function saveUserArticle(article: Article) {
    const index = userArticles.value.findIndex((item) => item.id === article.id)

    if (index >= 0) {
      userArticles.value[index] = article
    } else {
      userArticles.value.push(article)
    }

    persist()
  }

  function saveMockOverride(article: Article) {
    overrides.value = {
      ...overrides.value,
      [article.id]: {
        title: article.title,
        image: article.image,
        kind: article.kind,
        content: article.content,
        author: article.author,
        status: article.status,
        rating: article.rating,
        platform: article.platform,
        verdict: article.verdict,
      },
    }

    persist()
  }

  /** Creates or updates an article; mock seed entries are edited via overrides. */
  function saveArticle(article: Article) {
    if (mockArticles.some((item) => item.id === article.id)) {
      saveMockOverride(article)
    } else {
      saveUserArticle(article)
    }
  }

  function deleteArticle(id: string) {
    const index = userArticles.value.findIndex((article) => article.id === id)

    if (index >= 0) {
      userArticles.value.splice(index, 1)
    } else {
      deletedIds.value = [...deletedIds.value, id]
    }

    persist()
  }

  function articleKindLabel(kind: ArticleKind): string {
    return kind === 'review' ? 'Обзор' : 'Новость'
  }

  return {
    all,
    newsSections,
    reviewSections,
    homeSections,
    getArticle,
    nextUserArticleId,
    saveArticle,
    deleteArticle,
    articleKindLabel,
  }
})
