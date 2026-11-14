<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import UserAvatar from '@/components/ui/UserAvatar.vue'
import { useAuthStore } from '@/stores/auth'

const emit = defineEmits<{ select: [] }>()

const auth = useAuthStore()
const router = useRouter()

const items = computed(() => {
  const base = [
    { label: 'Закладки', to: { name: 'bookmarks' } },
    { label: 'Активность', to: { name: 'activity' } },
    { label: 'Настройки аккаунта', to: { name: 'account-settings' } },
  ]

  return auth.canEditArticles
    ? [...base, { label: 'Редактор', to: { name: 'editor-articles' } }]
    : base
})

async function logout() {
  await auth.endSession()
  router.push({ name: 'home' })
  emit('select')
}
</script>

<template>
  <div class="user-menu">
    <div class="user-menu__head">
      <UserAvatar
        class="user-menu__avatar"
        :name="auth.account?.username ?? ''"
        :avatar="auth.account?.avatar"
        :size="36"
      />
      <span class="user-menu__name">{{ auth.account?.username }}</span>
    </div>

    <RouterLink
      v-for="item in items"
      :key="item.label"
      class="user-menu__item"
      :to="item.to"
      @click="emit('select')"
    >
      {{ item.label }}
    </RouterLink>

    <button
      class="user-menu__item user-menu__item--logout"
      type="button"
      @click="logout"
    >
      Выйти
    </button>
  </div>
</template>

<style scoped>
.user-menu {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 224px;
  padding: 8px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.user-menu__head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  margin-bottom: 4px;
  border-bottom: 1px solid var(--color-divider);
}

.user-menu__name {
  overflow: hidden;
  font: 500 14px/17px Roboto, sans-serif;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-menu__item {
  display: flex;
  align-items: center;
  min-height: 36px;
  padding: 8px 12px;
  color: inherit;
  background: transparent;
  border: 0;
  border-radius: var(--radius);
  font: 400 13px/16px Roboto, sans-serif;
  text-align: left;
  cursor: pointer;
}

.user-menu__item:hover {
  background: var(--color-card-hover);
}

.user-menu__item--logout {
  color: var(--color-danger);
}
</style>
