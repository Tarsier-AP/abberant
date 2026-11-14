<script setup lang="ts">
import { computed, ref } from 'vue'
import LoadState from '@/components/ui/LoadState.vue'
import type { Article } from '@/types/article'
import type { ArticlePatchPayload } from '@/api/articles'
import { useArticlesStore } from '@/stores/articles'
import { useSectionsStore } from '@/stores/sections'
import { formatDate } from '@/utils/date'

const store = useArticlesStore()
const sectionsStore = useSectionsStore()

const errorMessage = ref('')
const busyId = ref<string | null>(null)

/** Section id → display title, resolved from the server section registry. */
function sectionTitle(article: Article): string {
  return article.section ? sectionsStore.titleOf(article.section) : ''
}

const drafts = computed(() =>
  store.all
    .filter((article) => article.status === 'draft')
    .sort((a, b) => b.date.localeCompare(a.date)),
)

const published = computed(() =>
  store.all
    .filter((article) => article.status === 'published')
    .sort((a, b) => b.date.localeCompare(a.date)),
)

const ready = computed(() => store.loadState === 'ready')

/**
 * Status switches travel with the full payload: the backend clears review
 * fields that are missing from a PATCH request.
 */
function buildStatusPatch(article: Article): ArticlePatchPayload {
  const status = article.status === 'published' ? 'draft' : 'published'

  const patch: ArticlePatchPayload = {
    title: article.title,
    image: article.image,
    kind: article.kind,
    content: article.content,
    status,
  }

  if (article.kind === 'review') {
    patch.rating = article.rating ?? 0
    patch.platform = article.platform ?? ''
    patch.verdict = article.verdict ?? ''
  }

  return patch
}

async function toggleStatus(article: Article) {
  if (busyId.value) return

  errorMessage.value = ''
  busyId.value = article.id

  try {
    await store.updateArticle(article.id, buildStatusPatch(article))
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Не удалось изменить статус.'
  } finally {
    busyId.value = null
  }
}

async function remove(article: Article) {
  if (busyId.value) return
  if (!window.confirm(`Удалить статью «${article.title}»? Это действие необратимо.`)) return

  errorMessage.value = ''
  busyId.value = article.id

  try {
    await store.deleteArticle(article.id)
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Не удалось удалить статью.'
  } finally {
    busyId.value = null
  }
}
</script>

<template>
  <div class="manage container">
    <header class="manage__header">
      <h1 class="manage__heading">Статьи</h1>
      <RouterLink class="manage__new" :to="{ name: 'editor' }">Новая статья</RouterLink>
    </header>

    <p v-if="errorMessage" class="manage__error" role="alert">{{ errorMessage }}</p>

    <LoadState
      v-if="!ready"
      class="manage__state"
      :error="store.errorMessage"
      @retry="store.reload()"
    />

    <template v-else>
      <section class="manage__section" aria-labelledby="manage-drafts-title">
        <h2 id="manage-drafts-title" class="manage__section-title">Черновики</h2>

        <p v-if="drafts.length === 0" class="manage__empty">Черновиков нет.</p>
        <ul v-else class="manage__list">
          <li v-for="article in drafts" :key="article.id" class="manage__item">
            <div class="manage__info">
              <RouterLink
                class="manage__title"
                :to="{ name: 'article', params: { id: article.id } }"
              >
                {{ article.title }}
              </RouterLink>
              <span class="manage__meta">
                {{ store.articleKindLabel(article.kind) }} · {{ article.author }} ·
                {{ formatDate(article.date) }}<template v-if="sectionTitle(article)">
                  · {{ sectionTitle(article) }}</template
                >
              </span>
            </div>

            <div class="manage__actions">
              <RouterLink
                class="manage__button"
                :to="{ name: 'editor-edit', params: { id: article.id } }"
              >
                Редактировать
              </RouterLink>
              <button
                class="manage__button manage__button--primary"
                type="button"
                :disabled="busyId !== null"
                @click="toggleStatus(article)"
              >
                Опубликовать
              </button>
              <button
                class="manage__button manage__button--danger"
                type="button"
                :disabled="busyId !== null"
                @click="remove(article)"
              >
                Удалить
              </button>
            </div>
          </li>
        </ul>
      </section>

      <section class="manage__section" aria-labelledby="manage-published-title">
        <h2 id="manage-published-title" class="manage__section-title">Опубликовано</h2>

        <p v-if="published.length === 0" class="manage__empty">Опубликованных статей нет.</p>
        <ul v-else class="manage__list">
          <li v-for="article in published" :key="article.id" class="manage__item">
            <div class="manage__info">
              <RouterLink
                class="manage__title"
                :to="{ name: 'article', params: { id: article.id } }"
              >
                {{ article.title }}
              </RouterLink>
              <span class="manage__meta">
                {{ store.articleKindLabel(article.kind) }} · {{ article.author }} ·
                {{ formatDate(article.date) }}<template v-if="sectionTitle(article)">
                  · {{ sectionTitle(article) }}</template
                >
              </span>
            </div>

            <div class="manage__actions">
              <RouterLink
                class="manage__button"
                :to="{ name: 'editor-edit', params: { id: article.id } }"
              >
                Редактировать
              </RouterLink>
              <button
                class="manage__button manage__button--muted"
                type="button"
                :disabled="busyId !== null"
                @click="toggleStatus(article)"
              >
                Снять с публикации
              </button>
              <button
                class="manage__button manage__button--danger"
                type="button"
                :disabled="busyId !== null"
                @click="remove(article)"
              >
                Удалить
              </button>
            </div>
          </li>
        </ul>
      </section>
    </template>
  </div>
</template>

<style scoped>
.manage {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  max-width: 1060px;
  padding-block: var(--gap);
  margin-inline: auto;
}

.manage__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 12px;
  background: var(--color-surface);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.manage__heading {
  font: 600 24px/30px Roboto, sans-serif;
}

.manage__new {
  display: inline-flex;
  align-items: center;
  min-height: 36px;
  padding: 6px 14px;
  color: var(--color-accent-text);
  background: var(--color-accent);
  border-radius: var(--radius);
  font: 500 13px/16px Roboto, sans-serif;
}

.manage__new:hover {
  background: var(--color-card-hover);
}

.manage__error {
  padding: 8px;
  color: var(--color-danger);
  font: 400 12px/16px Roboto, sans-serif;
}

.manage__state {
  width: 100%;
}

.manage__section {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 12px;
  background: var(--color-surface);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.manage__section-title {
  font: 600 16px/20px Roboto, sans-serif;
}

.manage__empty {
  padding: 8px;
  color: var(--color-muted);
  font: 400 13px/18px Roboto, sans-serif;
}

.manage__list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.manage__item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.manage__info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.manage__title {
  color: var(--color-text);
  font: 500 14px/18px Roboto, sans-serif;
}

.manage__title:hover {
  color: var(--color-accent);
}

.manage__meta {
  color: var(--color-muted);
  font: 400 12px/16px Roboto, sans-serif;
}

.manage__actions {
  display: flex;
  flex: none;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 6px;
}

.manage__button {
  display: inline-flex;
  align-items: center;
  min-height: 32px;
  padding: 4px 12px;
  color: var(--color-text);
  background: var(--color-surface);
  border: 0;
  border-radius: var(--radius);
  font: 500 12px/16px Roboto, sans-serif;
  cursor: pointer;
}

.manage__button:hover {
  background: var(--color-card-hover);
}

.manage__button:disabled {
  opacity: 0.6;
  cursor: default;
}

.manage__button--primary {
  color: var(--color-accent-text);
  background: var(--color-accent);
}

.manage__button--primary:hover {
  background: var(--color-accent);
  opacity: 0.9;
}

.manage__button--danger {
  color: var(--color-danger);
}

.manage__button--muted {
  color: var(--color-muted);
}

@media (max-width: 767px) {
  .manage__item {
    flex-direction: column;
    align-items: stretch;
  }

  .manage__actions {
    justify-content: flex-start;
  }
}
</style>
