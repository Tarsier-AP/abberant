<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useBookmarksStore } from '@/stores/bookmarks'

const props = defineProps<{ articleId: string }>()

const auth = useAuthStore()
const bookmarks = useBookmarksStore()
const router = useRouter()

const busy = ref(false)

const active = computed(() => bookmarks.isBookmarked(props.articleId))

const label = computed(() =>
  active.value ? 'Убрать из закладок' : 'Добавить в закладки',
)

async function toggle() {
  if (!auth.account) {
    void router.push({ name: 'login' })
    return
  }

  if (busy.value) return

  busy.value = true

  try {
    await bookmarks.toggle(props.articleId)
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <button
    class="bookmark"
    type="button"
    :class="{ 'bookmark--active': active }"
    :aria-label="label"
    :aria-pressed="active"
    :title="label"
    :disabled="busy"
    @click="toggle"
  >
    <svg width="12" height="12" viewBox="0 0 12 12" aria-hidden="true">
      <path
        d="M2.5 1h7a.5.5 0 0 1 .5.5v9.4a.3.3 0 0 1-.47.25L6 8.6l-3.53 2.55A.3.3 0 0 1 2 10.9V1.5a.5.5 0 0 1 .5-.5Z"
        :fill="active ? 'currentColor' : 'none'"
        stroke="currentColor"
        stroke-width="1.3"
        stroke-linejoin="round"
      />
    </svg>
  </button>
</template>

<style scoped>
.bookmark {
  display: grid;
  flex: none;
  place-items: center;
  width: 20px;
  height: 20px;
  padding: 0;
  color: inherit;
  background: none;
  border: 0;
  border-radius: 4px;
  cursor: pointer;
}

.bookmark:hover {
  background: var(--color-card);
}

.bookmark--active {
  color: var(--color-accent);
}

.bookmark:disabled {
  opacity: 0.6;
  cursor: default;
}
</style>
