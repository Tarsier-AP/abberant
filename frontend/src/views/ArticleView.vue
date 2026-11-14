<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ReviewStars from '@/components/article/ReviewStars.vue'
import ReaderStars from '@/components/article/ReaderStars.vue'
import GradeBadge from '@/components/article/GradeBadge.vue'
import CommentSection from '@/components/article/CommentSection.vue'
import ShareButton from '@/components/ui/ShareButton.vue'
import BookmarkButton from '@/components/ui/BookmarkButton.vue'
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
    auth.account?.role === 'editor' ||
    auth.account?.role === 'moderator' ||
    (article.value.author !== '' && article.value.author === auth.account?.username)
  )
})

/** Matches the editor page rules: managers and the author may edit. */
const canEdit = computed(() => {
  if (!article.value) return false

  return (
    auth.account?.role === 'owner' ||
    auth.account?.role === 'moderator' ||
    (article.value.author !== '' && article.value.author === auth.account?.username)
  )
})

const bodyHtml = computed(() => {
  if (!article.value?.content) return ''
  // Drop a leading markdown heading so it doesn't duplicate the page title.
  const stripped = article.value.content.replace(/^#{1,6}[^\n]*\n+/, '')
  return renderMarkdown(stripped)
})

const shareUrl = computed(() =>
  article.value ? new URL(article.value.url, window.location.origin).href : '',
)

const zoomed = ref<string | null>(null)

function openZoom(event: MouseEvent) {
  const target = event.target as HTMLElement
  if (target.tagName === 'IMG' && target.getAttribute('src')) {
    zoomed.value = target.getAttribute('src')
  }
}

function closeZoom() {
  zoomed.value = null
}

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

        <div v-if="article.kind === 'review'" class="article__review-meta">
          <ReviewStars v-if="article.rating" :rating="article.rating" />
          <span v-if="article.platform" class="article__platform">{{ article.platform }}</span>
        </div>

        <div class="article__meta">
          <time class="article__date" :datetime="article.date">{{ formatDate(article.date) }}</time>
          <span class="article__author">{{ article.author }}</span>
          <RouterLink
            v-if="canEdit"
            class="article__edit"
            :to="{ name: 'editor-edit', params: { id: article.id } }"
          >
            Редактировать
          </RouterLink>
          <ShareButton :url="shareUrl" label="Поделиться статьёй" />
          <BookmarkButton :article-id="article.id" />
        </div>
      </header>

      <div
        v-if="bodyHtml"
        class="article-body article__body"
        v-html="bodyHtml"
        @click="openZoom"
      ></div>
      <p v-else class="article__placeholder">Текст статьи скоро появится.</p>

      <div v-if="ready && store.isRateable(article)" class="article__rating">
        <GradeBadge
          v-if="article.readerRating && article.readerRating.votes > 0"
          :average="article.readerRating.average"
        />

        <ReaderStars
          :article-id="article.id"
          :rating="article.readerRating ?? { average: 0, votes: 0, myVote: 0 }"
        />
      </div>

      <CommentSection
        v-if="article.status === 'published'"
        class="article__comments"
        :article="article"
      />

      <Teleport to="body">
        <div v-if="zoomed" class="article__zoom" @click="closeZoom">
          <img :src="zoomed" alt="Увеличенное изображение" />
        </div>
      </Teleport>
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
  gap: clamp(10px, 1.5vw, 12px);
  width: 100%;
  max-width: 1060px;
  margin-inline: auto;
}

.article__header {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: clamp(8px, 1.5vw, 12px);
  background: var(--color-header);
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
  font: 600 clamp(18px, 2.4vw, 22px)/clamp(24px, 3vw, 28px) Roboto, sans-serif;
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

.article__edit {
  display: inline-flex;
  align-items: center;
  min-height: 30px;
  padding: 4px 12px;
  color: var(--color-text);
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  font: 500 12px/16px Roboto, sans-serif;
}

.article__edit:hover {
  background: var(--color-card-hover);
}

.article__review-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
}

.article__rating {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
  width: fit-content;
  padding: 8px 12px;
  background: var(--color-surface);
  border-radius: var(--radius);
}

.article__platform {
  color: var(--color-muted);
  font: 400 13px/16px Roboto, sans-serif;
}

.article__body {
  padding: clamp(8px, 1.5vw, 12px);
  background: var(--color-header);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.article__body :deep(img) {
  display: block;
  width: min(100%, 560px);
  height: auto;
  margin-inline: auto;
  border-radius: var(--radius);
  cursor: zoom-in;
}

.article__zoom {
  position: fixed;
  inset: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  background: rgb(0 0 0 / 80%);
  cursor: zoom-out;
}

.article__zoom img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.article__placeholder {
  padding: clamp(8px, 1.5vw, 12px);
  color: var(--color-muted);
  background: var(--color-header);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  font: 400 14px/18px Roboto, sans-serif;
}
</style>
