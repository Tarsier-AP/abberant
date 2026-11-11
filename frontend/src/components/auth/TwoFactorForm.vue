<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import type { ComponentPublicInstance } from 'vue'
import type { TwoFactorMethod } from '@/types/auth'
import AuthCard from './AuthCard.vue'

const props = defineProps<{
  method?: TwoFactorMethod
  loading?: boolean
}>()

const emit = defineEmits<{
  submit: [code: string]
  cancel: []
}>()

const digits = ref<string[]>(['', '', '', '', '', ''])
const errors = reactive({
  code: '',
})

const inputRefs = ref<Array<HTMLInputElement | null>>([])

const description = computed(() => {
  switch (props.method ?? 'sms') {
    case 'authenticator':
      return 'Введите 6-значный код из приложения-аутентификатора:'
    case 'email':
      return 'Введите 6-значный код, который мы отправили на вашу электронную почту:'
    default:
      return 'Введите 6-значный код, который мы отправили в SMS:'
  }
})

function setInputRef(el: Element | ComponentPublicInstance | null, index: number) {
  inputRefs.value[index] = el instanceof HTMLInputElement ? el : null
}

function focusInput(index: number) {
  const input = inputRefs.value[index]

  input?.focus()
  input?.select()
}

function handleInput(index: number, event: Event) {
  const input = event.target as HTMLInputElement
  const value = input.value.replace(/\D/g, '').slice(-1)

  digits.value[index] = value
  input.value = value
  errors.code = ''

  if (value && index < digits.value.length - 1) {
    focusInput(index + 1)
  }
}

function handleKeydown(index: number, event: KeyboardEvent) {
  if (event.key === 'Backspace' && !digits.value[index] && index > 0) {
    event.preventDefault()
    digits.value[index - 1] = ''
    focusInput(index - 1)
    return
  }

  if (event.key === 'ArrowLeft' && index > 0) {
    event.preventDefault()
    focusInput(index - 1)
  }

  if (event.key === 'ArrowRight' && index < digits.value.length - 1) {
    event.preventDefault()
    focusInput(index + 1)
  }
}

function handlePaste(event: ClipboardEvent) {
  const pasted = event.clipboardData?.getData('text').replace(/\D/g, '').slice(0, 6) ?? ''

  if (!pasted) return

  event.preventDefault()

  digits.value = Array.from({ length: 6 }, (_, index) => pasted[index] ?? '')
  errors.code = ''
  focusInput(Math.min(pasted.length, digits.value.length - 1))
}

function handleSubmit() {
  const code = digits.value.join('')

  if (!/^\d{6}$/.test(code)) {
    errors.code = 'Введите 6-значный код.'
    const firstEmpty = digits.value.findIndex((digit) => !digit)

    focusInput(firstEmpty === -1 ? 0 : firstEmpty)
    return
  }

  emit('submit', code)
}
</script>

<template>
  <AuthCard title="Подтвердите код" title-shadow>
    <form class="two-factor-form" novalidate @submit.prevent="handleSubmit">
      <p id="two-factor-description" class="two-factor-form__description">
        {{ description }}
      </p>

      <div
        class="two-factor-form__code"
        role="group"
        aria-labelledby="two-factor-description"
        @paste="handlePaste"
      >
        <input
          v-for="(digit, index) in digits"
          :key="index"
          :ref="(el) => setInputRef(el, index)"
          :value="digit"
          class="two-factor-form__digit"
          type="text"
          inputmode="numeric"
          pattern="[0-9]*"
          maxlength="1"
          :autocomplete="index === 0 ? 'one-time-code' : 'off'"
          :aria-label="`Цифра ${index + 1}`"
          :aria-invalid="errors.code ? 'true' : undefined"
          :aria-describedby="errors.code ? 'two-factor-code-error' : undefined"
          :disabled="loading"
          @input="handleInput(index, $event)"
          @keydown="handleKeydown(index, $event)"
        />
      </div>

      <p v-if="errors.code" id="two-factor-code-error" class="two-factor-form__error" role="alert">
        {{ errors.code }}
      </p>

      <div class="two-factor-form__help">
        <span>Потеряли доступ к 2FA?</span>
        <RouterLink class="two-factor-form__help-link" to="/help">Напишите нам</RouterLink>
      </div>

      <div class="two-factor-form__actions">
        <button
          class="two-factor-form__button"
          type="button"
          :disabled="loading"
          @click="emit('cancel')"
        >
          Отменить
        </button>
        <button class="two-factor-form__button" type="submit" :disabled="loading">
          {{ loading ? 'Проверяем…' : 'Подтвердить' }}
        </button>
      </div>
    </form>
  </AuthCard>
</template>

<style scoped>
.two-factor-form {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.two-factor-form__description {
  font: 600 14px/18px Roboto, sans-serif;
}

.two-factor-form__code {
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 4px;
  padding: 0 8px;
  background: var(--color-card);
  border-radius: var(--radius);
}

.two-factor-form__digit {
  width: 100%;
  height: 56px;
  padding: 0;
  color: var(--color-text);
  background: var(--color-auth-input);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius);
  font: 600 22px/1 Roboto, sans-serif;
  text-align: center;
}

.two-factor-form__digit:focus-visible {
  border-color: var(--color-accent);
  outline: 2px solid var(--color-accent);
  outline-offset: 0;
}

.two-factor-form__digit[aria-invalid='true'] {
  border-color: var(--color-danger);
}

.two-factor-form__digit:disabled {
  opacity: 0.7;
}

.two-factor-form__error {
  color: var(--color-danger);
  font: 400 12px/16px Roboto, sans-serif;
}

.two-factor-form__help {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 8px;
  font: 400 13px/16px Roboto, sans-serif;
  text-align: center;
}

.two-factor-form__help-link {
  font-weight: 500;
}

.two-factor-form__help-link:hover {
  text-decoration: underline;
}

.two-factor-form__actions {
  display: flex;
  gap: 8px;
  padding: 8px;
}

.two-factor-form__button {
  flex: 1;
  min-height: 36px;
  padding: 8px;
  color: var(--color-text);
  background: var(--color-card);
  border: 0;
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  font: 400 14px/18px Roboto, sans-serif;
  cursor: pointer;
}

.two-factor-form__button:hover:not(:disabled) {
  background: var(--color-card-hover);
}

.two-factor-form__button:disabled {
  cursor: default;
  opacity: 0.6;
}
</style>
