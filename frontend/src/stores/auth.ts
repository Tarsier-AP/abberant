import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import type { TwoFactorMethod, UserRole } from '@/types/auth'
import type { ProfilePatch } from '@/api/auth'
import * as api from '@/api/auth'

export interface AuthAccount {
  login: string
  username: string
  avatar: string
  additionalEmail: string
  phone: string
  activated: boolean
  twoFactorEnabled: boolean
  twoFactorMethod: TwoFactorMethod | ''
  smsNotifications: boolean
  emailNotifications: boolean
  role: UserRole
}

export interface PendingRegistration {
  login: string
  username: string
  registrationId: string
  /** Single-use code echoed back by the server in dev mode. */
  devCode?: string
}

const PENDING_KEY = 'abberant-auth-pending-registration'

function readStoredPending(): PendingRegistration | null {
  try {
    const raw = localStorage.getItem(PENDING_KEY)

    return raw ? (JSON.parse(raw) as PendingRegistration) : null
  } catch {
    return null
  }
}

function writeStoredPending(value: PendingRegistration | null) {
  try {
    if (value === null) {
      localStorage.removeItem(PENDING_KEY)
    } else {
      localStorage.setItem(PENDING_KEY, JSON.stringify(value))
    }
  } catch {
    // Ignore storage errors (e.g. private browsing).
  }
}

export const useAuthStore = defineStore('auth', () => {
  const account = ref<AuthAccount | null>(null)
  const pendingRegistration = ref<PendingRegistration | null>(readStoredPending())
  /** Set between the password step of the login and the 2FA code step. */
  const pendingLogin = ref<{ method: TwoFactorMethod | '' } | null>(null)
  const initialized = ref(false)

  const canUseTwoFactor = computed(
    () => account.value?.activated === true && account.value.twoFactorEnabled === true,
  )

  const canEditArticles = computed(
    () =>
      account.value?.role === 'editor' ||
      account.value?.role === 'owner' ||
      account.value?.role === 'creator',
  )

  const canModerate = computed(
    () =>
      account.value?.role === 'editor' ||
      account.value?.role === 'owner' ||
      account.value?.role === 'moderator',
  )

  function setAccount(next: AuthAccount | null) {
    account.value = next
  }

  /** Loads the current session account (or clears it) once on app start. */
  async function init() {
    if (initialized.value) return

    try {
      const { account: loaded } = await api.fetchMe()

      account.value = loaded
    } catch {
      account.value = null
    } finally {
      initialized.value = true
    }
  }

  function startRegistration(data: Omit<PendingRegistration, 'registrationId'>) {
    pendingRegistration.value = { ...data, registrationId: '' }
  }

  function completeStartRegistration(registrationId: string) {
    if (!pendingRegistration.value) return

    pendingRegistration.value.registrationId = registrationId
    writeStoredPending(pendingRegistration.value)
  }

  function cancelRegistration() {
    pendingRegistration.value = null
    writeStoredPending(null)
  }

  function setPendingLogin(method: TwoFactorMethod | '') {
    pendingLogin.value = { method }
  }

  function clearPendingLogin() {
    pendingLogin.value = null
  }

  async function applyPatch(patch: ProfilePatch) {
    const { account: updated } = await api.patchMe(patch)

    account.value = updated
  }

  async function endSession() {
    try {
      await api.logout()
    } finally {
      account.value = null
      pendingLogin.value = null
    }
  }

  function reset() {
    account.value = null
    pendingRegistration.value = null
    pendingLogin.value = null
    writeStoredPending(null)
  }

  return {
    account,
    pendingRegistration,
    pendingLogin,
    canUseTwoFactor,
    canEditArticles,
    canModerate,
    init,
    setAccount,
    startRegistration,
    completeStartRegistration,
    cancelRegistration,
    setPendingLogin,
    clearPendingLogin,
    applyPatch,
    endSession,
    reset,
  }
})
