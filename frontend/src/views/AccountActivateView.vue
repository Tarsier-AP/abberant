<script setup lang="ts">
import { reactive, ref } from 'vue'
import AuthCard from '@/components/auth/AuthCard.vue'
import * as api from '@/api/auth'

const form = reactive({
  email: '',
  code: '',
})

const errors = reactive({
  email: '',
  code: '',
})

const loading = ref(false)
const error = ref('')
const activated = ref(false)

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

function validateEmail(value: string): string {
  if (!value.trim()) return 'Введите электронную почту.'

  return EMAIL_PATTERN.test(value.trim()) ? '' : 'Введите электронную почту в корректном формате.'
}

function validateCode(value: string): string {
  return /^\d{6}$/.test(value) ? '' : 'Введите 6-значный код из сообщения.'
}

function handleSubmit() {
  errors.email = validateEmail(form.email)
  errors.code = validateCode(form.code)
  error.value = ''

  if (errors.email || errors.code || loading.value) return

  loading.value = true

  void api
    .activateByEmail(form.email.trim(), form.code)
    .then(() => {
      activated.value = true
    })
    .catch((requestError: unknown) => {
      error.value =
        requestError instanceof Error ? requestError.message : 'Не удалось активировать аккаунт.'
    })
    .finally(() => {
      loading.value = false
    })
}
</script>

<template>
  <div class="activate">
    <AuthCard title="Активация аккаунта" title-shadow>
      <div class="activate__content">
        <template v-if="activated">
          <p class="activate__text">
            Аккаунт активирован. Теперь вы можете войти, указав электронную почту
            <strong>{{ form.email }}</strong
            >.
          </p>

          <section class="activate__settings" aria-labelledby="activate-settings-title">
            <h2 id="activate-settings-title" class="activate__settings-title">
              Настройки аккаунта
            </h2>

            <p class="activate__text">
              В настройках аккаунта можно включить двухфакторную аутентификацию и выбрать
              способ подтверждения.
            </p>

            <RouterLink
              class="activate__settings-link"
              :to="{ name: 'account-settings' }"
            >
              Перейти в настройки
            </RouterLink>
          </section>

          <RouterLink class="activate__submit" :to="{ name: 'login' }">
            Перейти ко входу
          </RouterLink>
        </template>

        <template v-else>
          <p class="activate__text">
            Введите код активации, отправленный на вашу электронную почту.
            Если письмо не пришло, проверьте папку «Спам».
          </p>

          <form class="activate__form" novalidate @submit.prevent="handleSubmit">
            <div class="activate__field">
              <label class="activate__label" for="activate-email">Электронная почта:</label>
              <input
                id="activate-email"
                v-model="form.email"
                class="activate__input"
                type="email"
                autocomplete="email"
                :aria-invalid="errors.email ? 'true' : undefined"
                :aria-describedby="errors.email ? 'activate-email-error' : undefined"
                :disabled="loading"
              />
              <p v-if="errors.email" id="activate-email-error" class="activate__error" role="alert">
                {{ errors.email }}
              </p>
            </div>

            <div class="activate__field">
              <label class="activate__label" for="activate-code">Код активации:</label>
              <input
                id="activate-code"
                v-model="form.code"
                class="activate__input"
                type="text"
                inputmode="numeric"
                maxlength="6"
                :aria-invalid="errors.code ? 'true' : undefined"
                :aria-describedby="errors.code ? 'activate-code-error' : undefined"
                :disabled="loading"
              />
              <p v-if="errors.code" id="activate-code-error" class="activate__error" role="alert">
                {{ errors.code }}
              </p>
            </div>

            <p v-if="error" class="activate__error" role="alert">{{ error }}</p>

            <button class="activate__submit" type="submit" :disabled="loading">
              {{ loading ? 'Проверяем…' : 'Активировать' }}
            </button>
          </form>
        </template>
      </div>
    </AuthCard>
  </div>
</template>

<style scoped>
.activate {
  position: relative;
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
}

.activate__content {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 8px;
}

.activate__text {
  font: 400 14px/20px Roboto, sans-serif;
}

.activate__form {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.activate__field {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
  border-radius: 4px;
}

.activate__label {
  font: 600 14px/18px Roboto, sans-serif;
}

.activate__input {
  width: 100%;
  min-height: 30px;
  padding: 4px 8px;
  color: var(--color-text);
  background: var(--color-auth-input);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius);
  font: 400 14px/20px Roboto, sans-serif;
}

.activate__input:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 0;
}

.activate__input[aria-invalid='true'] {
  border-color: var(--color-danger);
}

.activate__input:disabled {
  opacity: 0.7;
}

.activate__error {
  color: var(--color-danger);
  font: 400 12px/16px Roboto, sans-serif;
}

.activate__settings {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
  background: var(--color-surface);
  border-radius: var(--radius);
}

.activate__settings-title {
  font: 700 14px/18px Roboto, sans-serif;
}

.activate__settings-link {
  align-self: flex-start;
  color: var(--color-auth-link);
  font: 500 14px/18px Roboto, sans-serif;
  text-decoration: underline;
  text-underline-offset: 2px;
}

.activate__settings-link:hover {
  opacity: 0.8;
}

.activate__submit {
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

.activate__submit:hover:not(:disabled) {
  background: var(--color-card-hover);
}

.activate__submit:disabled {
  cursor: default;
  opacity: 0.6;
}
</style>
