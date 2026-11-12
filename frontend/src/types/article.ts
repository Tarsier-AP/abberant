export type ArticleKind = 'news' | 'review'

export type ArticleStatus = 'draft' | 'published'

export interface Article {
  id: string
  title: string
  image: string
  date: string
  comments: number
  url: string
  kind: ArticleKind
  /** Markdown source of the body; an empty string means the body is not written yet. */
  content: string
  author: string
  status: ArticleStatus
  /** Review-only fields. */
  rating?: number
  platform?: string
  verdict?: string
}

export interface Headline {
  id: string
  title: string
  comments: number
  url: string
}

export interface ArticleSection {
  id: string
  title: string
  href: string
  articles: Article[]
}
