<script setup lang="ts">
import { reactive, ref } from 'vue'
import AuthCard from './AuthCard.vue'

defineProps<{
  loading?: boolean
}>()

const emit = defineEmits<{
  submit: [password: string]
  cancel: []
}>()

const form = reactive({
  password: '',
  confirmation: '',
})

const errors = reactive({
  password: '',
  confirmation: '',
})

const touched = reactive({
  password: false,
  confirmation: false,
})

const passwordInput = ref<HTMLInputElement | null>(null)
const confirmationInput = ref<HTMLInputElement | null>(null)

const requirements = [
  'Минимальная длина: 8 символов',
  'Хотя бы одна заглавной буквы',
  'Хотя бы одна цифра',
  'Хотя бы один специальный символ (например, !@#$)',
]

const UPPERCASE_PATTERN = /[A-ZА-ЯЁ]/
const DIGIT_PATTERN = /\d/
const SPECIAL_PATTERN = /[^A-Za-zА-Яа-яЁё0-9\s]/

function validatePassword(value: string): string {
  if (!value) return 'Введите пароль.'
  if (value.length < 8) return 'Пароль должен содержать не менее 8 символов.'
  if (!UPPERCASE_PATTERN.test(value)) return 'Добавьте хотя бы одну заглавную букву.'
  if (!DIGIT_PATTERN.test(value)) return 'Добавьте хотя бы одну цифру.'
  if (!SPECIAL_PATTERN.test(value)) return 'Добавьте хотя бы один специальный символ.'

  return ''
}

function validateConfirmation(value: string, password = form.password): string {
  if (!value) return 'Повторите пароль.'
  if (value !== password) return 'Пароли не совпадают.'

  return ''
}

function handlePasswordBlur() {
  touched.password = true
  errors.password = validatePassword(form.password)

  if (touched.confirmation && form.confirmation) {
    errors.confirmation = validateConfirmation(form.confirmation)
  }
}

function handleConfirmationBlur() {
  touched.confirmation = true
  errors.confirmation = validateConfirmation(form.confirmation)
}

function handlePasswordInput(event: Event) {
  if (!touched.password) return

  const value = (event.target as HTMLInputElement).value
  errors.password = validatePassword(value)

  if (touched.confirmation && form.confirmation) {
    errors.confirmation = validateConfirmation(form.confirmation, value)
  }
}

function handleConfirmationInput(event: Event) {
  if (!touched.confirmation) return

  const value = (event.target as HTMLInputElement).value
  errors.confirmation = validateConfirmation(value)
}

function handleSubmit() {
  touched.password = true
  touched.confirmation = true
  errors.password = validatePassword(form.password)
  errors.confirmation = validateConfirmation(form.confirmation)

  if (errors.password) {
    passwordInput.value?.focus()
    return
  }

  if (errors.confirmation) {
    confirmationInput.value?.focus()
    return
  }

  emit('submit', form.password)
}
</script>

<template>
  <AuthCard title="Создание аккаунта" title-shadow>
    <form class="register-password-form" novalidate @submit.prevent="handleSubmit">
      <div class="register-password-form__field">
        <label class="register-password-form__label" for="register-password">Пароль:</label>
        <input
          id="register-password"
          ref="passwordInput"
          v-model="form.password"
          class="register-password-form__input"
          type="password"
          name="password"
          autocomplete="new-password"
          :aria-invalid="errors.password ? 'true' : undefined"
          :aria-describedby="
            errors.password
              ? 'register-password-error register-password-requirements'
              : 'register-password-requirements'
          "
          :disabled="loading"
          @blur="handlePasswordBlur"
          @input="handlePasswordInput"
        />
        <p
          v-if="errors.password"
          id="register-password-error"
          class="register-password-form__error"
          role="alert"
        >
          {{ errors.password }}
        </p>
      </div>

      <div id="register-password-requirements" class="register-password-form__requirements">
        <p class="register-password-form__requirements-title">Требования к паролю:</p>
        <ul class="register-password-form__requirements-list">
          <li
            v-for="requirement in requirements"
            :key="requirement"
            class="register-password-form__requirement"
          >
            {{ requirement }}
          </li>
        </ul>
      </div>

      <div class="register-password-form__field">
        <label class="register-password-form__label" for="register-password-confirmation">
          Повторите пароль:
        </label>
        <input
          id="register-password-confirmation"
          ref="confirmationInput"
          v-model="form.confirmation"
          class="register-password-form__input"
          type="password"
          name="password-confirmation"
          autocomplete="new-password"
          :aria-invalid="errors.confirmation ? 'true' : undefined"
          :aria-describedby="errors.confirmation ? 'register-password-confirmation-error' : undefined"
          :disabled="loading"
          @blur="handleConfirmationBlur"
          @input="handleConfirmationInput"
        />
        <p
          v-if="errors.confirmation"
          id="register-password-confirmation-error"
          class="register-password-form__error"
          role="alert"
        >
          {{ errors.confirmation }}
        </p>
      </div>

      <div class="register-password-form__actions">
        <button
          class="register-password-form__button"
          type="button"
          :disabled="loading"
          @click="emit('cancel')"
        >
          Отменить
        </button>
        <button class="register-password-form__button" type="submit" :disabled="loading">
          {{ loading ? 'Регистрируем…' : 'Зарегистрироваться' }}
        </button>
      </div>
    </form>
  </AuthCard>
</template>

<style scoped>
.register-password-form {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.register-password-form__field {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
  border-radius: 4px;
}

.register-password-form__label {
  font: 600 14px/18px Roboto, sans-serif;
}

.register-password-form__input {
  width: 100%;
  min-height: 30px;
  padding: 4px 8px;
  color: var(--color-text);
  background: var(--color-auth-input);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius);
  font: 400 14px/20px Roboto, sans-serif;
}

.register-password-form__input:hover:not(:disabled) {
  border-color: var(--color-muted);
}

.register-password-form__input:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 0;
}

.register-password-form__input[aria-invalid='true'] {
  border-color: var(--color-danger);
}

.register-password-form__input:disabled {
  opacity: 0.7;
}

.register-password-form__requirements {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
}

.register-password-form__requirements-title {
  color: var(--color-text);
  font: 700 12px/16px Roboto, sans-serif;
  opacity: 0.5;
}

.register-password-form__requirements-list {
  display: flex;
  flex-direction: column;
  gap: 5px;
  padding: 5px;
}

.register-password-form__requirement {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 0 2px 0 5px;
  font: 700 12px/16px Roboto, sans-serif;
}

.register-password-form__requirement::before {
  content: '';
  flex: none;
  width: 7px;
  height: 7px;
  background: var(--color-text);
  border-radius: 50%;
}

.register-password-form__error {
  color: var(--color-danger);
  font: 400 12px/16px Roboto, sans-serif;
}

.register-password-form__actions {
  display: flex;
  gap: 4px;
  padding: 8px;
}

.register-password-form__button {
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

.register-password-form__button:hover:not(:disabled) {
  background: var(--color-card-hover);
}

.register-password-form__button:disabled {
  cursor: default;
  opacity: 0.6;
}
</style>
