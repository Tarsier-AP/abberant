<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import RegisterPasswordForm from '@/components/auth/RegisterPasswordForm.vue'
import { useAuthStore } from '@/stores/auth'
import * as api from '@/api/auth'

const router = useRouter()
const auth = useAuthStore()
const loading = ref(false)
const error = ref('')

async function handleSubmit(password: string) {
  if (loading.value) return

  const registrationId = auth.pendingRegistration?.registrationId

  if (!registrationId) {
    void router.push({ name: 'register' })
    return
  }

  loading.value = true
  error.value = ''

  try {
    const result = await api.setRegistrationPassword(registrationId, password)

    if (auth.pendingRegistration) {
      // Keep the dev-mode code hint for the activation screen.
      auth.pendingRegistration.devCode = result.devCode
    }

    void router.push({ name: 'register-activation' })
  } catch (requestError) {
    error.value =
      requestError instanceof Error ? requestError.message : 'Не удалось задать пароль.'
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
  <div class="register-password">
    <RegisterPasswordForm
      :loading="loading"
      :error="error"
      @submit="handleSubmit"
      @cancel="handleCancel"
    />
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
