<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import type { ComponentPublicInstance } from 'vue'
import AuthCard from './AuthCard.vue'

const props = defineProps<{
  login: string
  username: string
  loading?: boolean
  resending?: boolean
  error?: string
}>()

const emit = defineEmits<{
  submit: [code: string]
  resend: []
}>()

const digits = ref<string[]>(['', '', '', '', '', ''])
const errors = reactive({ code: '' })
const inputRefs = ref<Array<HTMLInputElement | null>>([])

const description = computed(
  () => `Мы отправили код активации на вашу электронную почту ${props.login}.`,
)

const title = computed(() => `Активация аккаунта${props.username ? ` — ${props.username}` : ''}`)

function setInputRef(el: Element | ComponentPublicInstance | null, index: number) {
  inputRefs.value[index] = el instanceof HTMLInputElement ? el : null
}

function focusInput(index: number) {
  inputRefs.value[index]?.focus()
  inputRefs.value[index]?.select()
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
    errors.code = 'Введите 6-значный код из сообщения.'
    const firstEmpty = digits.value.findIndex((digit) => !digit)

    focusInput(firstEmpty === -1 ? 0 : firstEmpty)
    return
  }

  emit('submit', code)
}
</script>

<template>
  <AuthCard :title="title" title-shadow>
    <form class="account-activation" novalidate @submit.prevent="handleSubmit">
      <p id="activation-description" class="account-activation__text">
        {{ description }}
      </p>

      <div
        class="account-activation__code"
        role="group"
        aria-labelledby="activation-description"
        @paste="handlePaste"
      >
        <input
          v-for="(digit, index) in digits"
          :key="index"
          :ref="(el) => setInputRef(el, index)"
          :value="digit"
          class="account-activation__digit"
          type="text"
          inputmode="numeric"
          pattern="[0-9]*"
          maxlength="1"
          :autocomplete="index === 0 ? 'one-time-code' : 'off'"
          :aria-label="`Цифра ${index + 1}`"
          :aria-invalid="errors.code || error ? 'true' : undefined"
          :aria-describedby="errors.code || error ? 'activation-error' : undefined"
          :disabled="loading"
          @input="handleInput(index, $event)"
          @keydown="handleKeydown(index, $event)"
        />
      </div>

      <p v-if="error || errors.code" id="activation-error" class="account-activation__warning" role="alert">
        {{ error || errors.code }}
      </p>

      <p class="account-activation__resend">
        <span>Письмо не пришло? Проверьте папку «Спам».</span>
        <button
          class="account-activation__resend-link"
          type="button"
          :disabled="loading || resending"
          @click="emit('resend')"
        >
          Отправить снова
        </button>
      </p>

      <p class="account-activation__note">Код действует 5 минут и может быть использован один раз.</p>

      <button class="account-activation__submit" type="submit" :disabled="loading">
        {{ loading ? 'Проверяем…' : 'Активировать' }}
      </button>
    </form>
  </AuthCard>
</template>

<style scoped>
.account-activation {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 8px;
}

.account-activation__text {
  font: 400 14px/20px Roboto, sans-serif;
}

.account-activation__code {
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 4px;
  padding: 0 8px;
  background: var(--color-card);
  border-radius: var(--radius);
}

.account-activation__digit {
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

.account-activation__digit:focus-visible {
  border-color: var(--color-accent);
  outline: 2px solid var(--color-accent);
  outline-offset: 0;
}

.account-activation__digit[aria-invalid='true'] {
  border-color: var(--color-danger);
}

.account-activation__digit:disabled {
  opacity: 0.7;
}

.account-activation__warning {
  padding: 8px;
  color: var(--color-warning);
  background: var(--color-surface);
  border-radius: var(--radius);
  font: 600 12px/16px Roboto, sans-serif;
}

.account-activation__resend {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 0 8px;
  font: 400 13px/16px Roboto, sans-serif;
  text-align: center;
}

.account-activation__resend-link {
  padding: 0;
  color: var(--color-auth-link);
  background: none;
  border: 0;
  font: 500 13px/16px Roboto, sans-serif;
  text-decoration: underline;
  text-underline-offset: 2px;
  cursor: pointer;
}

.account-activation__resend-link:hover:not(:disabled) {
  opacity: 0.8;
}

.account-activation__resend-link:disabled {
  cursor: default;
  opacity: 0.6;
}

.account-activation__note {
  color: var(--color-muted);
  font: 400 12px/16px Roboto, sans-serif;
  opacity: 0.7;
}

.account-activation__submit {
  display: flex;
  align-items: center;
  justify-content: center;
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

.account-activation__submit:hover:not(:disabled) {
  background: var(--color-card-hover);
}

.account-activation__submit:disabled {
  cursor: default;
  opacity: 0.6;
}
</style>
