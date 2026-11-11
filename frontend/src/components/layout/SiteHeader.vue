<script setup lang="ts">
import { onMounted, ref } from 'vue'
import ThemeIcon from '@/components/ui/ThemeIcon.vue'
import { useAuthStore } from '@/stores/auth'

type Theme = 'light' | 'dark'

const THEME_KEY = 'abberant-theme'

const menuOpen = ref(false)
const auth = useAuthStore()

function readStoredTheme(): Theme | null {
  try {
    const stored = localStorage.getItem(THEME_KEY)
    return stored === 'light' || stored === 'dark' ? stored : null
  } catch {
    return null
  }
}

function getSystemTheme(): Theme {
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

const theme = ref<Theme>(readStoredTheme() ?? getSystemTheme())

function applyTheme(value: Theme) {
  theme.value = value
  document.documentElement.dataset.theme = value

  try {
    localStorage.setItem(THEME_KEY, value)
  } catch {
    // Ignore storage errors (e.g. private browsing).
  }
}

function toggleTheme() {
  applyTheme(theme.value === 'dark' ? 'light' : 'dark')
}

onMounted(() => {
  document.documentElement.dataset.theme = theme.value
})

const links = [
  { label: 'Новости', to: '/news' },
  { label: 'Обзоры', to: '/reviews' },
  { label: 'Прохождения', to: '/walkthroughs' },
  { label: 'Помощь', to: '/help' },
]
</script>

<template>
  <header class="header" :class="{ 'header--open': menuOpen }">
    <div class="header__inner container">
      <RouterLink class="logo" to="/" aria-label="Abberant — на главную">
        <span class="logo__text">Abberant</span>
      </RouterLink>

      <nav class="nav" aria-label="Основная навигация">
        <RouterLink v-for="link in links" :key="link.to" class="nav__link" :to="link.to">
          {{ link.label }}
        </RouterLink>
      </nav>

      <div class="actions">
        <button class="icon-btn search" type="button" aria-label="Поиск">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <circle cx="10.5" cy="10.5" r="6.5" stroke="currentColor" stroke-width="2" />
            <path
              d="m15.5 15.5 4.5 4.5"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
            />
          </svg>
        </button>

        <button
          class="icon-btn theme-toggle"
          type="button"
          :aria-label="theme === 'dark' ? 'Включить светлую тему' : 'Включить тёмную тему'"
          :title="theme === 'dark' ? 'Включить светлую тему' : 'Включить тёмную тему'"
          @click="toggleTheme"
        >
          <ThemeIcon />
        </button>

        <RouterLink
          class="icon-btn account"
          :to="{ name: auth.account?.activated ? 'account-settings' : 'login' }"
          :aria-label="auth.account?.activated ? 'Аккаунт' : 'Войти'"
          :title="auth.account?.activated ? 'Аккаунт' : 'Войти'"
        >
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" aria-hidden="true">
            <circle cx="12" cy="8" r="4" stroke="currentColor" stroke-width="2" />
            <path
              d="M4 20c0-4.418 3.582-8 8-8s8 3.582 8 8"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
            />
          </svg>
        </RouterLink>

        <button
          class="burger"
          type="button"
          :aria-expanded="menuOpen"
          aria-controls="mobile-menu"
          :aria-label="menuOpen ? 'Закрыть меню' : 'Открыть меню'"
          @click="menuOpen = !menuOpen"
        >
          <span></span>
          <span></span>
          <span></span>
        </button>
      </div>
    </div>

    <nav
      v-show="menuOpen"
      id="mobile-menu"
      class="mobile-menu"
      aria-label="Мобильная навигация"
    >
      <div class="mobile-menu__inner">
        <RouterLink
          v-for="link in links"
          :key="link.to"
          class="mobile-menu__link"
          :to="link.to"
          @click="menuOpen = false"
        >
          <span>{{ link.label }}</span>
          <svg width="8" height="16" viewBox="0 0 8 16" fill="none" aria-hidden="true">
            <path
              d="m1 1 6 7-6 7"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </RouterLink>

        <div class="mobile-menu__actions">
          <button
            class="mobile-menu__action"
            type="button"
            :aria-label="theme === 'dark' ? 'Включить светлую тему' : 'Включить тёмную тему'"
            :title="theme === 'dark' ? 'Включить светлую тему' : 'Включить тёмную тему'"
            @click="toggleTheme"
          >
            <ThemeIcon />
          </button>

          <RouterLink
            class="mobile-menu__action"
            :to="{ name: auth.account?.activated ? 'account-settings' : 'login' }"
            @click="menuOpen = false"
          >
            {{ auth.account?.activated ? 'Аккаунт' : 'Войти' }}
          </RouterLink>
        </div>
      </div>
    </nav>
  </header>
</template>

<style scoped>
.header {
  position: relative;
  z-index: 20;
  width: 100%;
  color: var(--color-text);
  background: var(--color-header);
}

.header__inner {
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: var(--header-height);
  padding-block: 8px;
}

.logo {
  display: flex;
  flex: none;
  align-items: center;
  gap: 4px;
}

.logo__text {
  font: 600 20px/24px Inter, sans-serif;
}

.nav {
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  gap: 16px;
}

.nav__link {
  font: 400 14px/17px Roboto, sans-serif;
  transition: opacity 0.15s ease;
}

.nav__link:hover {
  opacity: 0.6;
}

.actions {
  display: flex;
  flex: none;
  align-items: center;
  gap: 8px;
  margin-left: auto;
}

.icon-btn {
  display: grid;
  flex: none;
  place-items: center;
  width: 40px;
  height: 40px;
  padding: 0;
  color: inherit;
  background: transparent;
  border: 0;
  border-radius: var(--radius);
  cursor: pointer;
}

.icon-btn svg {
  width: 24px;
  height: 24px;
}

.theme-toggle svg {
  width: 22px;
  height: 22px;
}

.icon-btn:hover {
  background: var(--color-card);
}

.theme-toggle {
  background: var(--color-card);
  box-shadow: var(--shadow);
}

.theme-toggle:hover {
  background: var(--color-card-hover);
}

.account {
  background: var(--color-card);
  box-shadow: var(--shadow);
}

.account:hover {
  background: var(--color-card-hover);
}

.burger {
  display: none;
  flex-direction: column;
  justify-content: space-between;
  width: 40px;
  height: 40px;
  padding: 9px 8px;
  background: transparent;
  border: 0;
  border-radius: var(--radius);
  cursor: pointer;
}

.burger span {
  display: block;
  width: 100%;
  height: 3px;
  background: var(--color-burger);
  border-radius: 8px;
}

.mobile-menu {
  padding: 8px 16px;
  background: var(--color-header);
  border-radius: 0 0 12px 12px;
}

.mobile-menu__inner {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
}

.mobile-menu__link {
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 35px;
  padding: 8px 12px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  font: 400 14px/17px Roboto, sans-serif;
}

.mobile-menu__link span {
  flex: 1;
}

.mobile-menu__link svg {
  flex: none;
  color: var(--color-icon);
}

.mobile-menu__link:hover {
  background: var(--color-card-hover);
}

.mobile-menu__actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.mobile-menu__action {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  min-height: 36px;
  padding: 8px;
  background: var(--color-card);
  border: 0;
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  font: 400 13px/16px Roboto, sans-serif;
  cursor: pointer;
}

.mobile-menu__action svg {
  width: 22px;
  height: 22px;
}

.mobile-menu__action:hover {
  background: var(--color-card-hover);
}

@media (min-width: 1024px) {
  .mobile-menu {
    display: none !important;
  }
}

@media (max-width: 1023px) {
  .header--open {
    border-radius: 0 0 12px 12px;
  }

  .nav,
  .search,
  .theme-toggle,
  .account {
    display: none;
  }

  .burger {
    display: flex;
  }
}
</style>
