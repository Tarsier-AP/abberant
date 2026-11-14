<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useNotificationsStore } from '@/stores/notifications'
import NotificationsPanel from '@/components/layout/NotificationsPanel.vue'

const auth = useAuthStore()
const store = useNotificationsStore()

const open = ref(false)
const root = ref<HTMLElement | null>(null)

const badge = computed(() =>
  store.unread > 99 ? '99+' : store.unread > 0 ? String(store.unread) : '',
)

const label = computed(() =>
  store.unread > 0 ? `Уведомления: ${store.unread} непрочитанных` : 'Уведомления',
)

function toggle() {
  if (!auth.account) return

  open.value = !open.value

  if (open.value) void store.open()
}

function onDocumentClick(event: MouseEvent) {
  if (open.value && root.value && !root.value.contains(event.target as Node)) {
    open.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', onDocumentClick)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onDocumentClick)
})
</script>

<template>
  <div ref="root" class="bell">
    <button
      class="bell__button"
      type="button"
      :aria-label="label"
      :title="label"
      :aria-expanded="open"
      @click="toggle"
    >
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" aria-hidden="true">
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

      <span v-if="badge" class="bell__badge">{{ badge }}</span>
    </button>

    <div v-if="open" class="bell__panel">
      <NotificationsPanel heading="Уведомления" />
    </div>
  </div>
</template>

<style scoped>
.bell {
  position: relative;
}

.bell__button {
  display: grid;
  flex: none;
  place-items: center;
  position: relative;
  width: 40px;
  height: 40px;
  padding: 0;
  color: inherit;
  background: var(--color-card);
  border: 0;
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  cursor: pointer;
}

.bell__button:hover {
  background: var(--color-card-hover);
}

.bell__badge {
  position: absolute;
  top: 2px;
  right: 0;
  min-width: 16px;
  padding: 1px 4px;
  color: var(--color-accent-text);
  background: var(--color-danger);
  border-radius: 999px;
  font: 500 10px/12px Roboto, sans-serif;
  text-align: center;
}

.bell__panel {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  z-index: 60;
  width: min(340px, calc(100vw - 32px));
  max-height: 420px;
  overflow-y: auto;
  padding: 12px;
  background: var(--color-header);
  border-radius: var(--radius);
  box-shadow: var(--shadow-hover);
}
</style>
