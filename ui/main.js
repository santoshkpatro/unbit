import '@/assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createNotivue } from 'notivue'

import App from '@/App.vue'
import router from '@/router'
import { settingMetaAPI } from '@/api/setting'
import { authStatusAPI } from '@/api/auth'

import { useSettingStore } from '@/stores/setting'
import { useAuthStore } from '@/stores/auth'

const app = createApp(App)

app.use(createPinia())
app.use(
  createNotivue({
    position: 'top-right',
    enqueue: true,
    max: 5,
  }),
)

async function bootstrap() {
  const metaData = await settingMetaAPI()
  const authStatus = await authStatusAPI()

  const settingStore = useSettingStore()
  settingStore.setSetting(metaData)

  const authStore = useAuthStore()
  if (authStatus.isLoggedIn) {
    authStore.setLoggedInUser(authStatus.userProfile)
  }

  app.use(router)
  app.mount('#app')
}

bootstrap()
