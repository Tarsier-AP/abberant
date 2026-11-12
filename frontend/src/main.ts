import { createApp } from 'vue'
import { createPinia } from 'pinia'

import './assets/main.css'
import './assets/article.css'
import App from './App.vue'
import router from './router'
import { useAuthStore } from '@/stores/auth'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// Restore the session account (server-side session cookie) before mounting.
useAuthStore(pinia)
  .init()
  .finally(() => {
    app.mount('#app')
  })
