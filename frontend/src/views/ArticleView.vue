<script setup lang="ts">
import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ReviewStars from '@/components/article/ReviewStars.vue'
import ShareButton from '@/components/ui/ShareButton.vue'
import LoadState from '@/components/ui/LoadState.vue'
import { useAuthStore } from '@/stores/auth'
import { useArticlesStore } from '@/stores/articles'
import { formatDate } from '@/utils/date'
import { renderMarkdown } from '@/utils/markdown'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const store = useArticlesStore()

const article = computed(() =>
  store.loadState === 'ready' && typeof route.params.id === 'string'
    ? store.getArticle(route.params.id)
    : undefined,
)

const ready = computed(() => store.loadState === 'ready')

const canView = computed(() => {
  if (!article.value) return false
  if (article.value.status === 'published') return true

  return (
    auth.account?.role === 'owner' ||
    (article.value.author !== '' && article.value.author === auth.account?.username)
  )
})

const bodyHtml = computed(() =>
  article.value?.content ? renderMarkdown(article.value.content) : '',
)

const shareUrl = computed(() =>
  article.value ? new URL(article.value.url, window.location.origin).href : '',
)

watch(
  ready,
  (isReady) => {
    if (isReady && !article.value) void router.replace({ name: 'home' })
  },
  { immediate: true },
)
</script>

<template>
  <div v-if="article && canView" class="article container">
    <article class="article__card">
      <header class="article__header">
        <span class="article__kind">{{ store.articleKindLabel(article.kind) }}</span>
        <span v-if="article.status === 'draft'" class="article__badge">Черновик</span>

        <h1 class="article__title">{{ article.title }}</h1>

        <div class="article__meta">
          <time class="article__date" :datetime="article.date">{{ formatDate(article.date) }}</time>
          <span class="article__author">{{ article.author }}</span>
          <ShareButton :url="shareUrl" label="Поделиться статьёй" />
        </div>

        <div v-if="article.kind === 'review'" class="article__review-meta">
          <ReviewStars v-if="article.rating" :rating="article.rating" />
          <span v-if="article.platform" class="article__platform">{{ article.platform }}</span>
        </div>
      </header>

      <img
        v-if="article.image"
        class="article__cover"
        :src="article.image"
        :alt="article.title"
      />

      <div v-if="bodyHtml" class="article-body article__body" v-html="bodyHtml"></div>
      <p v-else class="article__placeholder">Текст статьи скоро появится.</p>

      <footer v-if="article.kind === 'review' && article.verdict" class="article__verdict">
        <h2 class="article__verdict-title">Вердикт</h2>
        <p class="article__verdict-text">{{ article.verdict }}</p>
      </footer>
    </article>
  </div>
  <LoadState
    v-else-if="!ready"
    class="article container"
    :error="store.errorMessage"
    @retry="store.reload()"
  />
</template>

<style scoped>
.article {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  padding-block: var(--gap);
}

.article__card {
  display: flex;
  flex-direction: column;
  gap: 16px;
  width: 100%;
  max-width: 1060px;
  margin-inline: auto;
  padding: 12px;
  background: var(--color-surface);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.article__header {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.article__kind {
  align-self: flex-start;
  padding: 4px 10px;
  color: var(--color-accent-text);
  background: var(--color-accent);
  border-radius: var(--radius);
  font: 500 11px/14px Roboto, sans-serif;
  text-transform: uppercase;
}

.article__badge {
  align-self: flex-start;
  padding: 4px 10px;
  color: var(--color-warning);
  background: var(--color-surface);
  border-radius: var(--radius);
  font: 500 11px/14px Roboto, sans-serif;
  text-transform: uppercase;
}

.article__title {
  font: 600 24px/30px Roboto, sans-serif;
}

.article__meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.article__date,
.article__author {
  color: var(--color-muted);
  font: 400 13px/16px Roboto, sans-serif;
}

.article__review-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
}

.article__platform {
  color: var(--color-muted);
  font: 400 13px/16px Roboto, sans-serif;
}

.article__cover {
  width: 100%;
  border-radius: var(--radius);
}

.article__body {
  padding: 8px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.article__placeholder {
  padding: 8px;
  color: var(--color-muted);
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  font: 400 14px/18px Roboto, sans-serif;
}

.article__verdict {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 8px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.article__verdict-title {
  font: 600 16px/20px Roboto, sans-serif;
}

.article__verdict-text {
  color: var(--color-muted);
  font: 400 14px/20px Roboto, sans-serif;
}
</style>
