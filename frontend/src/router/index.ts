import { createRouter, createWebHistory } from 'vue-router'

import HomeView from '@/views/HomeView.vue'
import NewsView from '@/views/NewsView.vue'
import WalkthroughsView from '@/views/WalkthroughsView.vue'
import ReviewsView from '@/views/ReviewsView.vue'
import HelpView from '@/views/HelpView.vue'
import LoginView from '@/views/LoginView.vue'
import PlaceholderView from '@/views/PlaceholderView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    { path: '/news', component: NewsView, meta: { title: 'Новости' } },
    { path: '/news/:platform', component: NewsView },
    { path: '/reviews', component: ReviewsView, meta: { title: 'Обзоры' } },
    { path: '/reviews/:section', component: ReviewsView },
    { path: '/walkthroughs', component: WalkthroughsView },
    { path: '/walkthroughs/:section', component: WalkthroughsView },
    { path: '/help', component: HelpView, meta: { title: 'Помощь по сайту' } },
    { path: '/login', name: 'login', component: LoginView, meta: { title: 'Вход' } },
    {
      path: '/login/2fa',
      name: 'two-factor',
      component: PlaceholderView,
      meta: { title: 'Двухфакторная аутентификация' },
    },
    {
      path: '/register',
      component: PlaceholderView,
      meta: { title: 'Регистрация' },
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
    { path: '/article/:id', component: PlaceholderView, meta: { title: 'Статья' } },
    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
})

export default router
