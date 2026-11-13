import { createApp } from 'vue'
import { createPinia } from 'pinia'

import './assets/main.css'
import './assets/article.css'
import App from './App.vue'
import router from './router'
import { useAuthStore } from '@/stores/auth'
import { useArticlesStore } from '@/stores/articles'
import { migrateLocalArticles } from '@/utils/migrateLocalArticles'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// Restore the session account (server-side session cookie) before mounting,
// then load the articles from the API.
useAuthStore(pinia)
  .init()
  .then(() => {
    void migrateLocalArticles()

    return useArticlesStore(pinia).init()
  })
  .finally(() => {
    app.mount('#app')
  })
