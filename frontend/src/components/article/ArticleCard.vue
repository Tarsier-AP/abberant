<script setup lang="ts">
import { computed } from 'vue'
import type { Article } from '@/types/article'
import { formatDate } from '@/utils/date'
import CommentCount from '@/components/ui/CommentCount.vue'

const props = defineProps<{ article: Article }>()

const date = computed(() => formatDate(props.article.date))
</script>

<template>
  <article class="card">
    <RouterLink
      class="card__media"
      :to="article.url"
      tabindex="-1"
      aria-hidden="true"
    >
      <img :src="article.image" :alt="article.title" loading="lazy" />
    </RouterLink>

    <h3 class="card__title">
      <RouterLink :to="article.url">{{ article.title }}</RouterLink>
    </h3>

    <div class="card__meta">
      <time class="card__date" :datetime="article.date">{{ date }}</time>
      <CommentCount :count="article.comments" />
      <RouterLink class="card__open" :to="article.url" aria-label="Открыть статью">
        <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
          <path
            d="M10.5 1.5 1.5 10.5M10.5 1.5H4M10.5 1.5V8"
            stroke="currentColor"
            stroke-width="1.5"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
      </RouterLink>
    </div>
  </article>
</template>

<style scoped>
.card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 8px;
  background: var(--color-surface);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  transition:
    box-shadow 0.15s ease,
    transform 0.15s ease;
}

.card:hover {
  box-shadow: var(--shadow-hover);
  transform: translateY(-2px);
}

.card__media img {
  width: 100%;
  aspect-ratio: 330 / 208;
  object-fit: cover;
  border-radius: var(--radius);
}

.card__title {
  padding-inline: 8px;
  font: 500 14px/16px Roboto, sans-serif;
}

.card__title a:hover {
  text-decoration: underline;
}

.card__meta {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  margin-top: auto;
}

.card__date {
  flex: 1;
  font: 500 14px/16px Roboto, sans-serif;
}

.card__open {
  display: grid;
  flex: none;
  place-items: center;
  width: 20px;
  height: 20px;
  border-radius: 4px;
}

.card__open:hover {
  background: var(--color-card);
}
</style>
