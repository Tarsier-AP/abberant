<script setup lang="ts">
defineProps<{
  title: string
  href: string
}>()
</script>

<template>
  <div class="section-header">
    <h2 class="section-header__title">{{ title }}</h2>
    <RouterLink
      class="section-header__link"
      :to="href"
      :aria-label="`Все материалы раздела «${title}»`"
    >
      <svg width="8" height="15" viewBox="0 0 8 15" fill="none" aria-hidden="true">
        <path
          d="m1 1 6 6.5L1 14"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
    </RouterLink>
  </div>
</template>

<style scoped>
.section-header {
  display: flex;
  align-items: center;
  gap: 6px;
}

.section-header__title {
  font: 700 24px/28px Roboto, sans-serif;
}

.section-header__link {
  display: grid;
  flex: none;
  place-items: center;
  width: 24px;
  height: 24px;
  border-radius: 4px;
}

.section-header__link:hover {
  background: var(--color-surface);
}
</style>
