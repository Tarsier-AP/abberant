<script setup lang="ts">
import { onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import AuthDemoHint from '@/components/auth/AuthDemoHint.vue'
import RegisterForm from '@/components/auth/RegisterForm.vue'
import { useTwoFactorDemoShortcut } from '@/composables/useTwoFactorDemoShortcut'
import { useAuthStore } from '@/stores/auth'
import type { RegisterCredentials } from '@/types/auth'

const router = useRouter()
const auth = useAuthStore()
const loading = ref(false)
let timer: ReturnType<typeof setTimeout> | undefined

useTwoFactorDemoShortcut()

function handleSubmit(credentials: RegisterCredentials) {
  if (loading.value) return

  auth.startRegistration({
    login: credentials.login,
    username: credentials.username,
  })
  loading.value = true

  // Simulated request until the Go backend and the real auth flow exist.
  timer = setTimeout(() => {
    loading.value = false
    void router.push({ name: 'register-password' })
  }, 700)
}

function handleCancel() {
  if (timer) clearTimeout(timer)

  auth.cancelRegistration()
  loading.value = false
  void router.push({ name: 'login' })
}

onBeforeUnmount(() => {
  if (timer) clearTimeout(timer)
})
</script>

<template>
  <div class="register">
    <RegisterForm :loading="loading" @submit="handleSubmit" @cancel="handleCancel" />
    <AuthDemoHint />
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
