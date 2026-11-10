<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import HomeSection from '@/components/home/HomeSection.vue'
import PagePagination from '@/components/ui/PagePagination.vue'
import { walkthroughs, walkthroughPageCount } from '@/data/articles'

const route = useRoute()
const router = useRouter()
const root = ref<HTMLElement | null>(null)

function clampPage(value: unknown): number {
  const parsed = Number(value)

  if (!Number.isInteger(parsed)) return 1

  return Math.min(Math.max(parsed, 1), walkthroughPageCount)
}

const page = computed(() => clampPage(route.query.page))

function goToPage(next: number) {
  const target = clampPage(next)

  if (target === page.value) return

  void router.push({ query: { ...route.query, page: String(target) } })
}

watch(page, async () => {
  await nextTick()
  root.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
})
</script>

<template>
  <div ref="root" class="walkthroughs container">
    <h1 class="visually-hidden">Прохождения</h1>

    <HomeSection v-for="section in walkthroughs" :key="section.id" v-bind="section" />

    <PagePagination
      :page="page"
      :total-pages="walkthroughPageCount"
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
