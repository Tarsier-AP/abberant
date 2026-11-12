<script setup lang="ts">
import type { EditorCommand } from '@/types/editor'

const buttons: Array<{ command: EditorCommand; label: string; text: string }> = [
  { command: 'h2', label: 'Заголовок второго уровня', text: 'H2' },
  { command: 'h3', label: 'Заголовок третьего уровня', text: 'H3' },
  { command: 'bold', label: 'Полужирный', text: 'Ж' },
  { command: 'italic', label: 'Курсив', text: 'К' },
  { command: 'strike', label: 'Зачёркнутый', text: 'Ч' },
  { command: 'quote', label: 'Цитата', text: '«' },
  { command: 'ul', label: 'Маркированный список', text: '•' },
  { command: 'ol', label: 'Нумерованный список', text: '1.' },
  { command: 'link', label: 'Ссылка', text: '🔗' },
  { command: 'image', label: 'Изображение', text: '🖼' },
  { command: 'code', label: 'Код', text: '<>' },
  { command: 'hr', label: 'Разделитель', text: '—' },
]

defineEmits<{ exec: [command: EditorCommand] }>()
</script>

<template>
  <div class="toolbar" role="toolbar" aria-label="Форматирование статьи">
    <button
      v-for="button in buttons"
      :key="button.command"
      class="toolbar__button"
      type="button"
      :title="button.label"
      :aria-label="button.label"
      @click="$emit('exec', button.command)"
    >
      {{ button.text }}
    </button>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  padding: 6px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.toolbar__button {
  display: grid;
  place-items: center;
  min-width: 32px;
  min-height: 32px;
  padding: 0 6px;
  background: var(--color-surface);
  border: 0;
  border-radius: var(--radius);
  font: 500 13px/16px Roboto, sans-serif;
  cursor: pointer;
}

.toolbar__button:hover {
  background: var(--color-card-hover);
}
</style>
