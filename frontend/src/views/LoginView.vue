<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import LoginForm from '@/components/auth/LoginForm.vue'
import { useAuthStore } from '@/stores/auth'
import * as api from '@/api/auth'
import type { LoginCredentials } from '@/types/auth'

const router = useRouter()
const auth = useAuthStore()
const loading = ref(false)
const error = ref('')

async function handleSubmit(credentials: LoginCredentials) {
  if (loading.value) return

  loading.value = true
  error.value = ''

  try {
    const result = await api.login(credentials.login, credentials.password)

    if (result.status === '2fa_required') {
      auth.setPendingLogin(result.twoFactorMethod)
      void router.push({
        name: 'two-factor',
        query: { method: result.twoFactorMethod || undefined },
      })
      return
    }

    auth.clearPendingLogin()
    auth.setAccount(result.account)
    void router.push({ name: 'home' })
  } catch (requestError) {
    error.value = requestError instanceof Error ? requestError.message : 'Ошибка входа.'
  } finally {
    loading.value = false
  }
}

function handleClearError() {
  error.value = ''
}
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
