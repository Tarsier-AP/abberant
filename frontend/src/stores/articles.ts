import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import type { Article, ArticleKind, ArticleSection, Headline, ReaderRating } from '@/types/article'
import type { ArticlePatchPayload } from '@/api/articles'
import * as api from '@/api/articles'
import {
  newsPlatformSections,
  reviewGroupSections,
  walkthroughHighlightIds,
  walkthroughHubIds,
  walkthroughSections,
} from '@/data/sections'

export type ArticlesLoadState = 'idle' | 'loading' | 'ready' | 'error'

/** Number of headlines shown in the PopularNow block. */
const HEADLINE_LIMIT = 4

function byDateDesc(a: Article, b: Article): number {
  return b.date.localeCompare(a.date)
}

/** Adds the derived url to a backend article payload. */
function withUrl(payload: api.ArticlePayload): Article {
  return { ...payload, url: `/article/${payload.id}` }
}

export const useArticlesStore = defineStore('articles', () => {
  const articles = ref<Article[]>([])
  const loadState = ref<ArticlesLoadState>('idle')
  const errorMessage = ref('')

  /** Loads every article visible to the current session once per app start. */
  async function init() {
    if (loadState.value === 'loading' || loadState.value === 'ready') return

    loadState.value = 'loading'
    errorMessage.value = ''

    try {
      const { articles: loaded } = await api.fetchArticles()

      articles.value = loaded.map(withUrl)
      loadState.value = 'ready'
    } catch (error) {
      loadState.value = 'error'
      errorMessage.value =
        error instanceof Error ? error.message : 'Не удалось загрузить материалы.'
    }
  }

  /** Retries loading after a failure. */
  function reload() {
    if (loadState.value !== 'error') return

    loadState.value = 'idle'
    void init()
  }

  /** Re-fetches everything; used after the legacy data migration. */
  async function refresh() {
    loadState.value = 'idle'
    await init()
  }

  /** Inserts or replaces an article in the local cache. */
  function putArticle(article: Article) {
    const index = articles.value.findIndex((item) => item.id === article.id)

    if (index >= 0) {
      articles.value[index] = article
    } else {
      articles.value.push(article)
    }
  }

  /** Replaces the cached comment count after the comment section changes. */
  function setCommentCount(id: string, count: number) {
    const index = articles.value.findIndex((item) => item.id === id)
    const article = articles.value[index]

    if (article && article.comments !== count) {
      articles.value[index] = { ...article, comments: count }
    }
  }

  /** Every article visible to the current session. */
  const all = computed(() => articles.value)

  const publishedNews = computed(() =>
    all.value.filter((article) => article.kind === 'news' && article.status === 'published'),
  )

  const publishedReviews = computed(() =>
    all.value.filter((article) => article.kind === 'review' && article.status === 'published'),
  )

  /** Published articles assigned to a named section (platform feed or game hub). */
  function sectionFeed(sectionId: string): Article[] {
    return all.value
      .filter((article) => article.section === sectionId && article.status === 'published')
      .sort(byDateDesc)
  }

  const newsSections = computed<ArticleSection[]>(() => [
    {
      id: 'all',
      title: 'Все',
      href: '/news/all',
      articles: [...publishedNews.value].sort(byDateDesc),
    },
    ...newsPlatformSections.map((section) => ({
      ...section,
      articles: sectionFeed(section.id),
    })),
  ])

  const walkthroughFeeds = computed<ArticleSection[]>(() =>
    walkthroughSections.map((section) => ({
      ...section,
      articles: sectionFeed(section.id),
    })),
  )

  const reviewSections = computed<ArticleSection[]>(() => {
    const recentConfig = reviewGroupSections.find((config) => config.id === 'recent')
    const popularConfig = reviewGroupSections.find((config) => config.id === 'popular')

    if (!recentConfig || !popularConfig) return []

    const sorted = [...publishedReviews.value].sort(byDateDesc)

    return [
      { ...recentConfig, articles: sorted },
      {
        ...popularConfig,
        articles: [...publishedReviews.value].sort((a, b) => b.comments - a.comments),
      },
    ]
  })

  const homeSections = computed<ArticleSection[]>(() => {
    const byId = new Map(all.value.map((article) => [article.id, article]))

    const highlights = walkthroughHighlightIds
      .map((id) => byId.get(id))
      .filter((article): article is Article => article?.status === 'published')

    return [
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
        articles: highlights,
      },
    ]
  })

  /** PopularNow headlines, derived from the newest published articles. */
  const popularNow = computed<Headline[]>(() =>
    [...publishedNews.value]
      .sort(byDateDesc)
      .slice(0, HEADLINE_LIMIT)
      .map((article) => ({
        id: article.id,
        title: article.title,
        comments: article.comments,
        url: article.url,
      })),
  )

  function getArticle(id: string): Article | undefined {
    return all.value.find((article) => article.id === id)
  }

  async function createArticle(payload: api.ArticleCreatePayload): Promise<Article> {
    const { article } = await api.createArticle(payload)
    const mapped = withUrl(article)

    putArticle(mapped)

    return mapped
  }

  async function updateArticle(id: string, patch: ArticlePatchPayload): Promise<Article> {
    const { article } = await api.updateArticle(id, patch)
    const mapped = withUrl(article)

    putArticle(mapped)

    return mapped
  }

  async function deleteArticle(id: string): Promise<void> {
    await api.deleteArticle(id)

    articles.value = articles.value.filter((article) => article.id !== id)
  }

  function articleKindLabel(kind: ArticleKind): string {
    return kind === 'review' ? 'Обзор' : 'Новость'
  }

  /** Reviews and walkthrough guides accept reader votes; plain news do not. */
  function isRateable(article: Article): boolean {
    return article.kind === 'review' || walkthroughHubIds.includes(article.section ?? '')
  }

  /** Replaces the cached reader rating after a vote changes. */
  function setReaderRating(id: string, rating: ReaderRating) {
    const index = articles.value.findIndex((item) => item.id === id)
    const article = articles.value[index]

    if (article) {
      articles.value[index] = { ...article, readerRating: rating }
    }
  }

  async function voteArticle(id: string, stars: number): Promise<ReaderRating> {
    const { rating } = await api.voteArticle(id, stars)

    setReaderRating(id, rating)

    return rating
  }

  async function removeVote(id: string): Promise<ReaderRating> {
    const { rating } = await api.removeVote(id)

    setReaderRating(id, rating)

    return rating
  }

  return {
    articles,
    loadState,
    errorMessage,
    init,
    reload,
    refresh,
    all,
    publishedNews,
    publishedReviews,
    newsSections,
    walkthroughSections: walkthroughFeeds,
    reviewSections,
    homeSections,
    popularNow,
    getArticle,
    createArticle,
    updateArticle,
    deleteArticle,
    setCommentCount,
    articleKindLabel,
    isRateable,
    setReaderRating,
    voteArticle,
    removeVote,
  }
})
