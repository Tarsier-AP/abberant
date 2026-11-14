<script setup lang="ts">
import { computed, onMounted } from 'vue'
import ArticleCard from '@/components/article/ArticleCard.vue'
import LoadState from '@/components/ui/LoadState.vue'
import { useBookmarksStore } from '@/stores/bookmarks'

const store = useBookmarksStore()

onMounted(() => {
  void store.loadList()
})

/** Mirrors the live id set, so removing a bookmark hides the card at once. */
const visible = computed(() => store.list.filter((article) => store.isBookmarked(article.id)))
</script>

<template>
  <section class="bookmarks container">
    <header class="bookmarks__header">
      <h1 class="bookmarks__title">Закладки</h1>
      <p class="bookmarks__hint">
        Статьи, которые вы сохранили, чтобы вернуться к ним позже.
      </p>
    </header>

    <LoadState
      v-if="store.listState === 'loading' || store.listState === 'error'"
      class="bookmarks__state"
      :error="store.listError"
      @retry="store.loadList()"
    />

    <div v-else-if="visible.length" class="bookmarks__grid">
      <ArticleCard v-for="article in visible" :key="article.id" :article="article" />
    </div>

    <p v-else class="bookmarks__empty">
      Пока пусто. Сохраняйте статьи кнопкой закладки на карточке или на странице
      статьи.
    </p>
  </section>
</template>

<style scoped>
.bookmarks {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  padding-block: var(--gap);
}

.bookmarks__header {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.bookmarks__title {
  font: 600 clamp(20px, 2.6vw, 24px)/28px Roboto, sans-serif;
}

.bookmarks__hint {
  color: var(--color-muted);
  font: 400 14px/18px Roboto, sans-serif;
}

/* Same grid as the archive pages: three columns on desktop, no row cap. */
.bookmarks__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--card-gap);
}

@media (max-width: 1023px) {
  .bookmarks__grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 767px) {
  .bookmarks__grid {
    grid-template-columns: 1fr;
  }
}

.bookmarks__empty {
  padding: clamp(12px, 2vw, 16px);
  color: var(--color-muted);
  background: var(--color-surface);
  border-radius: var(--radius);
  font: 400 14px/20px Roboto, sans-serif;
}
</style>
