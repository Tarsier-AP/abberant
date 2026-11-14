<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import LoadState from '@/components/ui/LoadState.vue'
import { fetchActivity, type ActivityItem, type ActivityType } from '@/api/activity'
import { formatDate } from '@/utils/date'

type Filter = 'all' | ActivityType

const state = ref<'idle' | 'loading' | 'ready' | 'error'>('idle')
const errorMessage = ref('')
const items = ref<ActivityItem[]>([])
const filter = ref<Filter>('all')

const filters: Array<{ id: Filter; label: string }> = [
  { id: 'all', label: 'Всё' },
  { id: 'comment', label: 'Комментарии' },
  { id: 'vote', label: 'Оценки' },
  { id: 'bookmark', label: 'Закладки' },
]

const visible = computed(() =>
  filter.value === 'all' ? items.value : items.value.filter((item) => item.type === filter.value),
)

const counts = computed(() => {
  const result: Record<Filter, number> = { all: items.value.length, comment: 0, vote: 0, bookmark: 0 }

  for (const item of items.value) {
    result[item.type] += 1
  }

  return result
})

function itemText(item: ActivityItem): string {
  if (item.type === 'comment') return 'Комментарий'
  if (item.type === 'vote') return `Оценка: ${item.stars} из 5`
  if (item.type === 'bookmark') return 'Сохранили в закладки'

  return ''
}

async function load() {
  state.value = 'loading'
  errorMessage.value = ''

  try {
    const { activity } = await fetchActivity()

    items.value = activity
    state.value = 'ready'
  } catch (error) {
    state.value = 'error'
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось загрузить историю активности.'
  }
}

onMounted(() => {
  void load()
})
</script>

<template>
  <section class="activity container">
    <header class="activity__header">
      <h1 class="activity__title">История активности</h1>
      <p class="activity__hint">Ваши комментарии, оценки и закладки за последнее время.</p>
    </header>

    <LoadState v-if="state === 'loading' || state === 'error'" :error="errorMessage" @retry="load()" />

    <template v-else>
      <nav class="activity__filters" aria-label="Фильтр активности">
        <button
          v-for="option in filters"
          :key="option.id"
          class="activity__filter"
          type="button"
          :class="{ 'activity__filter--active': filter === option.id }"
          @click="filter = option.id"
        >
          {{ option.label }} <span class="activity__count">{{ counts[option.id] }}</span>
        </button>
      </nav>

      <p v-if="!visible.length" class="activity__empty">
        Здесь пока ничего нет. Комментарии, оценки и закладки появятся после
        первых действий на сайте.
      </p>

      <ol v-else class="activity__list">
        <li v-for="(item, index) in visible" :key="`${item.type}-${item.article.id}-${index}`" class="activity__item">
          <span class="activity__kind" :class="`activity__kind--${item.type}`">
            {{ itemText(item) }}
          </span>

          <div class="activity__body">
            <RouterLink class="activity__article" :to="`/article/${item.article.id}`">
              {{ item.article.title }}
            </RouterLink>

            <p v-if="item.type === 'comment' && item.body" class="activity__comment">
              {{ item.body }}
            </p>
          </div>

          <time class="activity__date" :datetime="item.createdAt">{{ formatDate(item.createdAt) }}</time>
        </li>
      </ol>
    </template>
  </section>
</template>

<style scoped>
.activity {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  padding-block: var(--gap);
}

.activity__header {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.activity__title {
  font: 600 clamp(20px, 2.6vw, 24px)/28px Roboto, sans-serif;
}

.activity__hint {
  color: var(--color-muted);
  font: 400 14px/18px Roboto, sans-serif;
}

.activity__filters {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.activity__filter {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  min-height: 32px;
  padding: 4px 12px;
  color: var(--color-text);
  background: var(--color-card);
  border: 0;
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  font: 500 12px/16px Roboto, sans-serif;
  cursor: pointer;
}

.activity__filter:hover {
  background: var(--color-card-hover);
}

.activity__filter--active {
  color: var(--color-accent-text);
  background: var(--color-accent);
}

.activity__count {
  color: inherit;
  opacity: 0.7;
  font-weight: 400;
}

.activity__list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.activity__item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 10px 12px;
  background: var(--color-surface);
  border-radius: var(--radius);
}

.activity__kind {
  flex: none;
  padding: 4px 10px;
  border-radius: var(--radius);
  font: 500 11px/14px Roboto, sans-serif;
  text-transform: uppercase;
  white-space: nowrap;
}

.activity__kind--comment {
  color: var(--color-accent-text);
  background: var(--color-accent);
}

.activity__kind--vote {
  color: var(--color-warning);
  background: var(--color-surface);
  box-shadow: var(--shadow);
}

.activity__kind--bookmark {
  color: var(--color-muted);
  background: var(--color-card);
  box-shadow: var(--shadow);
}

.activity__body {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.activity__article {
  font: 500 14px/18px Roboto, sans-serif;
}

.activity__article:hover {
  text-decoration: underline;
}

.activity__comment {
  color: var(--color-muted);
  font: 400 13px/17px Roboto, sans-serif;
}

.activity__date {
  flex: none;
  color: var(--color-muted);
  font: 400 12px/15px Roboto, sans-serif;
}

.activity__empty {
  padding: clamp(12px, 2vw, 16px);
  color: var(--color-muted);
  background: var(--color-surface);
  border-radius: var(--radius);
  font: 400 14px/20px Roboto, sans-serif;
}

@media (max-width: 600px) {
  .activity__item {
    flex-direction: column;
    gap: 6px;
  }
}
</style>
