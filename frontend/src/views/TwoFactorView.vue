<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import TwoFactorForm from '@/components/auth/TwoFactorForm.vue'
import { useAuthStore } from '@/stores/auth'
import * as api from '@/api/auth'
import type { TwoFactorMethod } from '@/types/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const loading = ref(false)
const error = ref('')

const allowed = computed(() => auth.pendingLogin !== null)

const method = computed<TwoFactorMethod | ''>(() => {
  const value = route.query.method

  if (value === 'authenticator' || value === 'email') return value

  return auth.pendingLogin?.method ?? auth.account?.twoFactorMethod ?? ''
})

function handleSubmit(code: string) {
  if (loading.value) return

  loading.value = true
  error.value = ''

  void confirmCode(code)
}

async function confirmCode(code: string) {
  try {
    const { account } = await api.completeTwoFactorLogin(code)

    auth.setAccount(account)
    auth.clearPendingLogin()
    void router.push({ name: 'home' })
  } catch (requestError) {
    error.value = requestError instanceof Error ? requestError.message : 'Ошибка подтверждения.'
    loading.value = false
  }
}

async function handleResend() {
  if (loading.value) return

  error.value = ''

  try {
    await api.sendEmailTwoFactorCode()
  } catch (requestError) {
    error.value = requestError instanceof Error ? requestError.message : 'Не удалось отправить код.'
  }
}

function handleCancel() {
  loading.value = false
  auth.clearPendingLogin()
  void router.push({ name: 'login' })
}
</script>

<template>
  <div v-if="allowed" class="two-factor">
    <TwoFactorForm
      :method="method || undefined"
      :loading="loading"
      :error="error"
      :resendable="method === 'email'"
      @submit="handleSubmit"
      @resend="handleResend"
      @cancel="handleCancel"
    />
  </div>
</template>

<style scoped>
.two-factor {
  position: relative;
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
}
</style>
