<script setup lang="ts">
import { useNotificationsStore } from '@/stores/notifications'

const store = useNotificationsStore()

defineProps<{ heading?: string }>()

defineExpose({ store })

function itemUrl(item: (typeof store.items)[number]): string {
  return item.article ? `/article/${item.article.id}` : '/'
}

function itemText(item: (typeof store.items)[number]): string {
  if (item.type === 'reply') return `${item.actor?.username ?? 'Читатель'} ответил на ваш комментарий`
  if (item.type === 'article') return 'Новый материал в вашем разделе'

  return 'Уведомление'
}
</script>

<template>
  <div class="panel">
    <header v-if="heading" class="panel__head">
      <span class="panel__title">{{ heading }}</span>
    </header>

    <p v-if="!store.items.length" class="panel__empty">
      Пока ничего не произошло. Ответы на ваши комментарии и новые материалы в
      ваших разделах появятся здесь.
    </p>

    <ul v-else class="panel__list">
      <li
        v-for="item in store.items"
        :key="item.id"
        class="panel__item"
        :class="{ 'panel__item--unread': !item.read }"
      >
        <RouterLink class="panel__link" :to="itemUrl(item)">
          <span class="panel__text">{{ itemText(item) }}</span>
          <span v-if="item.article" class="panel__article">{{ item.article.title }}</span>
        </RouterLink>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.panel {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.panel__head {
  padding-bottom: 8px;
  border-bottom: 1px solid var(--color-surface);
}

.panel__title {
  font: 600 15px/18px Roboto, sans-serif;
}

.panel__empty {
  padding-block: 12px;
  color: var(--color-muted);
  font: 400 13px/18px Roboto, sans-serif;
}

.panel__list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  overflow-y: auto;
  margin-block: 4px;
}

.panel__link {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 8px 10px;
  border-radius: var(--radius);
}

.panel__link:hover {
  background: var(--color-card);
}

.panel__item--unread .panel__text {
  font-weight: 600;
}

.panel__text {
  font: 400 13px/17px Roboto, sans-serif;
}

.panel__article {
  color: var(--color-muted);
  font: 400 12px/15px Roboto, sans-serif;
}
</style>
