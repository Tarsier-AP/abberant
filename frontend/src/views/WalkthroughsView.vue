<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { ArticleSection } from '@/types/article'
import HomeSection from '@/components/home/HomeSection.vue'
import PagePagination from '@/components/ui/PagePagination.vue'
import { useMediaQuery } from '@/composables/useMediaQuery'
import { walkthroughs } from '@/data/articles'

/** Maximum number of article cards on one page; above it a pager appears. */
const PAGE_LIMIT = 8

const route = useRoute()
const router = useRouter()
const root = ref<HTMLElement | null>(null)

const isLaptop = useMediaQuery('(min-width: 1024px)')
const previewCount = computed(() => (isLaptop.value ? 3 : 2))

const sectionId = computed(() => (route.params.section as string | undefined) ?? null)

const currentSection = computed(
  () => walkthroughs.find((section) => section.id === sectionId.value) ?? null,
)

/** Archive pages: whole sections only, filled until the card limit is hit. */
const archivePages = computed<ArticleSection[][]>(() => {
  const pages: ArticleSection[][] = []
  let page: ArticleSection[] = []
  let cards = 0

  for (const section of walkthroughs) {
    const preview = Math.min(section.articles.length, previewCount.value)

    if (page.length > 0 && cards + preview > PAGE_LIMIT) {
      pages.push(page)
      page = []
      cards = 0
    }

    page.push(section)
    cards += preview
  }

  if (page.length > 0) pages.push(page)

  return pages
})

const totalPages = computed(() => {
  if (currentSection.value) {
    return Math.max(1, Math.ceil(currentSection.value.articles.length / PAGE_LIMIT))
  }

  return Math.max(1, archivePages.value.length)
})

function clampPage(value: unknown): number {
  const parsed = Number(value)

  if (!Number.isInteger(parsed)) return 1

  return Math.min(Math.max(parsed, 1), totalPages.value)
}

const page = computed(() => clampPage(route.query.page))

const pageSections = computed<ArticleSection[]>(() => {
  if (currentSection.value) {
    const start = (page.value - 1) * PAGE_LIMIT

    return [
      {
        ...currentSection.value,
        articles: currentSection.value.articles.slice(start, start + PAGE_LIMIT),
      },
    ]
  }

  const sections = archivePages.value[page.value - 1] ?? []

  return sections.map((section) => ({
    ...section,
    articles: section.articles.slice(0, previewCount.value),
  }))
})

const heading = computed(() =>
  currentSection.value ? `Прохождения: ${currentSection.value.title}` : 'Прохождения',
)

function goToPage(next: number) {
  const target = clampPage(next)

  if (target === page.value) return

  void router.push({ query: { ...route.query, page: String(target) } })
}

watch(
  () => [page.value, sectionId.value],
  async () => {
    await nextTick()
    root.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  },
)

watch(
  currentSection,
  (section) => {
    if (sectionId.value && !section) {
      void router.replace({ path: '/walkthroughs' })
    }
  },
  { immediate: true },
)
</script>

<template>
  <div ref="root" class="walkthroughs container">
    <h1 class="visually-hidden">{{ heading }}</h1>

    <HomeSection v-for="section in pageSections" :key="section.id" v-bind="section" dense />

    <PagePagination
      v-if="totalPages > 1"
      :page="page"
      :total-pages="totalPages"
      @update:page="goToPage"
    />
  </div>
</template>

<style scoped>
.walkthroughs {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  padding-block: var(--gap);
}
</style>
