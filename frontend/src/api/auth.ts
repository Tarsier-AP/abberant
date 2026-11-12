import { api } from './client'
import type { AuthAccount, TwoFactorMethod } from '@/types/auth'

export interface RegistrationResponse {
  status: string
  registrationId: string
}

export interface DevCodeResponse {
  status: string
  devCode?: string
}

export interface TwoFactorChallenge {
  status: '2fa_required'
  twoFactorMethod: TwoFactorMethod | ''
  emailCodeSent?: boolean
  devCode?: string
}

export interface LoginResponse {
  status: 'ok'
  account: AuthAccount
}

export type LoginResult = TwoFactorChallenge | LoginResponse

export interface AccountResponse {
  account: AuthAccount
}

export interface TotpSetupResponse {
  status: string
  secret: string
  otpauthUrl: string
  qrDataURL: string
}

export interface BackupCodesResponse {
  status: string
  backupCodes: string[]
}

export interface ProfilePatch {
  username?: string
  avatar?: string
  additionalEmail?: string
  phone?: string
  smsNotifications?: boolean
  emailNotifications?: boolean
  twoFactorEnabled?: boolean
  twoFactorMethod?: TwoFactorMethod | ''
}

export function startRegistration(email: string, username: string) {
  return api<RegistrationResponse>('/register', {
    method: 'POST',
    body: { email, username },
  })
}

export function setRegistrationPassword(registrationId: string, password: string) {
  return api<DevCodeResponse>('/register/password', {
    method: 'POST',
    body: { registrationId, password },
  })
}

export function activateRegistration(registrationId: string, code: string) {
  return api<{ status: string }>('/register/activate', {
    method: 'POST',
    body: { registrationId, code },
  })
}

export function activateByEmail(email: string, code: string) {
  return api<{ status: string }>('/activate', {
    method: 'POST',
    body: { email, code },
  })
}

export function login(email: string, password: string) {
  return api<LoginResult>('/login', {
    method: 'POST',
    body: { email, password },
  })
}

export function completeTwoFactorLogin(code: string) {
  return api<AccountResponse>('/login/2fa', {
    method: 'POST',
    body: { code },
  })
}

export function logout() {
  return api<{ status: string }>('/logout', { method: 'POST' })
}

export function logoutAll() {
  return api<{ status: string }>('/logout/all', { method: 'POST' })
}

export function fetchMe() {
  return api<AccountResponse>('/me')
}

export function patchMe(patch: ProfilePatch) {
  return api<AccountResponse>('/me', { method: 'PATCH', body: patch })
}

export function setupTotp() {
  return api<TotpSetupResponse>('/2fa/totp/setup', { method: 'POST' })
}

export function confirmTotp(code: string) {
  return api<BackupCodesResponse>('/2fa/totp/confirm', { method: 'POST', body: { code } })
}

export function sendEmailTwoFactorCode() {
  return api<{ status: string; emailCodeSent: boolean; devCode?: string }>('/2fa/email/send', {
    method: 'POST',
  })
}

export function regenerateBackupCodes() {
  return api<BackupCodesResponse>('/2fa/backup/codes/regenerate', { method: 'POST' })
}
