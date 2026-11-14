<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import { uploadImage } from '@/api/media'
import type { EditorImagePayload } from '@/types/editor'

const props = defineProps<{ open: boolean }>()

const emit = defineEmits<{
  cancel: []
  confirm: [payload: EditorImagePayload]
}>()

const MAX_FILE_BYTES = 10 * 1024 * 1024

const url = ref('')
const alt = ref('')
const caption = ref('')
const error = ref('')
const fileInput = ref<HTMLInputElement | null>(null)
const urlInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)

watch(
  () => props.open,
  (open) => {
    if (!open) return

    url.value = ''
    alt.value = ''
    caption.value = ''
    error.value = ''
    uploading.value = false

    void nextTick(() => urlInput.value?.focus())
  },
)

function cancel() {
  emit('cancel')
}

function confirm() {
  const source = url.value.trim()

  if (!source) {
    error.value = 'Укажите ссылку на изображение или выберите файл.'
    return
  }

  emit('confirm', { url: source, alt: alt.value.trim(), caption: caption.value.trim() })
}

async function handleFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  input.value = ''
  error.value = ''

  if (!file) return

  if (!file.type.startsWith('image/')) {
    error.value = 'Выберите файл изображения (JPG, PNG или WebP).'
    return
  }

  if (file.size > MAX_FILE_BYTES) {
    error.value = 'Файл слишком большой. Максимальный размер — 10 МБ.'
    return
  }

  uploading.value = true

  try {
    const { url: uploaded } = await uploadImage(file, 'body')

    url.value = uploaded
    error.value = ''
  } catch (uploadError) {
    error.value =
      uploadError instanceof Error && uploadError.message
        ? uploadError.message
        : 'Не удалось обработать изображение. Попробуйте другой файл.'
  } finally {
    uploading.value = false
  }
}
</script>

<template>
  <div v-if="open" class="dialog" role="dialog" aria-modal="true" aria-labelledby="editor-image-title">
    <div class="dialog__backdrop" aria-hidden="true" @click="cancel"></div>

    <form class="dialog__card" novalidate @submit.prevent="confirm">
      <h2 id="editor-image-title" class="dialog__title">Вставить изображение</h2>

      <label class="dialog__label" for="editor-image-url">Ссылка на изображение</label>
      <input
        id="editor-image-url"
        ref="urlInput"
        v-model="url"
        class="dialog__input"
        type="text"
        placeholder="https://example.com/screenshot.jpg"
        :disabled="uploading"
      />

      <button
        class="dialog__upload"
        type="button"
        :disabled="uploading"
        @click="fileInput?.click()"
      >
        {{ uploading ? 'Обработка…' : 'Или загрузить файл' }}
      </button>

      <input
        ref="fileInput"
        hidden
        type="file"
        accept="image/*"
        @change="handleFileChange"
      />

      <label class="dialog__label" for="editor-image-alt">Описание (alt)</label>
      <input
        id="editor-image-alt"
        v-model="alt"
        class="dialog__input"
        type="text"
        placeholder="Что изображено на картинке"
      />

      <label class="dialog__label" for="editor-image-caption">Подпись под изображением</label>
      <input
        id="editor-image-caption"
        v-model="caption"
        class="dialog__input"
        type="text"
        placeholder="Необязательно"
      />

      <p v-if="error" class="dialog__error" role="alert">{{ error }}</p>

      <div class="dialog__actions">
        <button class="dialog__button" type="button" @click="cancel">Отмена</button>
        <button class="dialog__button" type="submit">Вставить</button>
      </div>
    </form>
  </div>
</template>

<style scoped>
.dialog {
  position: fixed;
  inset: 0;
  z-index: 50;
  display: grid;
  place-items: center;
  padding: 16px;
}

.dialog__backdrop {
  position: absolute;
  inset: 0;
  background: rgb(0 0 0 / 50%);
}

.dialog__card {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: min(420px, 100%);
  padding: 12px;
  background: var(--color-surface);
  border-radius: var(--radius);
  box-shadow: var(--shadow-hover);
}

.dialog__title {
  margin-bottom: 4px;
  font: 600 16px/20px Roboto, sans-serif;
}

.dialog__label {
  font: 500 12px/16px Roboto, sans-serif;
}

.dialog__input {
  width: 100%;
  min-height: 34px;
  padding: 4px 8px;
  color: var(--color-text);
  background: var(--color-auth-input);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius);
  font: 400 14px/18px Roboto, sans-serif;
}

.dialog__input:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 0;
}

.dialog__upload {
  align-self: flex-start;
  padding: 0;
  color: var(--color-auth-link);
  background: none;
  border: 0;
  font: 400 12px/16px Roboto, sans-serif;
  text-decoration: underline;
  text-underline-offset: 2px;
  cursor: pointer;
}

.dialog__upload:disabled {
  opacity: 0.6;
  cursor: default;
}

.dialog__error {
  color: var(--color-danger);
  font: 400 12px/16px Roboto, sans-serif;
}

.dialog__actions {
  display: flex;
  gap: 6px;
  margin-top: 6px;
}

.dialog__button {
  flex: 1;
  min-height: 34px;
  padding: 4px 8px;
  color: var(--color-text);
  background: var(--color-card);
  border: 0;
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  font: 400 13px/16px Roboto, sans-serif;
  cursor: pointer;
}

.dialog__button:hover {
  background: var(--color-card-hover);
}
</style>
