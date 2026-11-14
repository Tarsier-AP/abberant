import { api } from './client'

export type SectionType = 'news' | 'walkthrough'

/** Server-managed section row; review groups are frontend-only. */
export interface SectionPayload {
  id: string
  type: SectionType
  title: string
  archived: boolean
}

export interface SectionsResponse {
  sections: SectionPayload[]
}

export interface SectionResponse {
  section: SectionPayload
}

export interface SectionCreatePayload {
  type: SectionType
  title: string
}

export interface SectionPatchPayload {
  title?: string
  archived?: boolean
}

export function fetchSections() {
  return api<SectionsResponse>('/sections')
}

export function createSection(payload: SectionCreatePayload) {
  return api<SectionResponse>('/sections', { method: 'POST', body: payload })
}

export function updateSection(id: string, patch: SectionPatchPayload) {
  return api<SectionResponse>(`/sections/${encodeURIComponent(id)}`, {
    method: 'PATCH',
    body: patch,
  })
}
