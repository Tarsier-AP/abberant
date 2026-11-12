<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import RegisterForm from '@/components/auth/RegisterForm.vue'
import { useAuthStore } from '@/stores/auth'
import * as api from '@/api/auth'
import type { RegisterCredentials } from '@/types/auth'

const router = useRouter()
const auth = useAuthStore()
const loading = ref(false)
const error = ref('')

async function handleSubmit(credentials: RegisterCredentials) {
  if (loading.value) return

  loading.value = true
  error.value = ''

  try {
    const { registrationId } = await api.startRegistration(credentials.login, credentials.username)

    auth.startRegistration({
      login: credentials.login,
      username: credentials.username,
    })
    auth.completeStartRegistration(registrationId)

    void router.push({ name: 'register-password' })
  } catch (requestError) {
    error.value =
      requestError instanceof Error ? requestError.message : 'Не удалось начать регистрацию.'
  } finally {
    loading.value = false
  }
}

function handleCancel() {
  auth.cancelRegistration()
  loading.value = false
  void router.push({ name: 'login' })
}
</script>

<template>
  <div class="register">
    <RegisterForm :loading="loading" :error="error" @submit="handleSubmit" @cancel="handleCancel" />
  </div>
</template>

<style scoped>
.register {
  position: relative;
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
}
</style>
