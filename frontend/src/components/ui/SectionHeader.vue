<script setup lang="ts">
defineProps<{
  title: string
  href: string
}>()
</script>

<template>
  <RouterLink class="section-header" :to="href">
    <h2 class="section-header__title">{{ title }}</h2>
    <svg
      class="section-header__arrow"
      width="8"
      height="15"
      viewBox="0 0 8 15"
      fill="none"
      aria-hidden="true"
    >
      <path
        d="m1 1 6 6.5L1 14"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
    </svg>
  </RouterLink>
</template>

<style scoped>
.section-header {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  width: fit-content;
  padding-inline: 6px;
  border-radius: 4px;
}

.section-header:hover {
  background: var(--color-surface);
}

.section-header__title {
  font: 700 24px/28px Roboto, sans-serif;
}

.section-header__arrow {
  flex: none;
}
</style>
