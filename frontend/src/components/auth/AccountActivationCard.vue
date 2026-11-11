<script setup lang="ts">
import { computed } from 'vue'
import AuthCard from './AuthCard.vue'

const props = defineProps<{
  login: string
  username: string
}>()

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

const isEmail = computed(() => EMAIL_PATTERN.test(props.login))

const destination = computed(() =>
  isEmail.value
    ? `на вашу электронную почту ${props.login}`
    : `в SMS на номер ${props.login}`,
)

const activationUrl = computed(() => {
  const origin = typeof window === 'undefined' ? '' : window.location.origin

  return `${origin}/activate?token=${encodeURIComponent(props.login)}&username=${encodeURIComponent(props.username)}`
})
</script>

<template>
  <AuthCard title="Активация аккаунта" title-shadow>
    <div class="account-activation">
      <p class="account-activation__text">
        Мы отправили ссылку для активации аккаунта {{ destination }}.
      </p>
      <p class="account-activation__text">
        Перейдите по ссылке из сообщения, чтобы активировать аккаунт.
      </p>

      <RouterLink
        class="account-activation__url"
        :to="{ name: 'activate', query: { token: login, username } }"
      >
        {{ activationUrl }}
      </RouterLink>

      <p class="account-activation__warning" role="note">
        Внимание: эта ссылка временная и показана только для демонстрации. В рабочей версии
        она придёт в сообщении на указанный контакт.
      </p>

      <p class="account-activation__note">Ссылка действует 24 часа.</p>
    </div>
  </AuthCard>
</template>

<style scoped>
.account-activation {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 8px;
}

.account-activation__text {
  font: 400 14px/20px Roboto, sans-serif;
}

.account-activation__url {
  color: var(--color-auth-link);
  font: 400 13px/18px Roboto, sans-serif;
  overflow-wrap: anywhere;
  text-decoration: underline;
}

.account-activation__url:hover {
  opacity: 0.8;
}

.account-activation__warning {
  padding: 8px;
  color: var(--color-warning);
  background: var(--color-surface);
  border-radius: var(--radius);
  font: 600 12px/16px Roboto, sans-serif;
}

.account-activation__note {
  color: var(--color-muted);
  font: 400 12px/16px Roboto, sans-serif;
  opacity: 0.7;
}
</style>
