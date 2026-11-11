<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'

const props = defineProps<{
  count: number
  url?: string
}>()

const to = computed(() =>
  props.url ? { path: props.url, hash: '#comments' } : undefined,
)
</script>

<template>
  <component
    :is="to ? RouterLink : 'span'"
    class="comments"
    :class="{ 'comments--link': to }"
    :to="to"
  >
    <svg
      class="comments__icon"
      width="14"
      height="14"
      viewBox="0 0 24 24"
      fill="none"
      aria-hidden="true"
    >
      <path
        d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5Z"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </svg>
    <span class="visually-hidden">Комментарии:</span>
    <span class="comments__count">{{ count }}</span>
  </component>
</template>

<style scoped>
.comments {
  display: inline-flex;
  flex: none;
  align-items: center;
  gap: 4px;
  font: 400 13px/16px Inter, sans-serif;
}

.comments--link {
  margin: -2px -4px;
  padding: 2px 4px;
  border-radius: 4px;
  transition: background-color 0.15s ease;
}

.comments--link:hover {
  background: var(--color-card);
}

.comments--link:hover .comments__count {
  text-decoration: none;
  text-underline-offset: 2px;
}

.comments__icon {
  width: 14px;
  height: 14px;
}
</style>
