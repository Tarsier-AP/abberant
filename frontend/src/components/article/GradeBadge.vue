<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    /** Reader grade from 0 to 5, shown with one decimal and a comma. */
    average: number
    /** Larger variant for the article page header. */
    large?: boolean
  }>(),
  { large: false },
)

/** One decimal, comma as the decimal separator for the Russian UI. */
const text = computed(() => props.average.toFixed(1).replace('.', ','))

const title = computed(() => `Оценка читателей: ${text.value} из 5`)

const tone = computed(() => (props.average >= 4 ? 'high' : props.average >= 3 ? 'mid' : 'low'))
</script>

<template>
  <span
    class="grade"
    :class="[`grade--${tone}`, { 'grade--large': large }]"
    role="img"
    :aria-label="title"
    :title="title"
  >
    {{ text }}
  </span>
</template>

<style scoped>
.grade {
  display: inline-flex;
  flex: none;
  align-items: center;
  justify-content: center;
  min-width: 44px;
  padding: 3px 10px;
  color: #fff;
  border-radius: var(--radius);
  font: 600 14px/17px Roboto, sans-serif;
}

.grade--large {
  min-width: 64px;
  padding: 8px 14px;
  font: 600 22px/26px Roboto, sans-serif;
}

.grade--high {
  background: var(--color-grade-high);
}

.grade--mid {
  background: var(--color-grade-mid);
}

.grade--low {
  background: var(--color-grade-low);
}
</style>
