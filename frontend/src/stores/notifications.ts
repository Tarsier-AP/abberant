import { ref, watch } from 'vue'
import { defineStore } from 'pinia'
import type { NotificationPayload } from '@/api/notifications'
import * as api from '@/api/notifications'
import { useAuthStore } from '@/stores/auth'

/** How often the bell checks for new notifications. */
const POLL_INTERVAL = 30_000

export const useNotificationsStore = defineStore('notifications', () => {
  const auth = useAuthStore()

  const items = ref<NotificationPayload[]>([])
  const unread = ref(0)
  /** Sections the current account follows (news platform ids). */
  const sections = ref<string[]>([])

  let timer: number | undefined

  function stopPolling() {
    if (timer !== undefined) {
      window.clearInterval(timer)
      timer = undefined
    }
  }

  function startPolling() {
    if (timer !== undefined) return

    timer = window.setInterval(() => {
      if (!auth.account) return

      void api
        .fetchUnreadCount()
        .then(({ count }) => {
          unread.value = count
        })
        .catch(() => {
          // Transient failures just wait for the next tick.
        })
    }, POLL_INTERVAL)
  }

  /** Refreshes the badge and starts polling after login; clears on logout. */
  async function sync() {
    stopPolling()

    if (!auth.account) {
      items.value = []
      unread.value = 0
      sections.value = []
      return
    }

    startPolling()

    try {
      const [count, subscriptionList] = await Promise.all([
        api.fetchUnreadCount(),
        api.fetchSubscriptions(),
      ])

      unread.value = count.count
      sections.value = subscriptionList.sections
    } catch {
      // The bell simply stays empty until the next successful poll.
    }
  }

  /** Loads the panel list and marks everything read. */
  async function open() {
    if (!auth.account) return

    try {
      const { notifications: loaded } = await api.fetchNotifications()

      items.value = loaded
    } catch {
      return
    }

    if (unread.value > 0) {
      try {
        await api.markAllNotificationsRead()
        unread.value = 0
        items.value = items.value.map((item) => ({ ...item, read: true }))
      } catch {
        // The list stays; the badge catches up on the next poll.
      }
    }
  }

  async function subscribeTo(section: string) {
    await api.subscribe(section)

    if (!sections.value.includes(section)) {
      sections.value = [...sections.value, section]
    }
  }

  async function unsubscribeFrom(section: string) {
    await api.unsubscribe(section)

    sections.value = sections.value.filter((id) => id !== section)
  }

  function isSubscribed(section: string): boolean {
    return sections.value.includes(section)
  }

  watch(
    () => auth.account,
    () => {
      void sync()
    },
  )

  return {
    items,
    unread,
    sections,
    sync,
    open,
    subscribeTo,
    unsubscribeFrom,
    isSubscribed,
  }
})
