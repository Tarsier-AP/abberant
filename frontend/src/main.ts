import { createApp } from 'vue'
import { createPinia } from 'pinia'

// Self-hosted Roboto (was Google Fonts): same weights, no third-party request.
import '@fontsource/roboto/cyrillic-300.css'
import '@fontsource/roboto/cyrillic-400.css'
import '@fontsource/roboto/cyrillic-500.css'
import '@fontsource/roboto/cyrillic-700.css'
import '@fontsource/roboto/cyrillic-800.css'
import '@fontsource/roboto/300.css'
import '@fontsource/roboto/400.css'
import '@fontsource/roboto/500.css'
import '@fontsource/roboto/700.css'
import '@fontsource/roboto/800.css'

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

// Restore the session account and load the articles in parallel (the cookie is
// sent with both requests), then mount once with real content so the first
// paint is the page itself. The legacy migration needs the account role and
// runs right after the session is restored.
const auth = useAuthStore(pinia)
const articles = useArticlesStore(pinia)

Promise.all([auth.init(), articles.init()])
  .then(() => {
    void migrateLocalArticles()
  })
  .finally(() => {
    app.mount('#app')
  })
