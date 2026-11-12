<script setup lang="ts">
import { reactive, ref } from 'vue'
import type { RegisterCredentials } from '@/types/auth'
import AuthCard from './AuthCard.vue'

defineProps<{
  loading?: boolean
  error?: string
}>()

const emit = defineEmits<{
  submit: [credentials: RegisterCredentials]
  cancel: []
}>()

const form = reactive<RegisterCredentials>({
  login: '',
  username: '',
  acceptedTerms: false,
  acceptedPolicy: false,
})

const errors = reactive({
  login: '',
  username: '',
  terms: '',
  policy: '',
})

const touched = reactive({
  login: false,
  username: false,
})

const loginInput = ref<HTMLInputElement | null>(null)
const usernameInput = ref<HTMLInputElement | null>(null)
const termsInput = ref<HTMLInputElement | null>(null)
const policyInput = ref<HTMLInputElement | null>(null)

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const USERNAME_PATTERN = /^[\wА-Яа-яЁё-]+$/

function validateLogin(value: string): string {
  const login = value.trim()

  if (!login) return 'Введите электронную почту.'

  if (EMAIL_PATTERN.test(login)) return ''

  return 'Введите электронную почту в корректном формате.'
}

function validateUsername(value: string): string {
  const username = value.trim()

  if (!username) return 'Введите имя пользователя.'

  if (username.length < 3 || username.length > 20) {
    return 'Имя пользователя должно содержать от 3 до 20 символов.'
  }

  if (!USERNAME_PATTERN.test(username)) {
    return 'Используйте буквы, цифры, дефис или подчёркивание.'
  }

  return ''
}

function handleBlur(field: 'login' | 'username') {
  touched[field] = true

  errors[field] =
    field === 'login' ? validateLogin(form.login) : validateUsername(form.username)
}

function handleInput(field: 'login' | 'username', event: Event) {
  if (!touched[field]) return

  const value = (event.target as HTMLInputElement).value

  errors[field] = field === 'login' ? validateLogin(value) : validateUsername(value)
}

function clearConsentError(field: 'terms' | 'policy') {
  if (field === 'terms' && form.acceptedTerms) errors.terms = ''
  if (field === 'policy' && form.acceptedPolicy) errors.policy = ''
}

function handleSubmit() {
  touched.login = true
  touched.username = true
  errors.login = validateLogin(form.login)
  errors.username = validateUsername(form.username)
  errors.terms = form.acceptedTerms ? '' : 'Примите пользовательское соглашение.'
  errors.policy = form.acceptedPolicy ? '' : 'Примите политику обработки данных.'

  if (errors.login) {
    loginInput.value?.focus()
    return
  }

  if (errors.username) {
    usernameInput.value?.focus()
    return
  }

  if (errors.terms) {
    termsInput.value?.focus()
    return
  }

  if (errors.policy) {
    policyInput.value?.focus()
    return
  }

  emit('submit', {
    login: form.login.trim(),
    username: form.username.trim(),
    acceptedTerms: form.acceptedTerms,
    acceptedPolicy: form.acceptedPolicy,
  })
}
</script>

<template>
  <AuthCard title="Создание аккаунта" title-shadow>
    <form class="register-form" novalidate @submit.prevent="handleSubmit">
      <div class="register-form__field">
        <label class="register-form__label" for="register-login">
          Электронная почта:
        </label>
        <input
          id="register-login"
          ref="loginInput"
          v-model="form.login"
          class="register-form__input"
          type="email"
          name="login"
          autocomplete="email"
          inputmode="email"
          :aria-invalid="errors.login ? 'true' : undefined"
          :aria-describedby="errors.login ? 'register-login-error' : undefined"
          :disabled="loading"
          @blur="handleBlur('login')"
          @input="handleInput('login', $event)"
        />
        <p
          v-if="errors.login"
          id="register-login-error"
          class="register-form__error"
          role="alert"
        >
          {{ errors.login }}
        </p>
      </div>

      <div class="register-form__field">
        <label class="register-form__label" for="register-username">Имя пользователя:</label>
        <input
          id="register-username"
          ref="usernameInput"
          v-model="form.username"
          class="register-form__input"
          type="text"
          name="username"
          autocomplete="nickname"
          :aria-invalid="errors.username ? 'true' : undefined"
          :aria-describedby="errors.username ? 'register-username-error' : undefined"
          :disabled="loading"
          @blur="handleBlur('username')"
          @input="handleInput('username', $event)"
        />
        <p
          v-if="errors.username"
          id="register-username-error"
          class="register-form__error"
          role="alert"
        >
          {{ errors.username }}
        </p>
      </div>

      <div class="register-form__consents">
        <div class="register-form__consent">
          <input
            id="register-terms"
            ref="termsInput"
            v-model="form.acceptedTerms"
            class="register-form__checkbox"
            type="checkbox"
            name="terms"
            :aria-invalid="errors.terms ? 'true' : undefined"
            :aria-describedby="errors.terms ? 'register-terms-error' : undefined"
            :disabled="loading"
            @change="clearConsentError('terms')"
          />
          <label class="register-form__consent-text" for="register-terms">
            Я принимаю
            <RouterLink class="register-form__consent-link" to="/terms">
              пользовательское соглашение
            </RouterLink>
          </label>
        </div>
        <p
          v-if="errors.terms"
          id="register-terms-error"
          class="register-form__error"
          role="alert"
        >
          {{ errors.terms }}
        </p>

        <div class="register-form__consent">
          <input
            id="register-policy"
            ref="policyInput"
            v-model="form.acceptedPolicy"
            class="register-form__checkbox"
            type="checkbox"
            name="policy"
            :aria-invalid="errors.policy ? 'true' : undefined"
            :aria-describedby="errors.policy ? 'register-policy-error' : undefined"
            :disabled="loading"
            @change="clearConsentError('policy')"
          />
          <label class="register-form__consent-text" for="register-policy">
            Я принимаю
            <RouterLink class="register-form__consent-link" to="/personal-data">
              политику обработки данных
            </RouterLink>
          </label>
        </div>
        <p
          v-if="errors.policy"
          id="register-policy-error"
          class="register-form__error"
          role="alert"
        >
          {{ errors.policy }}
        </p>
      </div>

      <p v-if="error" class="register-form__error register-form__error--form" role="alert">
        {{ error }}
      </p>

      <div class="register-form__actions">
        <button
          class="register-form__button"
          type="button"
          :disabled="loading"
          @click="emit('cancel')"
        >
          Отменить
        </button>
        <button class="register-form__button" type="submit" :disabled="loading">
          {{ loading ? 'Создаём…' : 'Далее' }}
        </button>
      </div>
    </form>
  </AuthCard>
</template>

<style scoped>
.register-form {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.register-form__field {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
  border-radius: 4px;
}

.register-form__label {
  font: 600 14px/18px Roboto, sans-serif;
}

.register-form__input {
  width: 100%;
  min-height: 30px;
  padding: 4px 8px;
  color: var(--color-text);
  background: var(--color-auth-input);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius);
  font: 400 14px/20px Roboto, sans-serif;
}

.register-form__input:hover:not(:disabled) {
  border-color: var(--color-muted);
}

.register-form__input:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 0;
}

.register-form__input[aria-invalid='true'] {
  border-color: var(--color-danger);
}

.register-form__input:disabled {
  opacity: 0.7;
}

.register-form__consents {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 8px;
  border-radius: 4px;
}

.register-form__consent {
  display: flex;
  align-items: flex-start;
  gap: 4px;
  padding: 8px;
}

.register-form__checkbox {
  position: relative;
  flex: none;
  width: 16px;
  height: 16px;
  margin: 1px 0 0;
  background: var(--color-auth-input);
  border: 1px solid var(--color-input-border);
  border-radius: 4px;
  appearance: none;
  cursor: pointer;
}

.register-form__checkbox:checked {
  background: var(--color-accent);
  border-color: var(--color-accent);
}

.register-form__checkbox:checked::after {
  content: '';
  position: absolute;
  top: 1px;
  left: 5px;
  width: 4px;
  height: 8px;
  border: solid var(--color-accent-text);
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
}

.register-form__checkbox:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 2px;
}

.register-form__checkbox[aria-invalid='true'] {
  border-color: var(--color-danger);
}

.register-form__checkbox:disabled {
  cursor: default;
  opacity: 0.7;
}

.register-form__consent-text {
  flex: 1;
  font: 400 13px/18px Roboto, sans-serif;
}

.register-form__consent-link {
  color: var(--color-auth-link);
}

.register-form__consent-link:hover {
  text-decoration: underline;
}

.register-form__error {
  color: var(--color-danger);
  font: 400 12px/16px Roboto, sans-serif;
}

.register-form__error--form {
  padding: 0 8px;
}

.register-form__actions {
  display: flex;
  gap: 4px;
  padding: 8px;
}

.register-form__button {
  flex: 1;
  min-height: 32px;
  padding: 4px 8px;
  color: var(--color-text);
  background: var(--color-card);
  border: 0;
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  font: 400 12px/16px Roboto, sans-serif;
  cursor: pointer;
}

.register-form__button:hover:not(:disabled) {
  background: var(--color-card-hover);
}

.register-form__button:disabled {
  cursor: default;
  opacity: 0.6;
}
</style>
