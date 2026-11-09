export interface Article {
  id: string
  title: string
  image: string
  date: string
  comments: number
  url: string
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
