<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'

const props = defineProps<{
  url?: string
  label?: string
}>()

const visible = ref(false)
const failed = ref(false)

let hideTimer: number | undefined

const message = computed(() =>
  failed.value ? 'Не удалось скопировать ссылку' : 'Ссылка скопирована',
)

function resolveUrl() {
  if (!props.url) return window.location.href

  return new URL(props.url, window.location.origin).href
}

async function copyToClipboard(text: string) {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(text)
    return
  }

  const textarea = document.createElement('textarea')

  textarea.value = text
  textarea.setAttribute('readonly', '')
  textarea.style.position = 'fixed'
  textarea.style.opacity = '0'
  document.body.appendChild(textarea)
  textarea.select()

  const copied = document.execCommand('copy')

  document.body.removeChild(textarea)

  if (!copied) throw new Error('Copy failed')
}

function showToast() {
  visible.value = true
  window.clearTimeout(hideTimer)
  hideTimer = window.setTimeout(() => {
    visible.value = false
  }, 2200)
}

async function share() {
  try {
    await copyToClipboard(resolveUrl())
    failed.value = false
  } catch {
    failed.value = true
  }

  showToast()
}

onBeforeUnmount(() => {
  window.clearTimeout(hideTimer)
})
</script>

<template>
  <button class="share" type="button" :aria-label="label ?? 'Поделиться'" @click="share">
    <svg width="12" height="12" viewBox="0 0 12 12" fill="none" aria-hidden="true">
      <path
        d="M6 7.5v-6M6 1.5 3.5 4M6 1.5 8.5 4"
        stroke="currentColor"
        stroke-width="1.5"
        stroke-linecap="round"
        stroke-linejoin="round"
      />
      <path
        d="M2 6.5V10a.5.5 0 0 0 .5.5h7A.5.5 0 0 0 10 10V6.5"
        stroke="currentColor"
        stroke-width="1.5"
        stroke-linecap="round"
      />
    </svg>
  </button>

  <Teleport to="body">
    <Transition name="share-toast">
      <p
        v-if="visible"
        class="share-toast"
        :class="{ 'share-toast--error': failed }"
        role="status"
        aria-live="polite"
      >
        {{ message }}
      </p>
    </Transition>
  </Teleport>
</template>

<style scoped>
.share {
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

.share:hover {
  background: var(--color-card);
}

.share-toast {
  position: fixed;
  z-index: 1000;
  bottom: 24px;
  left: 50%;
  padding: 8px 16px;
  color: var(--color-bg);
  background: var(--color-text);
  border-radius: var(--radius);
  box-shadow: var(--shadow-hover);
  font: 500 14px/20px Roboto, sans-serif;
  transform: translateX(-50%);
}

.share-toast--error {
  color: #ffffff;
  background: var(--color-danger);
}

.share-toast-enter-active,
.share-toast-leave-active {
  transition:
    opacity 0.2s ease,
    transform 0.2s ease;
}

.share-toast-enter-from,
.share-toast-leave-to {
  opacity: 0;
  transform: translate(-50%, 8px);
}
</style>
