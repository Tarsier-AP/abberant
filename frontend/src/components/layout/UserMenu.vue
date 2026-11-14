<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import UserAvatar from '@/components/ui/UserAvatar.vue'
import { useAuthStore } from '@/stores/auth'

withDefaults(
  defineProps<{
    /** dropdown — floating panel in the header; stack — full-width block in the mobile menu. */
    variant?: 'dropdown' | 'stack'
  }>(),
  { variant: 'dropdown' },
)

const emit = defineEmits<{ select: [] }>()

const auth = useAuthStore()
const router = useRouter()

const items = computed(() => {
  const base = [
    { label: 'Закладки', to: { name: 'bookmarks' } },
    { label: 'Активность', to: { name: 'activity' } },
    { label: 'Настройки аккаунта', to: { name: 'account-settings' } },
  ]

  return auth.canModerate
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
  <div class="user-menu" :class="`user-menu--${variant}`">
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
      <span>{{ item.label }}</span>
      <svg
        v-if="variant === 'stack'"
        width="8"
        height="16"
        viewBox="0 0 8 16"
        fill="none"
        aria-hidden="true"
      >
        <path
          d="m1 1 6 7-6 7"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
    </RouterLink>

    <button
      class="user-menu__item user-menu__item--logout"
      type="button"
      @click="logout"
    >
      <span>Выйти</span>
    </button>
  </div>
</template>

<style scoped>
.user-menu {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 8px;
}

.user-menu--dropdown {
  min-width: 224px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.user-menu--stack {
  gap: 8px;
  padding: 0;
}

.user-menu__head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  margin-bottom: 4px;
  border-bottom: 1px solid var(--color-divider);
}

.user-menu--stack .user-menu__head {
  margin-bottom: 0;
  padding: 8px 12px;
  background: var(--color-card);
  border-bottom: 0;
  border-radius: var(--radius);
  box-shadow: var(--shadow);
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
  gap: 8px;
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

.user-menu__item span {
  flex: 1;
}

.user-menu__item svg {
  flex: none;
  color: var(--color-icon);
}

.user-menu__item:hover {
  background: var(--color-card-hover);
}

.user-menu--stack .user-menu__item {
  min-height: 35px;
  background: var(--color-card);
  box-shadow: var(--shadow);
  font: 400 14px/17px Roboto, sans-serif;
}

.user-menu--stack .user-menu__item:hover {
  background: var(--color-card-hover);
}

.user-menu__item--logout {
  color: var(--color-danger);
}
</style>
