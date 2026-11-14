<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import type { ReaderRating } from '@/types/article'
import { useAuthStore } from '@/stores/auth'
import { useArticlesStore } from '@/stores/articles'

const props = defineProps<{
  articleId: string
  rating: ReaderRating
}>()

const auth = useAuthStore()
const store = useArticlesStore()
const router = useRouter()

const busy = ref(false)
const hovered = ref(0)

const myVote = computed(() => props.rating.myVote ?? 0)
const average = computed(() => props.rating.average)
const votes = computed(() => props.rating.votes)

/** Stars the user sees as filled: hover preview or the own vote. */
const filledUpTo = computed(() => (auth.account && hovered.value ? hovered.value : myVote.value))

/** One decimal, comma as the decimal separator for the Russian UI. */
const averageLabel = computed(() =>
  votes.value ? average.value.toFixed(1).replace('.', ',') : '—',
)

const votesLabel = computed(() => {
  const mod10 = votes.value % 10
  const mod100 = votes.value % 100

  if (votes.value === 0) return 'нет оценок'

  if (mod10 === 1 && mod100 !== 11) return `${votes.value} голос`
  if (mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14)) return `${votes.value} голоса`

  return `${votes.value} голосов`
})

async function send(stars: number) {
  if (!auth.account) {
    void router.push({ name: 'login' })
    return
  }

  if (busy.value) return

  busy.value = true

  try {
    await store.voteArticle(props.articleId, stars)
  } finally {
    busy.value = false
  }
}

async function remove() {
  if (!auth.account || busy.value || !myVote.value) return

  busy.value = true

  try {
    await store.removeVote(props.articleId)
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <div class="reader-rating">
    <span class="reader-rating__label">Оценка читателей:</span>

    <span
      class="reader-rating__stars"
      role="group"
      aria-label="Оценить материал от 1 до 5 звёзд"
      @mouseleave="hovered = 0"
    >
      <button
        v-for="star in 5"
        :key="star"
        class="reader-rating__star"
        type="button"
        :class="{ 'reader-rating__star--filled': star <= filledUpTo }"
        :aria-label="`Оценить на ${star} из 5`"
        :title="auth.account ? `Оценить на ${star} из 5` : 'Войдите, чтобы оценить'"
        :disabled="busy"
        @click="send(star)"
        @mouseenter="hovered = star"
      >
        <svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true">
          <path
            d="M12 2.5 14.9 8.6l6.6.8-4.9 4.5 1.3 6.5L12 17.2l-5.9 3.2 1.3-6.5L2.5 9.4l6.6-.8L12 2.5Z"
            :fill="star <= filledUpTo ? 'currentColor' : 'none'"
            stroke="currentColor"
            stroke-width="1.6"
            stroke-linejoin="round"
          />
        </svg>
      </button>
    </span>

    <span class="reader-rating__value">
      {{ averageLabel }} · {{ votesLabel }}
    </span>

    <button
      v-if="auth.account && myVote"
      class="reader-rating__remove"
      type="button"
      :disabled="busy"
      @click="remove"
    >
      Убрать оценку
    </button>
  </div>
</template>

<style scoped>
.reader-rating {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
}

.reader-rating__label {
  color: var(--color-muted);
  font: 400 13px/16px Roboto, sans-serif;
}

.reader-rating__stars {
  display: inline-flex;
  gap: 2px;
}

.reader-rating__star {
  display: grid;
  place-items: center;
  width: 22px;
  height: 22px;
  padding: 0;
  color: var(--color-accent);
  background: none;
  border: 0;
  border-radius: 4px;
  cursor: pointer;
}

.reader-rating__star:hover {
  background: var(--color-card);
}

.reader-rating__value {
  color: var(--color-muted);
  font: 500 13px/16px Roboto, sans-serif;
}

.reader-rating__remove {
  padding: 2px 6px;
  color: var(--color-muted);
  background: none;
  border: 0;
  border-radius: 4px;
  font: 400 12px/15px Roboto, sans-serif;
  text-decoration: underline;
  cursor: pointer;
}

.reader-rating__remove:hover {
  color: var(--color-danger);
}
</style>
