/** Static section metadata for the article feeds.
 *
 * News platforms and walkthrough hubs are server-managed (`/api/v1/sections`,
 * loaded into the `sections` store). Only the feed groupings that never hold
 * article assignments stay configured here.
 */

export interface SectionConfig {
  id: string
  title: string
  href: string
}

/** «Недавние» walkthrough feed: every published walkthrough, newest first. */
export const walkthroughRecentSection: SectionConfig = {
  id: 'recent',
  title: 'Недавние',
  href: '/walkthroughs/recent',
}

/** Review groups for `/reviews/:section` (feed groupings, not assignments). */
export const reviewGroupSections: SectionConfig[] = [
  { id: 'recent', title: 'Недавние', href: '/reviews/recent' },
  { id: 'popular', title: 'Популярные', href: '/reviews/popular' },
]
