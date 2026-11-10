<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  page: number
  totalPages: number
}>()

const emit = defineEmits<{
  'update:page': [page: number]
}>()

const hasPrev = computed(() => props.page > 1)
const hasNext = computed(() => props.page < props.totalPages)
</script>

<template>
  <nav class="pagination" aria-label="Страницы раздела">
    <p class="pagination__status" aria-live="polite">
      Страница {{ page }} из {{ totalPages }}
    </p>

    <button
      class="pagination__arrow"
      type="button"
      aria-label="Предыдущая страница"
      :disabled="!hasPrev"
      @click="emit('update:page', page - 1)"
    >
      <svg width="8" height="14" viewBox="0 0 8 14" fill="none" aria-hidden="true">
        <path
          d="M7 1 1 7l6 6"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
    </button>

    <button
      class="pagination__arrow"
      type="button"
      aria-label="Следующая страница"
      :disabled="!hasNext"
      @click="emit('update:page', page + 1)"
    >
      <svg width="8" height="14" viewBox="0 0 8 14" fill="none" aria-hidden="true">
        <path
          d="m1 1 6 6-6 6"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
    </button>
  </nav>
</template>

<style scoped>
.pagination {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 10px;
}

.pagination__status {
  font: 700 24px/28px Roboto, sans-serif;
}

.pagination__arrow {
  display: grid;
  flex: none;
  place-items: center;
  width: 31px;
  height: 31px;
  padding: 0;
  color: var(--color-accent-text);
  background: var(--color-accent);
  border: 0;
  border-radius: 50%;
  cursor: pointer;
  transition: opacity 0.15s ease;
}

.pagination__arrow:hover:not(:disabled) {
  opacity: 0.8;
}

.pagination__arrow:disabled {
  opacity: 0.4;
  cursor: default;
}
</style>
