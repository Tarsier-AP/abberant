import { createRouter, createWebHistory } from 'vue-router'

import HomeView from '@/views/HomeView.vue'
import NewsView from '@/views/NewsView.vue'
import WalkthroughsView from '@/views/WalkthroughsView.vue'
import ReviewsView from '@/views/ReviewsView.vue'
import HelpView from '@/views/HelpView.vue'
import LoginView from '@/views/LoginView.vue'
import RegisterView from '@/views/RegisterView.vue'
import RegisterPasswordView from '@/views/RegisterPasswordView.vue'
import RegisterActivationView from '@/views/RegisterActivationView.vue'
import AccountActivateView from '@/views/AccountActivateView.vue'
import TwoFactorView from '@/views/TwoFactorView.vue'
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
      path: '/2fa',
      name: 'two-factor',
      component: TwoFactorView,
      meta: { title: 'Двухфакторная аутентификация' },
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView,
      meta: { title: 'Регистрация' },
    },
    {
      path: '/register/password',
      name: 'register-password',
      component: RegisterPasswordView,
      meta: { title: 'Создание пароля' },
    },
    {
      path: '/register/activate',
      name: 'register-activation',
      component: RegisterActivationView,
      meta: { title: 'Активация аккаунта' },
    },
    {
      path: '/activate',
      name: 'activate',
      component: AccountActivateView,
      meta: { title: 'Активация аккаунта' },
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
