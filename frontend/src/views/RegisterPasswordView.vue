<script setup lang="ts">
import { onBeforeUnmount, ref } from 'vue'
import { useRouter } from 'vue-router'
import AuthDemoHint from '@/components/auth/AuthDemoHint.vue'
import RegisterPasswordForm from '@/components/auth/RegisterPasswordForm.vue'
import { useTwoFactorDemoShortcut } from '@/composables/useTwoFactorDemoShortcut'

const router = useRouter()
const loading = ref(false)
let timer: ReturnType<typeof setTimeout> | undefined

useTwoFactorDemoShortcut()

function handleSubmit() {
  if (loading.value) return

  loading.value = true

  // Simulated request until the Go backend and the real auth flow exist.
  timer = setTimeout(() => {
    loading.value = false
    void router.push({ name: 'two-factor' })
  }, 700)
}

function handleCancel() {
  if (timer) clearTimeout(timer)

  loading.value = false
  void router.push({ name: 'login' })
}

onBeforeUnmount(() => {
  if (timer) clearTimeout(timer)
})
</script>

<template>
  <div class="register-password">
    <RegisterPasswordForm
      :loading="loading"
      @submit="handleSubmit"
      @cancel="handleCancel"
    />
    <AuthDemoHint />
  </div>
</template>

<style scoped>
.register-password {
  position: relative;
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
}
</style>
