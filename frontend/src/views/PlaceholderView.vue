<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const title = computed(() => String(route.meta.title ?? 'Раздел'))
</script>

<template>
  <div class="placeholder container">
    <h1 class="placeholder__title">{{ title }}</h1>
    <p class="placeholder__text">Раздел находится в разработке.</p>
    <RouterLink class="placeholder__back" to="/">← Вернуться на главную</RouterLink>
  </div>
</template>

<style scoped>
.placeholder {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 16px;
  padding-block: 64px;
}

.placeholder__title {
  font: 700 32px/39px Inter, sans-serif;
}

.placeholder__text {
  font: 400 16px/24px Roboto, sans-serif;
}

.placeholder__back {
  font: 500 16px/24px Roboto, sans-serif;
  text-decoration: underline;
}
</style>
