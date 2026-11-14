import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import type { Section, SectionType } from '@/types/sections'
import type { SectionPayload } from '@/api/sections'
import * as api from '@/api/sections'

export type SectionsLoadState = 'idle' | 'loading' | 'ready' | 'error'

function withDefaults(payload: SectionPayload): Section {
  return { ...payload }
}

export const useSectionsStore = defineStore('sections', () => {
  const sections = ref<Section[]>([])
  const loadState = ref<SectionsLoadState>('idle')
  const errorMessage = ref('')

  /** Loads the section registry once per app start. */
  async function init() {
    if (loadState.value === 'loading' || loadState.value === 'ready') return

    loadState.value = 'loading'
    errorMessage.value = ''

    try {
      const { sections: loaded } = await api.fetchSections()

      sections.value = loaded.map(withDefaults)
      loadState.value = 'ready'
    } catch (error) {
      loadState.value = 'error'
      errorMessage.value =
        error instanceof Error ? error.message : 'Не удалось загрузить секции.'
    }
  }

  /** Retries loading after a failure. */
  function reload() {
    if (loadState.value !== 'error') return

    loadState.value = 'idle'
    void init()
  }

  /** Re-fetches everything; used after a management change elsewhere. */
  async function refresh() {
    loadState.value = 'idle'
    await init()
  }

  function putSection(section: SectionPayload) {
    const index = sections.value.findIndex((item) => item.id === section.id)

    if (index >= 0) {
      sections.value[index] = withDefaults(section)
    } else {
      sections.value.push(withDefaults(section))
    }
  }

  /** Live (non-archived) sections of one type, oldest first. */
  function activeOfType(type: SectionType): Section[] {
    return sections.value.filter((section) => section.type === type && !section.archived)
  }

  const newsPlatforms = computed(() => activeOfType('news'))

  const walkthroughHubs = computed(() => activeOfType('walkthrough'))

  const archived = computed(() => sections.value.filter((section) => section.archived))

  function byId(id: string): Section | undefined {
    return sections.value.find((section) => section.id === id)
  }

  /** Display title for an article's section id, archived ones included. */
  function titleOf(id: string): string {
    return byId(id)?.title ?? id
  }

  /** Adds a section; the server derives the slug from the title. */
  async function create(type: SectionType, title: string): Promise<Section> {
    const { section } = await api.createSection({ type, title })

    putSection(section)

    return withDefaults(section)
  }

  async function rename(id: string, title: string): Promise<Section> {
    const { section } = await api.updateSection(id, { title })

    putSection(section)

    return withDefaults(section)
  }

  async function setArchived(id: string, value: boolean): Promise<Section> {
    const { section } = await api.updateSection(id, { archived: value })

    putSection(section)

    return withDefaults(section)
  }

  return {
    sections,
    loadState,
    errorMessage,
    init,
    reload,
    refresh,
    newsPlatforms,
    walkthroughHubs,
    archived,
    byId,
    titleOf,
    create,
    rename,
    setArchived,
  }
})
