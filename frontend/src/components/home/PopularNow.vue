<script setup lang="ts">
import type { Headline } from '@/types/article'
import CommentCount from '@/components/ui/CommentCount.vue'

defineProps<{ items: Headline[] }>()
</script>

<template>
  <section class="popular" aria-labelledby="popular-title">
    <h2 id="popular-title" class="popular__title">Популярное</h2>

    <ul class="popular__list">
      <li v-for="item in items" :key="item.id" class="popular__item">
        <RouterLink class="popular__link" :to="item.url">{{ item.title }}</RouterLink>
        <CommentCount :count="item.comments" :url="item.url" />
      </li>
    </ul>
  </section>
</template>

<style scoped>
.popular {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 17px 16px;
  color: var(--color-text);
  background: var(--color-surface);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.popular__title {
  font: 700 16px/19px Roboto, sans-serif;
}

.popular__list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.popular__item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.popular__link {
  flex: 1;
  font: 400 14px/17px Inter, sans-serif;
}

.popular__link:hover {
  text-decoration: underline;
}
</style>
