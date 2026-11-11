import { onBeforeUnmount, onMounted } from 'vue'
import { useRouter } from 'vue-router'

/** Temporary shortcut that opens the 2FA screen while the backend is absent. */
export const TWO_FACTOR_DEMO_HOTKEY = 'Ctrl + Shift + 2'

export function useTwoFactorDemoShortcut() {
  const router = useRouter()

  function handleKeydown(event: KeyboardEvent) {
    if (!event.ctrlKey || !event.shiftKey || event.altKey || event.metaKey) return
    if (event.code !== 'Digit2') return

    event.preventDefault()
    void router.push({ name: 'two-factor', query: { demo: '1' } })
  }

  onMounted(() => window.addEventListener('keydown', handleKeydown))
  onBeforeUnmount(() => window.removeEventListener('keydown', handleKeydown))
}
