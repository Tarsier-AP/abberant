import type { SectionType } from '@/api/sections'

/** Server-managed section: a news platform feed or a walkthrough game hub. */
export interface Section {
  id: string
  type: SectionType
  title: string
  /** Archived sections keep their article URLs but disappear from feeds. */
  archived: boolean
}

export type { SectionType } from '@/api/sections'
