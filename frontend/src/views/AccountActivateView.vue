<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AuthCard from '@/components/auth/AuthCard.vue'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const token = computed(() => (typeof route.query.token === 'string' ? route.query.token : ''))
const username = computed(() =>
  typeof route.query.username === 'string' ? route.query.username : '',
)

function resolveActivation() {
  if (!token.value) return auth.account?.activated === true

  if (auth.activateAccount(token.value)) return true

  return (
    Boolean(username.value) &&
    auth.createAccountFromActivation({
      login: token.value,
      username: username.value,
    })
  )
}

const activated = ref(resolveActivation())

if (token.value && activated.value) {
  void router.replace({ name: 'home' })
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
              Настройки аккаунта
            </h2>

            <p class="activate__text">
              В настройках аккаунта можно включить двухфакторную аутентификацию и выбрать
              способ подтверждения.
            </p>

            <RouterLink
              class="activate__settings-link"
              :to="{ name: 'account-settings' }"
            >
              Перейти в настройки
            </RouterLink>
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

.activate__settings-link {
  align-self: flex-start;
  color: var(--color-auth-link);
  font: 500 14px/18px Roboto, sans-serif;
  text-decoration: underline;
  text-underline-offset: 2px;
}

.activate__settings-link:hover {
  opacity: 0.8;
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
