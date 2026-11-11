<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import type { TwoFactorMethod } from '@/types/auth'

const router = useRouter()
const auth = useAuthStore()

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

const emailInput = ref<HTMLInputElement | null>(null)
const editingEmail = ref(false)
const emailDraft = ref('')
const emailError = ref('')

const isActivated = computed(() => auth.account?.activated === true)
const isEmailLogin = computed(() => EMAIL_PATTERN.test(auth.account?.login ?? ''))

const phone = computed(() => {
  if (!auth.account) return ''

  return isEmailLogin.value ? 'Не указан' : formatPhone(auth.account.login)
})

const primaryEmail = computed(() => {
  if (!auth.account) return ''

  return isEmailLogin.value ? auth.account.login : 'Не указана'
})

const additionalEmail = computed(() => auth.account?.additionalEmail ?? '')

function formatPhone(value: string) {
  const digits = value.replace(/\D/g, '')

  if (digits.length === 11 && (digits.startsWith('7') || digits.startsWith('8'))) {
    const normalized = digits.startsWith('8') ? `7${digits.slice(1)}` : digits

    return `+7 (${normalized.slice(1, 4)}) ${normalized.slice(4, 7)}-${normalized.slice(7, 9)}-${normalized.slice(9, 11)}`
  }

  return value
}

function toggleTwoFactor() {
  if (!auth.account || !isActivated.value) return

  auth.setTwoFactor(!auth.account.twoFactorEnabled)
}

function handleMethodChange(event: Event) {
  auth.setTwoFactorMethod((event.target as HTMLSelectElement).value as TwoFactorMethod)
}

function toggleSmsNotifications() {
  if (!auth.account) return

  auth.setSmsNotifications(!auth.account.smsNotifications)
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

  auth.setAdditionalEmail(email)
  editingEmail.value = false
  emailError.value = ''
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

        <div class="account-settings__rows">
          <div class="account-settings__row">
            <span class="account-settings__label">Телефон:</span>
            <span class="account-settings__value">{{ phone }}</span>
          </div>

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
            <span class="account-settings__value">{{ primaryEmail }}</span>
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
            :disabled="!isActivated"
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
            <select :value="auth.account.twoFactorMethod" @change="handleMethodChange">
              <option value="sms">SMS</option>
              <option value="email">Электронная почта</option>
              <option value="authenticator">Приложение-аутентификатор</option>
            </select>
          </label>

          <RouterLink
            class="account-settings__row account-settings__row--link"
            to="/password-reset"
          >
            <span class="account-settings__label">Сменить пароль:</span>
            <span class="account-settings__value">Сменить</span>
          </RouterLink>
        </div>

        <p v-if="!isActivated" class="account-settings__notice">
          Активируйте аккаунт, чтобы включить двухфакторную аутентификацию.
        </p>
      </section>

      <section class="account-settings__section" aria-labelledby="account-notifications">
        <h2 id="account-notifications" class="account-settings__section-title">Уведомления</h2>

        <div class="account-settings__rows">
          <button
            class="account-settings__row account-settings__row--interactive"
            type="button"
            role="switch"
            :aria-checked="auth.account.smsNotifications"
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
  max-width: 800px;
  margin-inline: auto;
  padding: 8px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
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

.account-settings__email-button:hover {
  background: var(--color-card-hover);
}

@media (max-width: 767px) {
  .account-settings {
    padding-block: 16px;
  }

  .account-settings__row {
    align-items: flex-start;
  }
}
</style>
