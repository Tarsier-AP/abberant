<script setup lang="ts">
defineProps<{
  title: string
}>()
</script>

<template>
  <section class="auth-card">
    <h1 class="auth-card__title">{{ title }}</h1>

    <div class="auth-card__panel">
      <slot />
    </div>
  </section>
</template>

<style scoped>
.auth-card {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
  max-width: 485px;
  padding: 8px;
  background: var(--color-surface);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.auth-card__title {
  padding: 8px;
  font: 700 26px/32px Inter, sans-serif;
  text-align: center;
}

.auth-card__panel {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}
</style>
