<script setup lang="ts">
import { onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import LoginForm from '@/components/auth/LoginForm.vue'
import { useAuthStore } from '@/stores/auth'
import type { LoginCredentials } from '@/types/auth'

const router = useRouter()
const auth = useAuthStore()
const loading = ref(false)
const error = ref('')
let timer: ReturnType<typeof setTimeout> | undefined

function handleSubmit(credentials: LoginCredentials) {
  if (loading.value) return

  loading.value = true
  error.value = ''

  // Simulated request until the Go backend and the real auth flow exist.
  timer = setTimeout(() => {
    loading.value = false

    if (!auth.matchesAccount(credentials.login)) {
      error.value = 'Неверный логин или пароль.'
      return
    }

    if (!auth.account?.activated) {
      error.value = 'Аккаунт не активирован. Перейдите по ссылке из сообщения.'
      return
    }

    if (auth.canUseTwoFactor) {
      void router.push({
        name: 'two-factor',
        query: { method: auth.account?.twoFactorMethod ?? 'sms' },
      })
      return
    }

    void router.push({ name: 'home' })
  }, 700)
}

function handleClearError() {
  error.value = ''
}

onBeforeUnmount(() => {
  if (timer) clearTimeout(timer)
})
</script>

<template>
  <div class="login">
    <LoginForm
      :loading="loading"
      :error="error"
      @submit="handleSubmit"
      @clear-error="handleClearError"
    />
  </div>
</template>

<style scoped>
.login {
  position: relative;
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
}
</style>
