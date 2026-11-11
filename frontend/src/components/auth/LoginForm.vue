<script setup lang="ts">
import { reactive, ref } from 'vue'
import type { LoginCredentials } from '@/types/auth'
import AuthCard from './AuthCard.vue'

defineProps<{
  loading?: boolean
}>()

const emit = defineEmits<{
  submit: [credentials: LoginCredentials]
}>()

const form = reactive<LoginCredentials>({
  login: '',
  password: '',
})

const errors = reactive({
  login: '',
  password: '',
})

const touched = reactive({
  login: false,
  password: false,
})

const loginInput = ref<HTMLInputElement | null>(null)
const passwordInput = ref<HTMLInputElement | null>(null)

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const PHONE_PATTERN = /^\+?[\d\s()-]+$/

function validateLogin(value: string): string {
  const login = value.trim()

  if (!login) return 'Введите электронную почту или номер телефона.'

  if (EMAIL_PATTERN.test(login)) return ''

  const digits = login.replace(/\D/g, '')

  if (PHONE_PATTERN.test(login) && digits.length >= 10 && digits.length <= 15) return ''

  return 'Введите электронную почту или номер телефона в корректном формате.'
}

function validatePassword(value: string): string {
  if (!value) return 'Введите пароль.'
  if (value.length < 6) return 'Пароль должен содержать не менее 6 символов.'

  return ''
}

function handleBlur(field: 'login' | 'password') {
  touched[field] = true
  errors[field] = field === 'login' ? validateLogin(form.login) : validatePassword(form.password)
}

function handleInput(field: 'login' | 'password', event: Event) {
  if (!touched[field]) return

  const value = (event.target as HTMLInputElement).value
  errors[field] = field === 'login' ? validateLogin(value) : validatePassword(value)
}

function handleSubmit() {
  touched.login = true
  touched.password = true
  errors.login = validateLogin(form.login)
  errors.password = validatePassword(form.password)

  if (errors.login) {
    loginInput.value?.focus()
    return
  }

  if (errors.password) {
    passwordInput.value?.focus()
    return
  }

  emit('submit', {
    login: form.login.trim(),
    password: form.password,
  })
}
</script>

<template>
  <AuthCard title="Вход в аккаунт">
    <form class="login-form" novalidate @submit.prevent="handleSubmit">
      <div class="login-form__field">
        <label class="login-form__label" for="login">
          Логин (электронная почта или номер телефона):
        </label>
        <input
          id="login"
          ref="loginInput"
          v-model="form.login"
          class="login-form__input"
          type="text"
          name="login"
          autocomplete="username"
          inputmode="email"
          :aria-invalid="errors.login ? 'true' : undefined"
          :aria-describedby="errors.login ? 'login-error' : undefined"
          :disabled="loading"
          @blur="handleBlur('login')"
          @input="handleInput('login', $event)"
        />
        <p v-if="errors.login" id="login-error" class="login-form__error" role="alert">
          {{ errors.login }}
        </p>
      </div>

      <div class="login-form__field">
        <label class="login-form__label" for="password">Пароль</label>
        <input
          id="password"
          ref="passwordInput"
          v-model="form.password"
          class="login-form__input"
          type="password"
          name="password"
          autocomplete="current-password"
          :aria-invalid="errors.password ? 'true' : undefined"
          :aria-describedby="errors.password ? 'password-error' : undefined"
          :disabled="loading"
          @blur="handleBlur('password')"
          @input="handleInput('password', $event)"
        />
        <p v-if="errors.password" id="password-error" class="login-form__error" role="alert">
          {{ errors.password }}
        </p>
      </div>

      <div class="login-form__links">
        <RouterLink class="login-form__link" to="/password-reset">Забыли пароль?</RouterLink>
        <RouterLink class="login-form__link" to="/register">Зарегистрируйтесь</RouterLink>
      </div>

      <div class="login-form__actions">
        <button class="login-form__submit" type="submit" :disabled="loading">
          {{ loading ? 'Проверяем…' : 'Далее' }}
        </button>
      </div>
    </form>
  </AuthCard>
</template>

<style scoped>
.login-form {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.login-form__field {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
  border-radius: 4px;
}

.login-form__label {
  font: 600 14px/18px Roboto, sans-serif;
}

.login-form__input {
  width: 100%;
  min-height: 30px;
  padding: 4px 8px;
  color: var(--color-text);
  background: var(--color-card);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius);
  font: 400 14px/20px Roboto, sans-serif;
}

.login-form__input:hover:not(:disabled) {
  border-color: var(--color-muted);
}

.login-form__input:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 0;
}

.login-form__input[aria-invalid='true'] {
  border-color: var(--color-danger);
}

.login-form__input:disabled {
  opacity: 0.7;
}

.login-form__error {
  color: var(--color-danger);
  font: 400 12px/16px Roboto, sans-serif;
}

.login-form__links {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 8px;
}

.login-form__link {
  font: 700 13px/16px Roboto, sans-serif;
  text-align: center;
}

.login-form__link:hover {
  text-decoration: underline;
}

.login-form__actions {
  display: flex;
  gap: 4px;
  padding: 8px;
}

.login-form__submit {
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

.login-form__submit:hover:not(:disabled) {
  background: var(--color-card-hover);
}

.login-form__submit:disabled {
  cursor: default;
  opacity: 0.6;
}
</style>
