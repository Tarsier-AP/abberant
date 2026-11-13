<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import type { Article } from '@/types/article'
import type { CommentPayload } from '@/api/comments'
import * as api from '@/api/comments'
import { useAuthStore } from '@/stores/auth'
import { useArticlesStore } from '@/stores/articles'
import { formatDate } from '@/utils/date'
import UserAvatar from '@/components/ui/UserAvatar.vue'
import LoadState from '@/components/ui/LoadState.vue'

const props = defineProps<{ article: Article }>()

const auth = useAuthStore()
const articlesStore = useArticlesStore()

const comments = ref<CommentPayload[]>([])
const loadState = ref<'loading' | 'ready' | 'error'>('loading')
const loadError = ref('')

const draft = ref('')
const replyDraft = ref('')
/** Id of the top-level comment the inline reply form is open for. */
const replyTo = ref<number | null>(null)
const submitting = ref(false)
const formError = ref('')
const replyError = ref('')

const account = computed(() => auth.account)

/** Replies are nested one level deep, so only top-level comments get the link. */
const canReply = (comment: CommentPayload) => account.value !== null && comment.parentId === null

const canDelete = (comment: CommentPayload) =>
  account.value !== null &&
  (auth.canModerate || comment.author.username === account.value.username)

/** Top-level comments, oldest first (the API order). */
const topLevel = computed(() => comments.value.filter((comment) => comment.parentId === null))

function repliesOf(id: number): CommentPayload[] {
  return comments.value.filter((comment) => comment.parentId === id)
}

/** Keeps the counter on article cards in sync with the section. */
function syncCount() {
  articlesStore.setCommentCount(props.article.id, comments.value.length)
}

async function load() {
  loadState.value = 'loading'
  loadError.value = ''

  try {
    const { comments: loaded } = await api.fetchComments(props.article.id)

    comments.value = loaded
    loadState.value = 'ready'
    syncCount()
  } catch (error) {
    loadState.value = 'error'
    loadError.value =
      error instanceof Error ? error.message : 'Не удалось загрузить комментарии.'
  }
}

onMounted(() => {
  void load()
})

function startReply(comment: CommentPayload) {
  replyTo.value = comment.id
  replyDraft.value = ''
  replyError.value = ''
}

function cancelReply() {
  replyTo.value = null
  replyDraft.value = ''
  replyError.value = ''
}

async function send(body: string, parentId: number | null, done: () => void) {
  if (submitting.value) return

  submitting.value = true
  formError.value = ''
  replyError.value = ''

  try {
    const payload = parentId === null ? { body } : { body, parentId }
    const { comment } = await api.createComment(props.article.id, payload)

    comments.value = [...comments.value, comment]
    syncCount()
    done()
  } catch (error) {
    const message =
      error instanceof Error ? error.message : 'Не удалось отправить комментарий.'

    if (parentId === null) {
      formError.value = message
    } else {
      replyError.value = message
    }
  } finally {
    submitting.value = false
  }
}

async function submitComment() {
  const body = draft.value.trim()

  if (body === '') {
    formError.value = 'Введите текст комментария.'
    return
  }

  await send(body, null, () => {
    draft.value = ''
  })
}

async function submitReply() {
  if (replyTo.value === null) return

  const body = replyDraft.value.trim()

  if (body === '') {
    replyError.value = 'Введите текст комментария.'
    return
  }

  const parentId = replyTo.value

  await send(body, parentId, () => {
    cancelReply()
  })
}

async function remove(comment: CommentPayload) {
  if (!window.confirm('Удалить этот комментарий?')) return

  try {
    await api.deleteComment(comment.id)

    // Replies are removed together with their parent on the server.
    comments.value = comments.value.filter(
      (item) => item.id !== comment.id && item.parentId !== comment.id,
    )
    syncCount()
  } catch (error) {
    formError.value =
      error instanceof Error ? error.message : 'Не удалось удалить комментарий.'
  }
}
</script>

<template>
  <section id="comments" class="comments-section" aria-label="Комментарии">
    <h2 class="comments-section__title">Комментарии ({{ comments.length }})</h2>

    <LoadState
      v-if="loadState !== 'ready'"
      class="comments-section__state"
      :error="loadError"
      @retry="load()"
    />

    <template v-else>
      <p v-if="topLevel.length === 0" class="comments-section__empty">
        Пока нет комментариев. Станьте первым!
      </p>

      <ul v-else class="comments-section__list">
        <li v-for="comment in topLevel" :key="comment.id" class="comment">
          <UserAvatar
            class="comment__avatar"
            :name="comment.author.username"
            :avatar="comment.author.avatar"
            :size="36"
          />

          <div class="comment__main">
            <div class="comment__head">
              <span class="comment__author">{{ comment.author.username }}</span>
              <time class="comment__date" :datetime="comment.createdAt">{{
                formatDate(comment.createdAt)
              }}</time>
            </div>

            <p class="comment__body">{{ comment.body }}</p>

            <div v-if="account" class="comment__actions">
              <button
                v-if="canReply(comment)"
                class="comment__action"
                type="button"
                @click="startReply(comment)"
              >
                Ответить
              </button>
              <button
                v-if="canDelete(comment)"
                class="comment__action comment__action--danger"
                type="button"
                @click="remove(comment)"
              >
                Удалить
              </button>
            </div>

            <form
              v-if="replyTo === comment.id"
              class="comment-form comment-form--reply"
              novalidate
              @submit.prevent="submitReply"
            >
              <textarea
                v-model="replyDraft"
                class="comment-form__input"
                rows="3"
                maxlength="500"
                placeholder="Ответить…"
                :disabled="submitting"
              ></textarea>
              <p v-if="replyError" class="comment-form__error">{{ replyError }}</p>
              <div class="comment-form__controls">
                <button
                  class="comment-form__submit"
                  type="submit"
                  :disabled="submitting || replyDraft.trim() === ''"
                >
                  Отправить
                </button>
                <button class="comment-form__cancel" type="button" @click="cancelReply">
                  Отмена
                </button>
              </div>
            </form>

            <ul v-if="repliesOf(comment.id).length > 0" class="comment__replies">
              <li v-for="reply in repliesOf(comment.id)" :key="reply.id" class="comment comment--reply">
                <UserAvatar
                  class="comment__avatar"
                  :name="reply.author.username"
                  :avatar="reply.author.avatar"
                  :size="28"
                />

                <div class="comment__main">
                  <div class="comment__head">
                    <span class="comment__author">{{ reply.author.username }}</span>
                    <time class="comment__date" :datetime="reply.createdAt">{{
                      formatDate(reply.createdAt)
                    }}</time>
                  </div>

                  <p class="comment__body">{{ reply.body }}</p>

                  <div v-if="canDelete(reply)" class="comment__actions">
                    <button
                      class="comment__action comment__action--danger"
                      type="button"
                      @click="remove(reply)"
                    >
                      Удалить
                    </button>
                  </div>
                </div>
              </li>
            </ul>
          </div>
        </li>
      </ul>

      <form v-if="account" class="comment-form" novalidate @submit.prevent="submitComment">
        <UserAvatar
          class="comment-form__avatar"
          :name="account.username"
          :avatar="account.avatar"
          :size="36"
        />

        <div class="comment-form__main">
          <textarea
            v-model="draft"
            class="comment-form__input"
            rows="3"
            maxlength="500"
            placeholder="Напишите комментарий…"
            :disabled="submitting"
          ></textarea>
          <p v-if="formError" class="comment-form__error">{{ formError }}</p>
          <div class="comment-form__controls">
            <button
              class="comment-form__submit"
              type="submit"
              :disabled="submitting || draft.trim() === ''"
            >
              Отправить
            </button>
          </div>
        </div>
      </form>
      <p v-else class="comments-section__guest">
        Комментировать могут только зарегистрированные пользователи.
        <RouterLink class="comments-section__guest-link" to="/login">Войти</RouterLink>
      </p>
    </template>
  </section>
</template>

<style scoped>
.comments-section {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: clamp(8px, 1.5vw, 12px);
  background: var(--color-header);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.comments-section__title {
  font: 600 16px/20px Roboto, sans-serif;
}

.comments-section__state,
.comments-section__empty {
  color: var(--color-muted);
  font: 400 14px/20px Roboto, sans-serif;
}

.comments-section__list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.comment {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 10px;
}

.comment__head {
  display: flex;
  flex-wrap: wrap;
  align-items: baseline;
  gap: 8px;
}

.comment__author {
  font: 500 13px/16px Roboto, sans-serif;
}

.comment__date {
  color: var(--color-muted);
  font: 400 12px/16px Roboto, sans-serif;
}

.comment__body {
  margin-top: 4px;
  font: 400 14px/20px Roboto, sans-serif;
  white-space: pre-line;
}

.comment__actions {
  display: flex;
  gap: 10px;
  margin-top: 4px;
}

.comment__action {
  padding: 0;
  color: var(--color-muted);
  background: none;
  border: none;
  cursor: pointer;
  font: 400 12px/16px Roboto, sans-serif;
  transition: color 0.15s ease;
}

.comment__action:hover,
.comment__action:focus-visible {
  color: var(--color-accent-text);
}

.comment__action--danger:hover,
.comment__action--danger:focus-visible {
  color: var(--color-warning);
}

.comment__replies {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 10px;
  padding-left: 12px;
  border-left: 2px solid var(--color-surface);
}

.comment-form {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 10px;
  align-items: start;
}

.comment-form--reply {
  grid-template-columns: 1fr;
  margin-top: 8px;
}

.comment-form__main {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.comment-form__input {
  width: 100%;
  padding: 10px;
  color: inherit;
  background: var(--color-surface);
  border: 1px solid transparent;
  border-radius: var(--radius);
  font: 400 14px/20px Roboto, sans-serif;
  resize: vertical;
}

.comment-form__input:focus {
  outline: none;
  border-color: var(--color-accent);
}

.comment-form__controls {
  display: flex;
  gap: 10px;
}

.comment-form__submit {
  padding: 8px 16px;
  color: var(--color-accent-text);
  background: var(--color-accent);
  border: none;
  border-radius: var(--radius);
  cursor: pointer;
  font: 500 13px/16px Roboto, sans-serif;
}

.comment-form__submit:disabled {
  opacity: 0.6;
  cursor: default;
}

.comment-form__cancel {
  padding: 8px 12px;
  color: var(--color-muted);
  background: none;
  border: none;
  cursor: pointer;
  font: 400 13px/16px Roboto, sans-serif;
}

.comment-form__cancel:hover {
  color: var(--color-accent-text);
}

.comment-form__error {
  color: var(--color-warning);
  font: 400 13px/16px Roboto, sans-serif;
}

.comments-section__guest {
  color: var(--color-muted);
  font: 400 15px/22px Roboto, sans-serif;
}

.comments-section__guest-link {
  color: var(--color-accent-text);
  font: 500 15px/22px Roboto, sans-serif;
}

.comments-section__guest-link:hover {
  text-decoration: underline;
}
</style>
