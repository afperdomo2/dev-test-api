import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { VueQueryPlugin } from '@tanstack/vue-query'

import App from './App.vue'
import router from './router'
import vuetify from './plugins/vuetify'
import { queryPluginOptions } from './plugins/query'

async function bootstrap() {
  const app = createApp(App)

  app.use(createPinia())

  const { useAuthStore } = await import('@/stores/auth.store')
  const authStore = useAuthStore()
  await authStore.initSession()
  await authStore.checkStatus()

  app.use(router)
  app.use(vuetify)
  app.use(VueQueryPlugin, queryPluginOptions)

  app.mount('#app')
}

bootstrap()
