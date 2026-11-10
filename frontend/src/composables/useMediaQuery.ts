import { onBeforeUnmount, onMounted, ref } from 'vue'

/**
 * Reactive `window.matchMedia`. Returns `false` until mounted, which keeps the
 * first render stable and avoids touching `window` during SSR.
 */
export function useMediaQuery(query: string) {
  const matches = ref(false)
  let list: MediaQueryList | null = null

  function update() {
    matches.value = list?.matches ?? false
  }

  onMounted(() => {
    list = window.matchMedia(query)
    update()
    list.addEventListener('change', update)
  })

  onBeforeUnmount(() => {
    list?.removeEventListener('change', update)
    list = null
  })

  return matches
}
