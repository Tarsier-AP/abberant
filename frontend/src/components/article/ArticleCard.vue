<script setup lang="ts">
import { computed } from 'vue'
import type { Article } from '@/types/article'
import { formatDate } from '@/utils/date'
import CommentCount from '@/components/ui/CommentCount.vue'
import ShareButton from '@/components/ui/ShareButton.vue'
import ReviewStars from '@/components/article/ReviewStars.vue'

const props = defineProps<{ article: Article }>()

const date = computed(() => formatDate(props.article.date))
const shareUrl = computed(() => new URL(props.article.url, window.location.origin).href)
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
      <span class="card__author">{{ article.author }}</span>
      <time class="card__date" :datetime="article.date">{{ date }}</time>
      <ReviewStars
        v-if="article.kind === 'review' && article.rating"
        :rating="article.rating"
        class="card__rating"
      />
      <CommentCount :count="article.comments" :url="article.url" />
      <ShareButton :url="shareUrl" label="Поделиться статьёй" />
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

.card__rating {
  font-size: 14px;
}

.card__author {
  color: var(--color-muted);
  font: 400 13px/16px Roboto, sans-serif;
}

.card__date {
  flex: 1;
  font: 500 14px/16px Roboto, sans-serif;
}
</style>
