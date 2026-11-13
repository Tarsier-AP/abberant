<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import * as api from '@/api/auth'
import type { TwoFactorMethod } from '@/types/auth'

const router = useRouter()
const auth = useAuthStore()

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const AVATAR_MAX_BYTES = 5 * 1024 * 1024
const AVATAR_SIZE = 160

const busy = ref(false)
const settingsError = ref('')

const avatarInput = ref<HTMLInputElement | null>(null)
const avatarError = ref('')

const emailInput = ref<HTMLInputElement | null>(null)
const editingEmail = ref(false)
const emailDraft = ref('')
const emailError = ref('')

const phoneInput = ref<HTMLInputElement | null>(null)
const editingPhone = ref(false)
const phoneDraft = ref('')
const phoneError = ref('')

const totpCodeInput = ref<HTMLInputElement | null>(null)
const enrollingTotp = ref(false)
const totpSetup = ref<{ secret: string; qrDataURL: string } | null>(null)
const totpCode = ref('')
const totpError = ref('')
const backupCodes = ref<string[] | null>(null)
const backupCopied = ref(false)

const isActivated = computed(() => auth.account?.activated === true)

const account = computed(() => auth.account)

const phone = computed(() => account.value?.phone ?? '')
const additionalEmail = computed(() => account.value?.additionalEmail ?? '')

async function withBusy(action: () => Promise<void>) {
  if (busy.value) return

  busy.value = true
  settingsError.value = ''

  try {
    await action()
  } catch (requestError) {
    settingsError.value =
      requestError instanceof Error ? requestError.message : 'Не удалось сохранить изменения.'
  } finally {
    busy.value = false
  }
}

async function createAvatarDataUrl(file: File): Promise<string> {
  const bitmap = await createImageBitmap(file)
  const canvas = document.createElement('canvas')
  const context = canvas.getContext('2d')

  canvas.width = AVATAR_SIZE
  canvas.height = AVATAR_SIZE

  if (!context) {
    bitmap.close()
    throw new Error('Canvas is not supported')
  }

  const scale = Math.max(AVATAR_SIZE / bitmap.width, AVATAR_SIZE / bitmap.height)
  const width = bitmap.width * scale
  const height = bitmap.height * scale

  context.drawImage(bitmap, (AVATAR_SIZE - width) / 2, (AVATAR_SIZE - height) / 2, width, height)
  bitmap.close()

  return canvas.toDataURL('image/webp', 0.85)
}

async function handleAvatarChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]

  input.value = ''
  avatarError.value = ''

  if (!file) return

  if (!file.type.startsWith('image/')) {
    avatarError.value = 'Выберите файл изображения (JPG, PNG или WebP).'
    return
  }

  if (file.size > AVATAR_MAX_BYTES) {
    avatarError.value = 'Файл слишком большой. Максимальный размер — 5 МБ.'
    return
  }

  try {
    const dataUrl = await createAvatarDataUrl(file)

    await withBusy(() => auth.applyPatch({ avatar: dataUrl }))
  } catch {
    avatarError.value = 'Не удалось обработать изображение. Попробуйте другой файл.'
  }
}

function removeAvatar() {
  void withBusy(() => auth.applyPatch({ avatar: '' }))
}

function startEmailEdit() {
  emailDraft.value = additionalEmail.value
  emailError.value = ''
  editingEmail.value = true

  void nextTick(() => emailInput.value?.focus())
}

function cancelEmailEdit() {
  editingEmail.value = false
  emailError.value = ''
}

function saveEmail() {
  const email = emailDraft.value.trim()

  if (email && !EMAIL_PATTERN.test(email)) {
    emailError.value = 'Введите корректный адрес электронной почты.'
    return
  }

  void withBusy(async () => {
    await auth.applyPatch({ additionalEmail: email })
    editingEmail.value = false
  })
}

function startPhoneEdit() {
  phoneDraft.value = phone.value
  phoneError.value = ''
  editingPhone.value = true

  void nextTick(() => phoneInput.value?.focus())
}

function cancelPhoneEdit() {
  editingPhone.value = false
  phoneError.value = ''
}

function savePhone() {
  const value = phoneDraft.value.trim()

  if (value.replace(/\D/g, '').length > 15) {
    phoneError.value = 'Введите корректный номер телефона.'
    return
  }

  void withBusy(async () => {
    await auth.applyPatch({ phone: value })
    editingPhone.value = false
  })
}

function toggleTwoFactor() {
  if (!account.value || !isActivated.value) return

  void withBusy(async () => {
    if (account.value?.twoFactorEnabled) {
      backupCodes.value = null
      enrollingTotp.value = false
      totpSetup.value = null

      await auth.applyPatch({ twoFactorEnabled: false })
      return
    }

    if (account.value?.twoFactorMethod === 'authenticator') {
      await startTotpEnrollment()
      return
    }

    await auth.applyPatch({
      twoFactorEnabled: true,
      twoFactorMethod: account.value?.twoFactorMethod || 'email',
    })
  })
}

function handleMethodChange(event: Event) {
  const method = (event.target as HTMLSelectElement).value as TwoFactorMethod

  if (method === 'authenticator') {
    void withBusy(startTotpEnrollment)
    return
  }

  void withBusy(() => auth.applyPatch({ twoFactorMethod: method }))
}

async function startTotpEnrollment() {
  const setup = await api.setupTotp()

  totpSetup.value = { secret: setup.secret, qrDataURL: setup.qrDataURL }
  enrollingTotp.value = true
  totpCode.value = ''
  totpError.value = ''

  void nextTick(() => totpCodeInput.value?.focus())
}

function cancelTotpEnrollment() {
  enrollingTotp.value = false
  totpSetup.value = null
  totpError.value = ''
}

function confirmTotpEnrollment() {
  const code = totpCode.value.replace(/\D/g, '')

  if (!/^\d{6}$/.test(code)) {
    totpError.value = 'Введите 6-значный код из приложения-аутентификатора.'
    return
  }

  if (busy.value) return

  busy.value = true
  settingsError.value = ''
  totpError.value = ''

  void (async () => {
    try {
      const result = await api.confirmTotp(code)

      backupCodes.value = result.backupCodes
      backupCopied.value = false
      cancelTotpEnrollment()

      const { account: updated } = await api.fetchMe()
      auth.setAccount(updated)
    } catch (requestError) {
      totpError.value =
        requestError instanceof Error ? requestError.message : 'Не удалось подтвердить код.'
    } finally {
      busy.value = false
    }
  })()
}

async function copyBackupCodes() {
  if (!backupCodes.value) return

  try {
    await navigator.clipboard.writeText(backupCodes.value.join('\n'))
    backupCopied.value = true
  } catch {
    backupCopied.value = false
  }
}

function regenerateBackupCodes() {
  void withBusy(async () => {
    const result = await api.regenerateBackupCodes()

    backupCodes.value = result.backupCodes
    backupCopied.value = false
  })
}

function toggleSmsNotifications() {
  if (!account.value) return

  void withBusy(() => auth.applyPatch({ smsNotifications: !account.value?.smsNotifications }))
}

function toggleEmailNotifications() {
  if (!account.value) return

  void withBusy(() => auth.applyPatch({ emailNotifications: !account.value?.emailNotifications }))
}

function handleLogout() {
  void withBusy(async () => {
    await auth.endSession()
    await router.push({ name: 'home' })
  })
}

function handleLogoutAll() {
  void withBusy(async () => {
    await auth.endAllSessions()
    await router.push({ name: 'home' })
  })
}

onMounted(() => {
  if (!auth.account) {
    void router.replace({ name: 'login' })
  }
})
</script>

<template>
  <div class="account-settings container">
    <h1 class="visually-hidden">Настройки аккаунта</h1>

    <div v-if="auth.account" class="account-settings__card">
      <section class="account-settings__section" aria-labelledby="account-personal">
        <h2 id="account-personal" class="account-settings__section-title">Личная информация</h2>

        <div class="account-settings__profile">
          <button
            class="account-settings__avatar"
            type="button"
            aria-label="Изменить фото профиля"
            @click="avatarInput?.click()"
          >
            <img
              v-if="auth.account.avatar"
              class="account-settings__avatar-image"
              :src="auth.account.avatar"
              alt=""
            />
            <svg
              v-else
              class="account-settings__avatar-placeholder"
              width="48"
              height="48"
              viewBox="0 0 48 48"
              fill="currentColor"
              aria-hidden="true"
            >
              <circle cx="24" cy="17" r="8" opacity="0.7" />
              <path d="M8 46c0-8.837 7.163-16 16-16s16 7.163 16 16" opacity="0.7" />
            </svg>

            <span class="account-settings__avatar-overlay" aria-hidden="true">
              Изменить
            </span>
          </button>

          <div class="account-settings__profile-text">
            <span class="account-settings__profile-title">Фото профиля</span>
            <span class="account-settings__profile-hint">
              Наведите на аватар или нажмите, чтобы изменить
            </span>
          </div>

          <button
            v-if="auth.account.avatar"
            class="account-settings__profile-button account-settings__profile-button--muted"
            type="button"
            @click="removeAvatar"
          >
            Удалить фото
          </button>

          <input
            ref="avatarInput"
            hidden
            type="file"
            accept="image/*"
            @change="handleAvatarChange"
          />

          <p
            v-if="avatarError"
            class="account-settings__error account-settings__error--profile"
            role="alert"
          >
            {{ avatarError }}
          </p>
        </div>

        <div class="account-settings__rows">
          <button
            class="account-settings__row account-settings__row--interactive"
            type="button"
            @click="startPhoneEdit"
          >
            <span class="account-settings__label">Телефон:</span>
            <span class="account-settings__value">{{ phone || 'Добавить' }}</span>
          </button>

          <button
            class="account-settings__row account-settings__row--interactive"
            type="button"
            @click="startEmailEdit"
          >
            <span class="account-settings__label">Дополнительная почта:</span>
            <span class="account-settings__value">
              {{ additionalEmail || 'Добавить' }}
            </span>
          </button>

          <div class="account-settings__row">
            <span class="account-settings__label">Электронная почта:</span>
            <span class="account-settings__value">{{ auth.account.login }}</span>
          </div>
        </div>

        <form
          v-if="editingEmail"
          class="account-settings__email-form"
          novalidate
          @submit.prevent="saveEmail"
        >
          <label class="visually-hidden" for="account-additional-email">
            Дополнительная почта
          </label>
          <input
            id="account-additional-email"
            ref="emailInput"
            v-model="emailDraft"
            class="account-settings__email-input"
            type="email"
            autocomplete="email"
            placeholder="name@example.com"
            :aria-invalid="emailError ? 'true' : undefined"
            :aria-describedby="emailError ? 'account-email-error' : undefined"
          />
          <p
            v-if="emailError"
            id="account-email-error"
            class="account-settings__error"
            role="alert"
          >
            {{ emailError }}
          </p>
          <div class="account-settings__email-actions">
            <button
              class="account-settings__email-button"
              type="button"
              @click="cancelEmailEdit"
            >
              Отмена
            </button>
            <button class="account-settings__email-button" type="submit">Сохранить</button>
          </div>
        </form>

        <form
          v-if="editingPhone"
          class="account-settings__email-form"
          novalidate
          @submit.prevent="savePhone"
        >
          <label class="visually-hidden" for="account-phone">Телефон</label>
          <input
            id="account-phone"
            ref="phoneInput"
            v-model="phoneDraft"
            class="account-settings__email-input"
            type="tel"
            autocomplete="tel"
            placeholder="+7 900 000-00-00"
            :aria-invalid="phoneError ? 'true' : undefined"
            :aria-describedby="phoneError ? 'account-phone-error' : undefined"
          />
          <p
            v-if="phoneError"
            id="account-phone-error"
            class="account-settings__error"
            role="alert"
          >
            {{ phoneError }}
          </p>
          <div class="account-settings__email-actions">
            <button
              class="account-settings__email-button"
              type="button"
              @click="cancelPhoneEdit"
            >
              Отмена
            </button>
            <button class="account-settings__email-button" type="submit">Сохранить</button>
          </div>
        </form>
      </section>

      <div class="account-settings__divider" aria-hidden="true"></div>

      <section class="account-settings__section" aria-labelledby="account-security">
        <h2 id="account-security" class="account-settings__section-title">Безопасность</h2>

        <div class="account-settings__rows">
          <button
            class="account-settings__row account-settings__row--interactive"
            type="button"
            role="switch"
            :aria-checked="auth.account.twoFactorEnabled"
            :disabled="!isActivated || busy"
            @click="toggleTwoFactor"
          >
            <span class="account-settings__label">
              Двухфакторная аутентификация (2FA):
            </span>
            <span class="account-settings__value">
              {{
                !isActivated
                  ? 'Недоступна'
                  : auth.account.twoFactorEnabled
                    ? 'Включена'
                    : 'Выключена'
              }}
            </span>
          </button>

          <label
            v-if="auth.account.twoFactorEnabled"
            class="account-settings__row account-settings__row--select"
          >
            <span class="account-settings__label">Способ подтверждения:</span>
            <select
              :value="auth.account.twoFactorMethod"
              :disabled="busy"
              @change="handleMethodChange"
            >
              <option value="email">Электронная почта</option>
              <option value="authenticator">Приложение-аутентификатор</option>
              <option value="sms" disabled>SMS (скоро)</option>
            </select>
          </label>

          <RouterLink
            class="account-settings__row account-settings__row--link"
            to="/password-reset"
          >
            <span class="account-settings__label">Сменить пароль:</span>
            <span class="account-settings__value">Сменить</span>
          </RouterLink>

          <button
            class="account-settings__row account-settings__row--interactive"
            type="button"
            :disabled="busy"
            @click="handleLogout"
          >
            <span class="account-settings__label">Завершить сеанс:</span>
            <span class="account-settings__value">Выйти</span>
          </button>

          <button
            class="account-settings__row account-settings__row--interactive"
            type="button"
            :disabled="busy"
            @click="handleLogoutAll"
          >
            <span class="account-settings__label">Выйти со всех устройств:</span>
            <span class="account-settings__value">Завершить все сеансы</span>
          </button>
        </div>

        <p v-if="!isActivated" class="account-settings__notice">
          Активируйте аккаунт, чтобы включить двухфакторную аутентификацию.
        </p>

        <div v-if="enrollingTotp && totpSetup" class="account-settings__totp">
          <p class="account-settings__totp-hint">
            Отсканируйте QR-код в приложении-аутентификаторе и введите текущий код.
          </p>

          <img
            class="account-settings__totp-qr"
            :src="totpSetup.qrDataURL"
            alt="QR-код с секретным ключом приложения-аутентификатора"
          />

          <p class="account-settings__totp-secret">
            Если QR-код не сканируется, введите ключ вручную:
            <strong>{{ totpSetup.secret }}</strong>
          </p>

          <input
            ref="totpCodeInput"
            v-model="totpCode"
            class="account-settings__email-input account-settings__totp-input"
            type="text"
            inputmode="numeric"
            maxlength="6"
            autocomplete="one-time-code"
            placeholder="000000"
            :aria-invalid="totpError ? 'true' : undefined"
            :disabled="busy"
          />

          <p v-if="totpError" class="account-settings__error" role="alert">{{ totpError }}</p>

          <div class="account-settings__email-actions">
            <button
              class="account-settings__email-button"
              type="button"
              :disabled="busy"
              @click="cancelTotpEnrollment"
            >
              Отмена
            </button>
            <button
              class="account-settings__email-button"
              type="button"
              :disabled="busy"
              @click="confirmTotpEnrollment"
            >
              Подтвердить
            </button>
          </div>
        </div>

        <div v-if="backupCodes" class="account-settings__backup">
          <p class="account-settings__backup-hint">
            Одноразовые резервные коды. Сохраните их в надёжном месте: каждый код можно
            использовать один раз вместо кода подтверждения.
          </p>

          <ul class="account-settings__backup-list">
            <li v-for="code in backupCodes" :key="code" class="account-settings__backup-code">
              {{ code }}
            </li>
          </ul>

          <div class="account-settings__email-actions">
            <button
              class="account-settings__email-button"
              type="button"
              @click="copyBackupCodes"
            >
              {{ backupCopied ? 'Скопировано' : 'Скопировать' }}
            </button>
            <button
              class="account-settings__email-button"
              type="button"
              :disabled="busy"
              @click="regenerateBackupCodes"
            >
              Сгенерировать новые
            </button>
          </div>
        </div>

        <p v-if="settingsError" class="account-settings__error" role="alert">
          {{ settingsError }}
        </p>
      </section>

      <section class="account-settings__section" aria-labelledby="account-notifications">
        <h2 id="account-notifications" class="account-settings__section-title">Уведомления</h2>

        <div class="account-settings__rows">
          <button
            class="account-settings__row account-settings__row--interactive"
            type="button"
            role="switch"
            :aria-checked="auth.account.emailNotifications"
            :disabled="busy"
            @click="toggleEmailNotifications"
          >
            <span class="account-settings__label">
              <svg
                class="account-settings__icon"
                width="18"
                height="14"
                viewBox="0 0 18 14"
                fill="none"
                aria-hidden="true"
              >
                <rect
                  x="0.75"
                  y="0.75"
                  width="16.5"
                  height="12.5"
                  rx="2.25"
                  stroke="currentColor"
                  stroke-width="1.5"
                  opacity="0.6"
                />
                <path
                  d="M1.5 2.5 9 7.5l7.5-5"
                  stroke="currentColor"
                  stroke-width="1.5"
                  stroke-linejoin="round"
                  opacity="0.6"
                />
              </svg>
              Электронная почта:
            </span>
            <span class="account-settings__value">
              {{ auth.account.emailNotifications ? 'Включено' : 'Выключено' }}
            </span>
          </button>

          <button
            class="account-settings__row account-settings__row--interactive"
            type="button"
            role="switch"
            :aria-checked="auth.account.smsNotifications"
            :disabled="busy"
            @click="toggleSmsNotifications"
          >
            <span class="account-settings__label">
              <svg
                class="account-settings__icon"
                width="18"
                height="14"
                viewBox="0 0 18 14"
                fill="none"
                aria-hidden="true"
              >
                <rect
                  x="0.75"
                  y="0.75"
                  width="16.5"
                  height="12.5"
                  rx="3.25"
                  stroke="currentColor"
                  stroke-width="1.5"
                  opacity="0.6"
                />
                <circle cx="6" cy="7" r="1" fill="currentColor" />
                <circle cx="9" cy="7" r="1" fill="currentColor" />
                <circle cx="12" cy="7" r="1" fill="currentColor" />
              </svg>
              СМС:
            </span>
            <span class="account-settings__value">
              {{ auth.account.smsNotifications ? 'Включено' : 'Выключено' }}
            </span>
          </button>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
.account-settings {
  padding-block: var(--gap);
}

.account-settings__card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
}

.account-settings__section {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
  background: var(--color-surface);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.account-settings__section-title {
  font: 600 16px/20px Roboto, sans-serif;
}

.account-settings__profile {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  padding: 8px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.account-settings__avatar {
  position: relative;
  display: flex;
  flex: none;
  align-items: center;
  justify-content: center;
  width: 72px;
  height: 72px;
  padding: 0;
  color: var(--color-muted);
  background: var(--color-surface);
  border: 0;
  border-radius: 50%;
  overflow: hidden;
  cursor: pointer;
}

.account-settings__avatar:hover {
  color: var(--color-text);
}

.account-settings__avatar-image {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.account-settings__avatar-placeholder {
  width: 48px;
  height: 48px;
}

.account-settings__avatar-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #ffffff;
  background: rgb(0 0 0 / 55%);
  font: 500 11px/13px Roboto, sans-serif;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.account-settings__avatar:hover .account-settings__avatar-overlay,
.account-settings__avatar:focus-visible .account-settings__avatar-overlay {
  opacity: 1;
}

.account-settings__profile-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.account-settings__profile-title {
  font: 500 14px/18px Roboto, sans-serif;
}

.account-settings__profile-hint {
  color: var(--color-muted);
  font: 400 12px/16px Roboto, sans-serif;
}

.account-settings__profile-button {
  margin-left: auto;
  padding: 0;
  color: var(--color-auth-link);
  background: none;
  border: 0;
  font: 400 12px/16px Roboto, sans-serif;
  text-decoration: underline;
  text-underline-offset: 2px;
  cursor: pointer;
}

.account-settings__profile-button:hover {
  opacity: 0.8;
}

.account-settings__profile-button--muted {
  color: var(--color-muted);
}

.account-settings__error--profile {
  flex-basis: 100%;
}

.account-settings__rows {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.account-settings__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  width: 100%;
  padding: 8px;
  color: var(--color-text);
  background: var(--color-card);
  border: 0;
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  font: inherit;
  text-align: left;
}

.account-settings__row--interactive {
  cursor: pointer;
}

.account-settings__row--interactive:hover:not(:disabled) {
  background: var(--color-card-hover);
}

.account-settings__row:disabled {
  cursor: default;
  opacity: 0.6;
}

.account-settings__row--link:hover .account-settings__value {
  text-decoration: underline;
  text-underline-offset: 2px;
}

.account-settings__label {
  display: flex;
  align-items: center;
  gap: 4px;
  font: 500 14px/18px Roboto, sans-serif;
}

.account-settings__icon {
  flex: none;
}

.account-settings__value {
  flex: none;
  color: var(--color-muted);
  font: 400 12px/16px Roboto, sans-serif;
  opacity: 0.75;
  text-align: right;
}

.account-settings__row--interactive:hover .account-settings__value {
  text-decoration: underline;
  text-underline-offset: 2px;
}

.account-settings__row--select {
  flex-wrap: wrap;
}

.account-settings__row--select select {
  min-height: 30px;
  max-width: 100%;
  padding: 4px 8px;
  color: var(--color-text);
  background: var(--color-auth-input);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius);
  font: 400 12px/16px Roboto, sans-serif;
}

.account-settings__divider {
  height: 2px;
  background: var(--color-surface);
  border-radius: var(--radius);
}

.account-settings__notice {
  padding: 8px;
  color: var(--color-warning);
  background: var(--color-card);
  border-radius: var(--radius);
  font: 600 12px/16px Roboto, sans-serif;
}

.account-settings__email-form {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.account-settings__email-input {
  width: 100%;
  min-height: 30px;
  padding: 4px 8px;
  color: var(--color-text);
  background: var(--color-auth-input);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius);
  font: 400 14px/20px Roboto, sans-serif;
}

.account-settings__email-input:hover {
  border-color: var(--color-muted);
}

.account-settings__email-input:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 0;
}

.account-settings__email-input[aria-invalid='true'] {
  border-color: var(--color-danger);
}

.account-settings__error {
  color: var(--color-danger);
  font: 400 12px/16px Roboto, sans-serif;
}

.account-settings__email-actions {
  display: flex;
  gap: 4px;
}

.account-settings__email-button {
  flex: 1;
  min-height: 32px;
  padding: 4px 8px;
  color: var(--color-text);
  background: var(--color-surface);
  border: 0;
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  font: 400 12px/16px Roboto, sans-serif;
  cursor: pointer;
}

.account-settings__email-button:hover:not(:disabled) {
  background: var(--color-card-hover);
}

.account-settings__email-button:disabled {
  cursor: default;
  opacity: 0.6;
}

.account-settings__totp {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.account-settings__totp-hint {
  font: 400 13px/18px Roboto, sans-serif;
}

.account-settings__totp-qr {
  align-self: center;
  width: 160px;
  height: 160px;
  background: #ffffff;
  border-radius: var(--radius);
}

.account-settings__totp-secret {
  color: var(--color-muted);
  font: 400 12px/16px Roboto, sans-serif;
  overflow-wrap: anywhere;
}

.account-settings__totp-input {
  letter-spacing: 4px;
  font: 600 18px/24px Roboto, sans-serif;
  text-align: center;
}

.account-settings__backup {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.account-settings__backup-hint {
  font: 400 13px/18px Roboto, sans-serif;
}

.account-settings__backup-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 4px;
  padding: 8px;
  background: var(--color-surface);
  border-radius: var(--radius);
  font: 600 13px/18px Roboto, sans-serif;
  list-style: none;
}

@media (max-width: 767px) {
  .account-settings {
    padding-block: 16px;
  }

  .account-settings__card {
    gap: 12px;
  }

  .account-settings__section {
    padding: 12px;
  }

  .account-settings__profile {
    display: grid;
    grid-template-columns: auto minmax(0, 1fr);
    align-items: center;
    column-gap: 12px;
    row-gap: 8px;
  }

  .account-settings__avatar {
    grid-row: 1;
    grid-column: 1;
  }

  .account-settings__profile-text {
    grid-row: 1;
    grid-column: 2;
  }

  .account-settings__profile-button {
    grid-column: 1 / -1;
    justify-self: end;
    margin-left: 0;
  }

  .account-settings__error--profile {
    grid-column: 1 / -1;
  }

  .account-settings__row {
    flex-direction: column;
    align-items: stretch;
    gap: 2px;
  }

  .account-settings__label {
    flex: none;
    min-width: 0;
  }

  .account-settings__value {
    flex: none;
    min-width: 0;
    text-align: left;
    overflow-wrap: anywhere;
  }

  .account-settings__row--select select {
    width: 100%;
  }
}
</style>
