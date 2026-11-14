import { createRouter, createWebHistory } from 'vue-router'

import PlaceholderView from '@/views/PlaceholderView.vue'
import { useAuthStore } from '@/stores/auth'

/**
 * Views live in their own chunks and are fetched on the first navigation to
 * their route, keeping the initial bundle small.
 */
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', name: 'home', component: () => import('@/views/HomeView.vue') },
    { path: '/news', component: () => import('@/views/NewsView.vue'), meta: { title: 'Новости' } },
    { path: '/news/:platform', component: () => import('@/views/NewsView.vue') },
    { path: '/reviews', component: () => import('@/views/ReviewsView.vue'), meta: { title: 'Обзоры' } },
    { path: '/reviews/:section', component: () => import('@/views/ReviewsView.vue') },
    { path: '/walkthroughs', component: () => import('@/views/WalkthroughsView.vue') },
    { path: '/walkthroughs/:section', component: () => import('@/views/WalkthroughsView.vue') },
    {
      path: '/help',
      component: () => import('@/views/HelpView.vue'),
      meta: { title: 'Помощь по сайту' },
    },
    {
      path: '/bookmarks',
      name: 'bookmarks',
      component: () => import('@/views/BookmarksView.vue'),
      meta: { title: 'Закладки', requiresAuth: true },
    },
    {
      path: '/activity',
      name: 'activity',
      component: () => import('@/views/ActivityView.vue'),
      meta: { title: 'История активности', requiresAuth: true },
    },
    { path: '/login', name: 'login', component: () => import('@/views/LoginView.vue'), meta: { title: 'Вход' } },
    {
      path: '/2fa',
      name: 'two-factor',
      component: () => import('@/views/TwoFactorView.vue'),
      meta: { title: 'Двухфакторная аутентификация' },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/RegisterView.vue'),
      meta: { title: 'Регистрация' },
    },
    {
      path: '/register/password',
      name: 'register-password',
      component: () => import('@/views/RegisterPasswordView.vue'),
      meta: { title: 'Создание пароля' },
    },
    {
      path: '/register/activate',
      name: 'register-activation',
      component: () => import('@/views/RegisterActivationView.vue'),
      meta: { title: 'Активация аккаунта' },
    },
    {
      path: '/activate',
      name: 'activate',
      component: () => import('@/views/AccountActivateView.vue'),
      meta: { title: 'Активация аккаунта' },
    },
    {
      path: '/account/settings',
      name: 'account-settings',
      component: () => import('@/views/AccountSettingsView.vue'),
      meta: { title: 'Настройки аккаунта' },
    },
    {
      path: '/password-reset',
      component: PlaceholderView,
      meta: { title: 'Восстановление пароля' },
    },
    { path: '/about', component: PlaceholderView, meta: { title: 'О нас' } },
    { path: '/career', component: PlaceholderView, meta: { title: 'Карьера в Abberant' } },
    {
      path: '/privacy',
      component: PlaceholderView,
      meta: { title: 'Политика конфиденциальности' },
    },
    {
      path: '/terms',
      component: PlaceholderView,
      meta: { title: 'Условия использования сайта' },
    },
    {
      path: '/personal-data',
      component: PlaceholderView,
      meta: { title: 'Политика обработки персональных данных' },
    },
    { path: '/article/:id', name: 'article', component: () => import('@/views/ArticleView.vue'), meta: { title: 'Статья' } },
    {
      path: '/editor/sections',
      name: 'editor-sections',
      component: () => import('@/views/EditorSectionsView.vue'),
      meta: { title: 'Секции', requiresEditorRole: true },
    },
    {
      path: '/editor/articles',
      name: 'editor-articles',
      component: () => import('@/views/EditorArticlesView.vue'),
      meta: { title: 'Статьи', requiresEditorRole: true },
    },
    {
      path: '/editor',
      name: 'editor',
      component: () => import('@/views/EditorView.vue'),
      meta: { title: 'Новая статья', requiresEditorRole: true },
    },
    {
      path: '/editor/:id',
      name: 'editor-edit',
      component: () => import('@/views/EditorView.vue'),
      meta: { title: 'Редактирование статьи', requiresEditorRole: true },
    },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
  scrollBehavior(to, _from, savedPosition) {
    if (savedPosition) return savedPosition
    if (to.hash) return { el: to.hash, top: 16 }

    return { top: 0 }
  },
})

/**
 * Authenticated pages send guests to the login form; the session is restored
 * first on hard navigation straight to such a route.
 */
router.beforeEach(async (to) => {
  if (!to.meta.requiresAuth) return true

  const auth = useAuthStore()

  if (!auth.initialized) await auth.init()

  if (!auth.account) return { name: 'login' }

  return true
})

/** Only editors, owners and moderators may open the article editor. */
router.beforeEach(async (to) => {
  if (!to.meta.requiresEditorRole) return true

  const auth = useAuthStore()

  // The app mounts before the session is restored, so load it here on hard
  // navigation straight to an editor route.
  if (!auth.initialized) await auth.init()

  if (!auth.account) return { name: 'login' }
  if (!auth.canModerate) return { path: '/' }

  return true
})

export default router
