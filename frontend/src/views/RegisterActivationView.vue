<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import AuthCard from '@/components/auth/AuthCard.vue'
import AccountActivationCard from '@/components/auth/AccountActivationCard.vue'
import { useAuthStore } from '@/stores/auth'
import * as api from '@/api/auth'

const router = useRouter()
const auth = useAuthStore()

const loading = ref(false)
const error = ref('')
const activated = ref(false)

const login = computed(() => auth.pendingRegistration?.login ?? '')
const username = computed(() => auth.pendingRegistration?.username ?? '')
const devCode = computed(() => auth.pendingRegistration?.devCode ?? '')

async function handleSubmit(code: string) {
  if (loading.value) return

  const registrationId = auth.pendingRegistration?.registrationId

  if (!registrationId) {
    void router.replace({ name: 'register' })
    return
  }

  loading.value = true
  error.value = ''

  try {
    await api.activateRegistration(registrationId, code)

    auth.cancelRegistration()
    activated.value = true
  } catch (requestError) {
    error.value =
      requestError instanceof Error ? requestError.message : 'Не удалось активировать аккаунт.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="register-activation">
    <AuthCard v-if="activated" title="Активация аккаунта" title-shadow>
      <div class="register-activation__success">
        <p class="register-activation__text">
          Аккаунт активирован. Теперь вы можете войти, указав электронную почту
          <strong>{{ login }}</strong
          >.
        </p>

        <RouterLink class="register-activation__link" :to="{ name: 'login' }">
          Перейти ко входу
        </RouterLink>
      </div>
    </AuthCard>

    <AccountActivationCard
      v-else-if="login"
      :login="login"
      :username="username"
      :loading="loading"
      :error="error"
      :dev-code="devCode"
      @submit="handleSubmit"
    />
  </div>
</template>

<style scoped>
.register-activation {
  position: relative;
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
}

.register-activation__success {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 8px;
}

.register-activation__text {
  font: 400 14px/20px Roboto, sans-serif;
}

.register-activation__link {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 36px;
  padding: 8px;
  color: var(--color-text);
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  font: 400 14px/18px Roboto, sans-serif;
}

.register-activation__link:hover {
  background: var(--color-card-hover);
}
</style>
