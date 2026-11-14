/** Static section metadata for the article feeds; article rows live in PostgreSQL. */

export interface SectionConfig {
  id: string
  title: string
  href: string
}

/** Platform feeds for `/news/:platform`. */
export const newsPlatformSections: SectionConfig[] = [
  { id: 'pc', title: 'PC', href: '/news/pc' },
  { id: 'xbox', title: 'Xbox Series X|S', href: '/news/xbox' },
  { id: 'playstation', title: 'PlayStation 5', href: '/news/playstation' },
  { id: 'nintendo', title: 'Nintendo', href: '/news/nintendo' },
  { id: 'vr', title: 'VR', href: '/news/vr' },
  { id: 'hardware', title: 'Железо', href: '/news/hardware' },
]

/** Game hubs for `/walkthroughs/:section`. */
export const walkthroughSections: SectionConfig[] = [
  { id: 'recent', title: 'Недавние', href: '/walkthroughs/recent' },
  { id: 'control-resonant', title: 'CONTROL RESONANT', href: '/walkthroughs/control-resonant' },
  { id: 'silent-hill-townfall', title: 'SILENT HILL: TOWNFALL', href: '/walkthroughs/silent-hill-townfall' },
  { id: 'ghost-of-tsushima', title: 'GHOST OF TSUSHIMA', href: '/walkthroughs/ghost-of-tsushima' },
  { id: 'baldurs-gate-3', title: 'BALDUR’S GATE 3', href: '/walkthroughs/baldurs-gate-3' },
]

/**
 * Hub ids that mark walkthrough articles (the grouping `recent` excluded).
 * Mirrors the rateable set on the backend (articles.Rateable).
 */
export const walkthroughHubIds: string[] = [
  'control-resonant',
  'silent-hill-townfall',
  'ghost-of-tsushima',
  'baldurs-gate-3',
]

/** Walkthrough cards featured in the home page «Прохождения» section. */
export const walkthroughHighlightIds = ['walkthrough-1', 'walkthrough-2', 'walkthrough-3']

/** Review groups for `/reviews/:section`. */
export const reviewGroupSections: SectionConfig[] = [
  { id: 'recent', title: 'Недавние', href: '/reviews/recent' },
  { id: 'popular', title: 'Популярные', href: '/reviews/popular' },
]
