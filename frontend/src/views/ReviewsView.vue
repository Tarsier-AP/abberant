<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { ArticleSection } from '@/types/article'
import HomeSection from '@/components/home/HomeSection.vue'
import PagePagination from '@/components/ui/PagePagination.vue'
import { reviewGroups } from '@/data/articles'

/** Maximum number of article cards on one archive page; above it a pager appears. */
const PAGE_LIMIT = 8

/** Number of cards shown for each group on the /reviews index. */
const PREVIEW_LIMIT = 3

const route = useRoute()
const router = useRouter()
const root = ref<HTMLElement | null>(null)

const sectionId = computed(() => (route.params.section as string | undefined) ?? null)

const currentGroup = computed(
  () => reviewGroups.find((group) => group.id === sectionId.value) ?? null,
)

const totalPages = computed(() => {
  if (!currentGroup.value) return 1

  return Math.max(1, Math.ceil(currentGroup.value.articles.length / PAGE_LIMIT))
})

function clampPage(value: unknown): number {
  const parsed = Number(value)

  if (!Number.isInteger(parsed)) return 1

  return Math.min(Math.max(parsed, 1), totalPages.value)
}

const page = computed(() => clampPage(route.query.page))

const groups = computed<ArticleSection[]>(() => {
  if (currentGroup.value) {
    const start = (page.value - 1) * PAGE_LIMIT

    return [
      {
        ...currentGroup.value,
        articles: currentGroup.value.articles.slice(start, start + PAGE_LIMIT),
      },
    ]
  }

  return reviewGroups.map((group) => ({
    ...group,
    articles: group.articles.slice(0, PREVIEW_LIMIT),
  }))
})

const heading = computed(() =>
  currentGroup.value ? `Обзоры: ${currentGroup.value.title}` : 'Обзоры',
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
  currentGroup,
  (group) => {
    if (sectionId.value && !group) {
      void router.replace({ path: '/reviews' })
    }
  },
  { immediate: true },
)
</script>

<template>
  <div ref="root" class="reviews container">
    <h1 class="visually-hidden">{{ heading }}</h1>

    <HomeSection v-for="group in groups" :key="group.id" v-bind="group" dense />

    <PagePagination
      v-if="totalPages > 1"
      :page="page"
      :total-pages="totalPages"
      @update:page="goToPage"
    />
  </div>
</template>

<style scoped>
.reviews {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  padding-block: var(--gap);
}
</style>
