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
  /** Platform feed, game hub or review group id; empty for unassigned articles. */
  section?: string
  /** Aggregated reader vote; present on reviews and walkthrough guides. */
  readerRating?: ReaderRating
}

/** Reader vote summary as returned by the backend. */
export interface ReaderRating {
  average: number
  votes: number
  /** The current user's own stars; 0 means no vote yet. */
  myVote?: number
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
