<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import TwoFactorForm from '@/components/auth/TwoFactorForm.vue'
import { useAuthStore } from '@/stores/auth'
import type { TwoFactorMethod } from '@/types/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const loading = ref(false)
let timer: ReturnType<typeof setTimeout> | undefined

const isDemo = computed(() => route.query.demo === '1')
const allowed = computed(() => isDemo.value || auth.canUseTwoFactor)

const method = computed<TwoFactorMethod>(() => {
  const value = route.query.method

  if (value === 'authenticator' || value === 'email' || value === 'sms') return value

  return auth.account?.twoFactorMethod ?? 'sms'
})

onMounted(() => {
  if (!allowed.value) {
    void router.replace({ name: 'login' })
  }
})

function handleSubmit() {
  if (loading.value) return

  loading.value = true

  // Simulated verification until the Go backend and the real auth flow exist.
  timer = setTimeout(() => {
    loading.value = false
    void router.push({ name: 'home' })
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
  <div v-if="allowed" class="two-factor">
    <TwoFactorForm
      :method="method"
      :loading="loading"
      @submit="handleSubmit"
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
