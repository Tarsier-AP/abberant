<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useNotificationsStore } from '@/stores/notifications'

const props = defineProps<{
  /** News platform section id, e.g. `pc`. */
  section: string
  /** Human-readable section title for the accessible labels. */
  title: string
}>()

const auth = useAuthStore()
const store = useNotificationsStore()
const router = useRouter()

const busy = ref(false)

const subscribed = computed(() => store.isSubscribed(props.section))

const label = computed(() =>
  subscribed.value ? `Отписаться от раздела «${props.title}»` : `Подписаться на раздел «${props.title}»`,
)

async function toggle() {
  if (!auth.account) {
    void router.push({ name: 'login' })
    return
  }

  if (busy.value) return

  busy.value = true

  try {
    if (subscribed.value) {
      await store.unsubscribeFrom(props.section)
    } else {
      await store.subscribeTo(props.section)
    }
  } finally {
    busy.value = false
  }
}
</script>

<template>
  <button
    class="subscribe"
    type="button"
    :class="{ 'subscribe--active': subscribed }"
    :aria-pressed="subscribed"
    :aria-label="label"
    :title="label"
    :disabled="busy"
    @click="toggle"
  >
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" aria-hidden="true">
      <path
        d="M6 16v-5a6 6 0 1 1 12 0v5l1.5 2.5h-15L6 16Z"
        stroke="currentColor"
        stroke-width="2"
        stroke-linejoin="round"
      />
      <path
        d="M10 21a2.4 2.4 0 0 0 4 0"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
      />
    </svg>

    <span>{{ subscribed ? 'Вы подписаны' : 'Подписаться' }}</span>
  </button>
</template>

<style scoped>
.subscribe {
  display: inline-flex;
  flex: none;
  align-items: center;
  gap: 6px;
  min-height: 32px;
  padding: 4px 12px;
  color: var(--color-text);
  background: var(--color-card);
  border: 0;
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  font: 500 12px/16px Roboto, sans-serif;
  cursor: pointer;
}

.subscribe:hover {
  background: var(--color-card-hover);
}

.subscribe--active {
  color: var(--color-accent-text);
  background: var(--color-accent);
}

.subscribe:disabled {
  opacity: 0.6;
  cursor: default;
}
</style>
