<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{ rating: number; label?: string }>(), {
  label: undefined,
})

const STARS = [1, 2, 3, 4, 5]

const title = computed(() => props.label ?? `Оценка: ${props.rating} из 5`)
</script>

<template>
  <span class="stars" role="img" :aria-label="title" :title="title">
    <span
      v-for="star in STARS"
      :key="star"
      class="stars__star"
      :class="{ 'stars__star--filled': star <= rating }"
      aria-hidden="true"
    >
      ★
    </span>
  </span>
</template>

<style scoped>
.stars {
  display: inline-flex;
  gap: 2px;
  color: var(--color-muted);
  font-size: 18px;
  line-height: 1;
}

.stars__star--filled {
  color: var(--color-accent);
}
</style>
