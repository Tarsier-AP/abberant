import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import type { TwoFactorMethod } from '@/types/auth'

export interface AuthAccount {
  login: string
  username: string
  additionalEmail: string
  activated: boolean
  twoFactorEnabled: boolean
  twoFactorMethod: TwoFactorMethod
  smsNotifications: boolean
}

export interface PendingRegistration {
  login: string
  username: string
}

const ACCOUNT_KEY = 'abberant-auth-account'
const PENDING_KEY = 'abberant-auth-pending-registration'
const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

function readStored<T>(key: string): T | null {
  try {
    const raw = localStorage.getItem(key)
    return raw ? (JSON.parse(raw) as T) : null
  } catch {
    return null
  }
}

function writeStored(key: string, value: unknown) {
  try {
    if (value === null) {
      localStorage.removeItem(key)
    } else {
      localStorage.setItem(key, JSON.stringify(value))
    }
  } catch {
    // Ignore storage errors (e.g. private browsing).
  }
}

function normalizeLogin(value: string): string {
  const login = value.trim()

  return EMAIL_PATTERN.test(login) ? login.toLowerCase() : login.replace(/\D/g, '')
}

export const useAuthStore = defineStore('auth', () => {
  const storedAccount = readStored<AuthAccount>(ACCOUNT_KEY)

  const account = ref<AuthAccount | null>(
    storedAccount
      ? {
          ...storedAccount,
          additionalEmail: storedAccount.additionalEmail ?? '',
          twoFactorMethod: storedAccount.twoFactorMethod ?? 'sms',
          smsNotifications: storedAccount.smsNotifications ?? true,
        }
      : null,
  )
  const pendingRegistration = ref<PendingRegistration | null>(
    readStored<PendingRegistration>(PENDING_KEY),
  )

  const canUseTwoFactor = computed(
    () => account.value?.activated === true && account.value.twoFactorEnabled === true,
  )

  function startRegistration(data: PendingRegistration) {
    pendingRegistration.value = data
    writeStored(PENDING_KEY, data)
  }

  function cancelRegistration() {
    pendingRegistration.value = null
    writeStored(PENDING_KEY, null)
  }

  function completeRegistration() {
    if (!pendingRegistration.value) return false

    account.value = {
      login: pendingRegistration.value.login,
      username: pendingRegistration.value.username,
      additionalEmail: '',
      activated: false,
      twoFactorEnabled: false,
      twoFactorMethod: 'sms',
      smsNotifications: true,
    }

    cancelRegistration()
    writeStored(ACCOUNT_KEY, account.value)

    return true
  }

  function createAccountFromActivation(data: PendingRegistration) {
    if (account.value) return false

    account.value = {
      login: data.login,
      username: data.username,
      additionalEmail: '',
      activated: true,
      twoFactorEnabled: false,
      twoFactorMethod: 'sms',
      smsNotifications: true,
    }

    writeStored(ACCOUNT_KEY, account.value)

    return true
  }

  function activateAccount(login?: string) {
    if (!account.value) return false
    if (login !== undefined && !matchesAccount(login)) return false

    account.value.activated = true
    writeStored(ACCOUNT_KEY, account.value)

    return true
  }

  function setTwoFactor(enabled: boolean) {
    if (!account.value) return

    account.value.twoFactorEnabled = enabled
    writeStored(ACCOUNT_KEY, account.value)
  }

  function setTwoFactorMethod(method: TwoFactorMethod) {
    if (!account.value) return

    account.value.twoFactorMethod = method
    writeStored(ACCOUNT_KEY, account.value)
  }

  function setAdditionalEmail(email: string) {
    if (!account.value) return

    account.value.additionalEmail = email
    writeStored(ACCOUNT_KEY, account.value)
  }

  function setSmsNotifications(enabled: boolean) {
    if (!account.value) return

    account.value.smsNotifications = enabled
    writeStored(ACCOUNT_KEY, account.value)
  }

  function matchesAccount(login: string) {
    if (!account.value) return false

    return normalizeLogin(account.value.login) === normalizeLogin(login)
  }

  function reset() {
    account.value = null
    pendingRegistration.value = null
    writeStored(ACCOUNT_KEY, null)
    writeStored(PENDING_KEY, null)
  }

  return {
    account,
    pendingRegistration,
    canUseTwoFactor,
    startRegistration,
    cancelRegistration,
    completeRegistration,
    createAccountFromActivation,
    activateAccount,
    setTwoFactor,
    setTwoFactorMethod,
    setAdditionalEmail,
    setSmsNotifications,
    matchesAccount,
    reset,
  }
})
