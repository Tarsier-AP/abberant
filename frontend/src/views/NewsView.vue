<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { ArticleSection } from '@/types/article'
import HomeSection from '@/components/home/HomeSection.vue'
import PagePagination from '@/components/ui/PagePagination.vue'
import { useArticlesStore } from '@/stores/articles'

/** Maximum number of article cards on one news page (3x2 on desktop). */
const PAGE_LIMIT = 6

/** Number of cards shown for each section on the /news index. */
const PREVIEW_LIMIT = 6

const route = useRoute()
const router = useRouter()
const store = useArticlesStore()
const root = ref<HTMLElement | null>(null)

const sectionId = computed(() => (route.params.platform as string | undefined) ?? null)

const currentSection = computed(
  () => store.newsSections.find((section) => section.id === sectionId.value) ?? null,
)

const totalPages = computed(() => {
  if (!currentSection.value) return 1

  return Math.max(1, Math.ceil(currentSection.value.articles.length / PAGE_LIMIT))
})

function clampPage(value: unknown): number {
  const parsed = Number(value)

  if (!Number.isInteger(parsed)) return 1

  return Math.min(Math.max(parsed, 1), totalPages.value)
}

const page = computed(() => clampPage(route.query.page))

const groups = computed<ArticleSection[]>(() => {
  if (currentSection.value) {
    const start = (page.value - 1) * PAGE_LIMIT

    return [
      {
        ...currentSection.value,
        articles: currentSection.value.articles.slice(start, start + PAGE_LIMIT),
      },
    ]
  }

  return store.newsSections.map((section) => ({
    ...section,
    // Every section fills two desktop rows with six cards.
    articles: section.articles.slice(0, PREVIEW_LIMIT),
  }))
})

const heading = computed(() =>
  currentSection.value ? `Новости: ${currentSection.value.title}` : 'Новости',
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
      void router.replace({ path: '/news' })
    }
  },
  { immediate: true },
)
</script>

<template>
  <div ref="root" class="news container">
    <h1 class="visually-hidden">{{ heading }}</h1>

    <HomeSection v-for="group in groups" :key="group.id" v-bind="group" show-all />

    <PagePagination
      v-if="totalPages > 1"
      :page="page"
      :total-pages="totalPages"
      @update:page="goToPage"
    />
  </div>
</template>

<style scoped>
.news {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  padding-block: var(--gap);
}
</style>
