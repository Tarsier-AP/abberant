<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AuthCard from '@/components/auth/AuthCard.vue'
import { useAuthStore } from '@/stores/auth'
import type { TwoFactorMethod } from '@/types/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const token = computed(() => (typeof route.query.token === 'string' ? route.query.token : ''))
const activated = ref(
  token.value ? auth.activateAccount(token.value) : auth.account?.activated === true,
)

if (token.value && activated.value) {
  void router.replace({ name: 'home' })
}

function handleTwoFactorChange(event: Event) {
  auth.setTwoFactor((event.target as HTMLInputElement).checked)
}

function handleMethodChange(event: Event) {
  auth.setTwoFactorMethod((event.target as HTMLSelectElement).value as TwoFactorMethod)
}
</script>

<template>
  <div class="activate">
    <AuthCard title="Активация аккаунта" title-shadow>
      <div class="activate__content">
        <template v-if="auth.account && activated">
          <p class="activate__text">
            Аккаунт активирован. Теперь вы можете войти под логином
            <strong>{{ auth.account.login }}</strong
            >.
          </p>

          <section class="activate__settings" aria-labelledby="activate-settings-title">
            <h2 id="activate-settings-title" class="activate__settings-title">
              Настройки аккаунта (демо)
            </h2>

            <label class="activate__checkbox">
              <input
                type="checkbox"
                :checked="auth.account.twoFactorEnabled"
                @change="handleTwoFactorChange"
              />
              Двухфакторная аутентификация
            </label>

            <label v-if="auth.account.twoFactorEnabled" class="activate__method">
              <span>Способ подтверждения</span>
              <select :value="auth.account.twoFactorMethod" @change="handleMethodChange">
                <option value="sms">SMS</option>
                <option value="email">Электронная почта</option>
                <option value="authenticator">Приложение-аутентификатор</option>
              </select>
            </label>
          </section>

          <RouterLink class="activate__submit" :to="{ name: 'login' }">
            Перейти ко входу
          </RouterLink>
        </template>

        <template v-else>
          <p class="activate__text">
            Ссылка активации недействительна или устарела. Зарегистрируйтесь заново.
          </p>
          <RouterLink class="activate__submit" :to="{ name: 'register' }">
            К регистрации
          </RouterLink>
        </template>
      </div>
    </AuthCard>
  </div>
</template>

<style scoped>
.activate {
  position: relative;
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
}

.activate__content {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 8px;
}

.activate__text {
  font: 400 14px/20px Roboto, sans-serif;
}

.activate__settings {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
  background: var(--color-surface);
  border-radius: var(--radius);
}

.activate__settings-title {
  font: 700 14px/18px Roboto, sans-serif;
}

.activate__checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  font: 400 14px/18px Roboto, sans-serif;
  cursor: pointer;
}

.activate__checkbox input {
  width: 16px;
  height: 16px;
  margin: 0;
  accent-color: var(--color-accent);
  cursor: pointer;
}

.activate__method {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font: 400 13px/16px Roboto, sans-serif;
}

.activate__method select {
  min-height: 30px;
  padding: 4px 8px;
  color: var(--color-text);
  background: var(--color-auth-input);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius);
  font: 400 14px/20px Roboto, sans-serif;
}

.activate__submit {
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

.activate__submit:hover {
  background: var(--color-card-hover);
}
</style>
