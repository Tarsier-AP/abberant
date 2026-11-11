<script setup lang="ts">
import { onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import LoginForm from '@/components/auth/LoginForm.vue'
import type { LoginCredentials } from '@/types/auth'

const router = useRouter()
const loading = ref(false)
let timer: ReturnType<typeof setTimeout> | undefined

function handleSubmit(credentials: LoginCredentials) {
  if (loading.value) return

  loading.value = true

  // Simulated request until the Go backend and the real auth flow exist.
  timer = setTimeout(() => {
    loading.value = false
    void router.push({
      name: 'two-factor',
      state: { login: credentials.login },
    })
  }, 700)
}

onBeforeUnmount(() => {
  if (timer) clearTimeout(timer)
})
</script>

<template>
  <div class="login">
    <LoginForm :loading="loading" @submit="handleSubmit" />
  </div>
</template>

<style scoped>
.login {
  position: relative;
  display: flex;
  flex: 1;
  align-items: flex-start;
  justify-content: center;
  padding: 32px 16px 56px;
}
</style>
