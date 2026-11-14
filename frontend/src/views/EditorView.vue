<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { onBeforeRouteLeave, useRoute, useRouter } from 'vue-router'
import EditorToolbar from '@/components/editor/EditorToolbar.vue'
import EditorImageDialog from '@/components/editor/EditorImageDialog.vue'
import type { EditorCommand, EditorImagePayload } from '@/types/editor'
import type { ArticleKind, ArticleStatus } from '@/types/article'
import type { ArticleCreatePayload, ArticlePatchPayload } from '@/api/articles'
import { useAuthStore } from '@/stores/auth'
import { useArticlesStore } from '@/stores/articles'
import { newsPlatformSections, reviewGroupSections, walkthroughSections } from '@/data/sections'
import { uploadImage } from '@/api/media'
import { renderMarkdown } from '@/utils/markdown'

const FALLBACK_IMAGE = '/img/news-placeholder.svg'

/** Section picker groups; ids go into the article's `section` column. */
const sectionGroups = [
  { label: 'Новости', sections: newsPlatformSections },
  { label: 'Прохождения', sections: walkthroughSections },
  { label: 'Обзоры', sections: reviewGroupSections },
]

interface EditorForm {
  title: string
  kind: ArticleKind
  image: string
  content: string
  rating: number
  platform: string
  verdict: string
  status: ArticleStatus
  section: string
  date: string
}

/** Local calendar date as YYYY-MM-DD for the publication date input. */
function todayIsoDate(): string {
  const now = new Date()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')

  return `${now.getFullYear()}-${month}-${day}`
}

function emptyForm(): EditorForm {
  return {
    title: '',
    kind: 'news',
    image: '',
    content: '',
    rating: 3,
    platform: '',
    verdict: '',
    status: 'draft',
    section: '',
    date: todayIsoDate(),
  }
}

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const store = useArticlesStore()

const textarea = ref<HTMLTextAreaElement | null>(null)
const coverInput = ref<HTMLInputElement | null>(null)

const form = reactive<EditorForm>(emptyForm())

/** Identity fields that the editor form does not change. */
const editingId = ref<string | null>(null)
const author = ref('')
const createdDate = ref('')
const comments = ref(0)

const coverError = ref('')
const errorMessage = ref('')
const savedNotice = ref('')
const imageDialogOpen = ref(false)
const saving = ref(false)

const isNew = computed(() => editingId.value === null)
/** Owner and moderators manage any article; others may only touch their own. */
const isManager = computed(
  () => auth.account?.role === 'owner' || auth.account?.role === 'moderator',
)
const canModify = computed(
  () => isManager.value || (author.value !== '' && author.value === auth.account?.username),
)
const canDelete = canModify

const previewHtml = computed(() => renderMarkdown(form.content))

const snapshot = ref('')

function formSnapshot(): string {
  return JSON.stringify({
    title: form.title,
    kind: form.kind,
    image: form.image,
    content: form.content,
    rating: form.rating,
    platform: form.platform,
    verdict: form.verdict,
    status: form.status,
    section: form.section,
    date: form.date,
  })
}

const dirty = computed(() => formSnapshot() !== snapshot.value)

function takeSnapshot() {
  snapshot.value = formSnapshot()
}

function initForm() {
  if (store.loadState !== 'ready') return

  errorMessage.value = ''
  savedNotice.value = ''
  coverError.value = ''

  const id = typeof route.params.id === 'string' ? route.params.id : null

  if (!id) {
    editingId.value = null
    author.value = auth.account?.username ?? ''
    createdDate.value = ''
    comments.value = 0

    Object.assign(form, emptyForm())
    takeSnapshot()

    return
  }

  const article = store.getArticle(id)

  if (!article) {
    void router.replace({ name: 'editor' })
    return
  }

  if (!isManager.value && article.author !== auth.account?.username) {
    void router.replace({ name: 'article', params: { id } })
    return
  }

  editingId.value = article.id
  author.value = article.author
  createdDate.value = article.date
  comments.value = article.comments

  Object.assign(form, emptyForm(), {
    title: article.title,
    kind: article.kind,
    image: article.image === FALLBACK_IMAGE ? '' : article.image,
    content: article.content,
    rating: article.rating ?? 3,
    platform: article.platform ?? '',
    verdict: article.verdict ?? '',
    status: article.status,
    section: article.section ?? '',
    date: article.date,
  })
  takeSnapshot()
}

function buildCreatePayload(): ArticleCreatePayload {
  const payload: ArticleCreatePayload = {
    title: form.title.trim(),
    image: form.image.trim() || FALLBACK_IMAGE,
    kind: form.kind,
    content: form.content,
    status: form.status,
    date: form.date,
    section: form.section,
  }

  if (form.kind === 'review') {
    payload.rating = form.rating
    payload.platform = form.platform.trim()
    payload.verdict = form.verdict.trim()
  }

  return payload
}

function buildPatchPayload(): ArticlePatchPayload {
  return {
    title: form.title.trim(),
    image: form.image.trim() || FALLBACK_IMAGE,
    kind: form.kind,
    content: form.content,
    status: form.status,
    date: form.date,
    section: form.section,
    ...(form.kind === 'review'
      ? {
          rating: form.rating,
          platform: form.platform.trim(),
          verdict: form.verdict.trim(),
        }
      : {}),
  }
}

function validate(): string {
  if (!form.title.trim()) return 'Введите заголовок статьи.'

  if (form.status === 'published' && !form.content.trim()) {
    return 'Добавьте текст статьи перед публикацией.'
  }

  return ''
}

async function save(nextStatus?: ArticleStatus) {
  if (saving.value) return

  errorMessage.value = ''
  savedNotice.value = ''

  if (nextStatus) form.status = nextStatus

  const error = validate()

  if (error) {
    errorMessage.value = error
    return
  }

  saving.value = true

  try {
    const article =
      editingId.value === null
        ? await store.createArticle(buildCreatePayload())
        : await store.updateArticle(editingId.value, buildPatchPayload())

    editingId.value = article.id
    author.value = article.author
    createdDate.value = article.date
    comments.value = article.comments

    if (nextStatus === 'published') {
      void router.push({ name: 'article', params: { id: article.id } })
      return
    }

    if (route.name === 'editor') {
      void router.replace({ name: 'editor-edit', params: { id: article.id } })
    }

    takeSnapshot()
    savedNotice.value = form.status === 'published' ? 'Статья сохранена.' : 'Черновик сохранён.'
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Не удалось сохранить статью.'
  } finally {
    saving.value = false
  }
}

async function remove() {
  if (!editingId.value || !window.confirm('Удалить статью? Это действие необратимо.')) return

  try {
    await store.deleteArticle(editingId.value)
  } catch (err) {
    errorMessage.value = err instanceof Error ? err.message : 'Не удалось удалить статью.'

    return
  }

  void router.replace({ name: 'editor-articles' })
}

function goBack() {
  const target = editingId.value ? { name: 'article', params: { id: editingId.value } } : { name: 'home' }

  if (dirty.value && !window.confirm('Несохранённые изменения будут потеряны. Покинуть страницу?')) {
    return
  }

  void router.push(target)
}

function handleCoverChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  input.value = ''
  coverError.value = ''

  if (!file) return

  if (!file.type.startsWith('image/')) {
    coverError.value = 'Выберите файл изображения (JPG, PNG или WebP).'
    return
  }

  if (file.size > 10 * 1024 * 1024) {
    coverError.value = 'Файл слишком большой. Максимальный размер — 10 МБ.'
    return
  }

  uploadImage(file, 'cover')
    .then(({ url }) => {
      form.image = url
      coverError.value = ''
    })
    .catch((uploadError) => {
      coverError.value =
        uploadError instanceof Error && uploadError.message
          ? uploadError.message
          : 'Не удалось обработать изображение. Попробуйте другой файл.'
    })
}

function wrapSelection(before: string, after: string, placeholder: string) {
  const element = textarea.value

  if (!element) return

  const start = element.selectionStart ?? 0
  const end = element.selectionEnd ?? 0
  const selected = form.content.slice(start, end)
  const text = selected || placeholder

  form.content = form.content.slice(0, start) + before + text + after + form.content.slice(end)

  void nextTick(() => {
    element.focus()
    element.setSelectionRange(start + before.length, start + before.length + text.length)
  })
}

function prefixLines(prefix: string) {
  const element = textarea.value

  if (!element) return

  const start = element.selectionStart ?? 0
  const end = element.selectionEnd ?? 0
  const lineStart = form.content.lastIndexOf('\n', Math.max(0, start - 1)) + 1
  let lineEnd = form.content.indexOf('\n', end)

  if (lineEnd === -1) lineEnd = form.content.length

  const block = form.content.slice(lineStart, lineEnd)
  const prefixed = block
    .split('\n')
    .map((line) => (line.trim() ? prefix + line : line))
    .join('\n')

  form.content = form.content.slice(0, lineStart) + prefixed + form.content.slice(lineEnd)

  void nextTick(() => {
    element.focus()
    element.setSelectionRange(lineStart + prefix.length, lineStart + prefix.length)
  })
}

function insertSnippet(snippet: string) {
  const element = textarea.value

  if (!element) return

  const start = element.selectionStart ?? 0

  form.content = form.content.slice(0, start) + snippet + form.content.slice(start)

  void nextTick(() => {
    element.focus()
    element.setSelectionRange(start + snippet.length, start + snippet.length)
  })
}

function insertImage(payload: EditorImagePayload) {
  imageDialogOpen.value = false

  const element = textarea.value
  const start = element ? (element.selectionStart ?? form.content.length) : form.content.length
  const before = start > 0 && form.content[start - 1] !== '\n' ? '\n' : ''
  const alt = payload.alt || 'Изображение'
  const caption = payload.caption ? ` "${payload.caption}"` : ''
  const snippet = `${before}![${alt}](${payload.url}${caption})\n`

  form.content = form.content.slice(0, start) + snippet + form.content.slice(start)

  void nextTick(() => element?.focus())
}

function applyCommand(command: EditorCommand) {
  switch (command) {
    case 'h2':
      prefixLines('## ')
      break
    case 'h3':
      prefixLines('### ')
      break
    case 'bold':
      wrapSelection('**', '**', 'текст')
      break
    case 'italic':
      wrapSelection('*', '*', 'текст')
      break
    case 'strike':
      wrapSelection('~~', '~~', 'текст')
      break
    case 'quote':
      prefixLines('> ')
      break
    case 'ul':
      prefixLines('- ')
      break
    case 'ol':
      prefixLines('1. ')
      break
    case 'link':
      wrapSelection('[', '](https://example.com)', 'Текст ссылки')
      break
    case 'image':
      imageDialogOpen.value = true
      break
    case 'code':
      wrapSelection('`', '`', 'код')
      break
    case 'hr':
      insertSnippet('\n\n---\n\n')
      break
  }
}

function handleBeforeUnload(event: BeforeUnloadEvent) {
  if (dirty.value) {
    event.preventDefault()
    event.returnValue = ''
  }
}

watch(
  () => route.params.id,
  () => initForm(),
)

watch(
  () => store.loadState,
  (state) => {
    if (state === 'ready') initForm()
  },
)

onMounted(() => {
  initForm()
  window.addEventListener('beforeunload', handleBeforeUnload)
})

onBeforeUnmount(() => {
  window.removeEventListener('beforeunload', handleBeforeUnload)
})

onBeforeRouteLeave(() => {
  if (dirty.value && !window.confirm('Несохранённые изменения будут потеряны. Покинуть страницу?')) {
    return false
  }

  return true
})
</script>

<template>
  <div class="editor container">
    <h1 class="editor__heading">{{ isNew ? 'Новая статья' : 'Редактирование статьи' }}</h1>

    <p v-if="savedNotice" class="editor__notice" role="status">{{ savedNotice }}</p>
    <p v-if="errorMessage" class="editor__error" role="alert">{{ errorMessage }}</p>

    <section class="editor__section" aria-labelledby="editor-fields-title">
      <h2 id="editor-fields-title" class="editor__section-title">Основное</h2>

      <div class="editor__field">
        <label class="editor__label" for="editor-title">Заголовок</label>
        <input
          id="editor-title"
          v-model="form.title"
          class="editor__input"
          type="text"
          placeholder="О чём статья"
        />
      </div>

      <div class="editor__field">
        <span class="editor__label" id="editor-kind-label">Тип материала</span>
        <div class="editor__kinds" role="radiogroup" aria-labelledby="editor-kind-label">
          <button
            class="editor__kind"
            :class="{ 'editor__kind--active': form.kind === 'news' }"
            type="button"
            :aria-pressed="form.kind === 'news'"
            @click="form.kind = 'news'"
          >
            Новость
          </button>
          <button
            class="editor__kind"
            :class="{ 'editor__kind--active': form.kind === 'review' }"
            type="button"
            :aria-pressed="form.kind === 'review'"
            @click="form.kind = 'review'"
          >
            Обзор
          </button>
        </div>
      </div>

      <div class="editor__field">
        <label class="editor__label" for="editor-date">Дата публикации</label>
        <input
          id="editor-date"
          v-model="form.date"
          class="editor__input"
          type="date"
        />
      </div>

      <div class="editor__field">
        <label class="editor__label" for="editor-section">Раздел</label>
        <select id="editor-section" v-model="form.section" class="editor__input editor__select">
          <option value="">Без раздела</option>
          <optgroup
            v-for="group in sectionGroups"
            :key="group.label"
            :label="group.label"
          >
            <option v-for="section in group.sections" :key="section.id" :value="section.id">
              {{ section.title }}
            </option>
          </optgroup>
        </select>
      </div>

      <div class="editor__field">
        <label class="editor__label" for="editor-cover-url">Обложка</label>
        <div class="editor__cover">
          <img
            v-if="form.image"
            class="editor__cover-preview"
            :src="form.image"
            alt="Предпросмотр обложки"
          />
          <div v-else class="editor__cover-empty">Нет обложки</div>

          <div class="editor__cover-controls">
            <input
              id="editor-cover-url"
              v-model="form.image"
              class="editor__input"
              type="text"
              placeholder="Ссылка на изображение"
            />
            <div class="editor__cover-buttons">
              <button class="editor__button" type="button" @click="coverInput?.click()">
                Загрузить файл
              </button>
              <button
                v-if="form.image"
                class="editor__button editor__button--muted"
                type="button"
                @click="form.image = ''"
              >
                Убрать
              </button>
            </div>
          </div>

          <input ref="coverInput" hidden type="file" accept="image/*" @change="handleCoverChange" />
        </div>
        <p v-if="coverError" class="editor__error" role="alert">{{ coverError }}</p>
      </div>

      <template v-if="form.kind === 'review'">
        <div class="editor__field">
          <span class="editor__label" id="editor-rating-label">Оценка</span>
          <div class="editor__rating" role="radiogroup" aria-labelledby="editor-rating-label">
            <button
              v-for="star in 5"
              :key="star"
              class="editor__star"
              :class="{ 'editor__star--filled': star <= form.rating }"
              type="button"
              :aria-label="`Оценка ${star} из 5`"
              :aria-pressed="form.rating === star"
              @click="form.rating = star"
            >
              ★
            </button>
            <span class="editor__rating-value">{{ form.rating }} из 5</span>
          </div>
        </div>

        <div class="editor__field">
          <label class="editor__label" for="editor-platform">Платформы</label>
          <input
            id="editor-platform"
            v-model="form.platform"
            class="editor__input"
            type="text"
            placeholder="PC, PlayStation 5, Xbox Series X|S"
          />
        </div>

        <div class="editor__field">
          <label class="editor__label" for="editor-verdict">Вердикт</label>
          <input
            id="editor-verdict"
            v-model="form.verdict"
            class="editor__input"
            type="text"
            placeholder="Короткий итог обзора"
          />
        </div>
      </template>
    </section>

    <section class="editor__section" aria-labelledby="editor-content-title">
      <h2 id="editor-content-title" class="editor__section-title">Текст статьи</h2>

      <EditorToolbar @exec="applyCommand" />

      <div class="editor__split">
        <div class="editor__pane">
          <label class="visually-hidden" for="editor-content">Текст статьи в формате Markdown</label>
          <textarea
            id="editor-content"
            ref="textarea"
            v-model="form.content"
            class="editor__source"
            placeholder="Пишите текст в формате Markdown…"
          ></textarea>
        </div>

        <div class="editor__pane">
          <span class="editor__pane-label">Предпросмотр</span>
          <div v-if="form.content.trim()" class="article-body editor__preview" v-html="previewHtml"></div>
          <p v-else class="editor__preview-empty">Начните писать текст, и он появится здесь.</p>
        </div>
      </div>
    </section>

    <section class="editor__section" aria-labelledby="editor-actions-title">
      <h2 id="editor-actions-title" class="visually-hidden">Действия</h2>

      <div class="editor__actions">
        <button
          v-if="canModify"
          :disabled="saving"
          class="editor__button"
          type="button"
          @click="save()"
        >
          {{ saving ? 'Сохранение…' : 'Сохранить' }}
        </button>
        <button
          v-if="canModify && form.status === 'draft'"
          :disabled="saving"
          class="editor__button editor__button--primary"
          type="button"
          @click="save('published')"
        >
          {{ saving ? 'Публикация…' : 'Опубликовать' }}
        </button>
        <span v-else-if="canModify" class="editor__status">Опубликовано</span>
        <button
          v-if="canDelete && !isNew"
          class="editor__button editor__button--danger"
          type="button"
          @click="remove"
        >
          Удалить
        </button>
        <button class="editor__button editor__button--muted" type="button" @click="goBack">
          {{ isNew ? 'Отмена' : 'К статье' }}
        </button>
      </div>

      <p v-if="!canModify" class="editor__notice">
        У вас нет прав на редактирование этой статьи.
      </p>
    </section>

    <EditorImageDialog
      :open="imageDialogOpen"
      @cancel="imageDialogOpen = false"
      @confirm="insertImage"
    />
  </div>
</template>

<style scoped>
.editor {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  padding-block: var(--gap);
}

.editor__heading {
  font: 600 24px/30px Roboto, sans-serif;
}

.editor__section {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 12px;
  background: var(--color-surface);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.editor__section-title {
  font: 600 16px/20px Roboto, sans-serif;
}

.editor__field {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.editor__label {
  font: 500 13px/16px Roboto, sans-serif;
}

.editor__input {
  width: 100%;
  min-height: 34px;
  padding: 4px 8px;
  color: var(--color-text);
  background: var(--color-auth-input);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius);
  font: 400 14px/18px Roboto, sans-serif;
}

.editor__input:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 0;
}

.editor__select {
  appearance: auto;
}

.editor__kinds {
  display: flex;
  gap: 6px;
}

.editor__kind {
  padding: 6px 12px;
  background: var(--color-surface);
  border: 0;
  border-radius: var(--radius);
  font: 400 13px/16px Roboto, sans-serif;
  cursor: pointer;
}

.editor__kind--active {
  color: var(--color-accent-text);
  background: var(--color-accent);
}

.editor__cover {
  display: flex;
  gap: 10px;
  align-items: flex-start;
}

.editor__cover-preview {
  width: 160px;
  aspect-ratio: 330 / 208;
  object-fit: cover;
  border-radius: var(--radius);
}

.editor__cover-empty {
  display: grid;
  place-items: center;
  width: 160px;
  aspect-ratio: 330 / 208;
  color: var(--color-muted);
  background: var(--color-surface);
  border-radius: var(--radius);
  font: 400 12px/16px Roboto, sans-serif;
}

.editor__cover-controls {
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.editor__cover-buttons {
  display: flex;
  gap: 6px;
}

.editor__rating {
  display: flex;
  align-items: center;
  gap: 4px;
}

.editor__star {
  padding: 0;
  color: var(--color-muted);
  background: none;
  border: 0;
  font-size: 22px;
  line-height: 1;
  cursor: pointer;
}

.editor__star--filled {
  color: var(--color-accent);
}

.editor__rating-value {
  margin-left: 6px;
  color: var(--color-muted);
  font: 400 12px/16px Roboto, sans-serif;
}

.editor__split {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.editor__pane {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 8px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  min-width: 0;
}

.editor__pane-label {
  font: 500 12px/16px Roboto, sans-serif;
  color: var(--color-muted);
}

.editor__source {
  width: 100%;
  min-height: 460px;
  padding: 8px;
  color: var(--color-text);
  background: var(--color-auth-input);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius);
  font: 400 14px/20px 'Roboto Mono', ui-monospace, monospace;
  resize: vertical;
}

.editor__source:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 0;
}

.editor__preview {
  min-height: 460px;
  padding: 8px;
  background: var(--color-auth-input);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius);
}

.editor__preview-empty {
  min-height: 460px;
  padding: 8px;
  color: var(--color-muted);
  font: 400 13px/18px Roboto, sans-serif;
}

.editor__actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.editor__button {
  min-height: 36px;
  padding: 6px 14px;
  color: var(--color-text);
  background: var(--color-card);
  border: 0;
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  font: 500 13px/16px Roboto, sans-serif;
  cursor: pointer;
}

.editor__button:hover {
  background: var(--color-card-hover);
}

.editor__button--primary {
  color: var(--color-accent-text);
  background: var(--color-accent);
}

.editor__button--danger {
  color: var(--color-danger);
}

.editor__button--muted {
  color: var(--color-muted);
}

.editor__status {
  color: var(--color-muted);
  font: 500 13px/16px Roboto, sans-serif;
}

.editor__notice {
  padding: 8px;
  color: var(--color-warning);
  background: var(--color-card);
  border-radius: var(--radius);
  font: 600 12px/16px Roboto, sans-serif;
}

.editor__error {
  color: var(--color-danger);
  font: 400 12px/16px Roboto, sans-serif;
}

@media (max-width: 1023px) {
  .editor__split {
    grid-template-columns: 1fr;
  }

  .editor__source,
  .editor__preview,
  .editor__preview-empty {
    min-height: 280px;
  }
}

@media (max-width: 767px) {
  .editor__cover {
    flex-direction: column;
  }
}
</style>
