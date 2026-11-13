<script setup lang="ts">
defineProps<{
  /** Russian error message; a retry button is shown when set. */
  error?: string
}>()

const emit = defineEmits<{
  retry: []
}>()
</script>

<template>
  <div class="load-state">
    <template v-if="error">
      <p class="load-state__message">{{ error }}</p>
      <button class="load-state__retry" type="button" @click="emit('retry')">Повторить</button>
    </template>
    <p v-else class="load-state__message">Загрузка…</p>
  </div>
</template>

<style scoped>
.load-state {
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: center;
  padding: 24px 12px;
  background: var(--color-surface);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.load-state__message {
  color: var(--color-muted);
  font: 400 14px/18px Roboto, sans-serif;
}

.load-state__retry {
  min-height: 36px;
  padding: 6px 14px;
  color: var(--color-accent-text);
  background: var(--color-accent);
  border: 0;
  border-radius: var(--radius);
  font: 500 13px/16px Roboto, sans-serif;
  cursor: pointer;
}
</style>
