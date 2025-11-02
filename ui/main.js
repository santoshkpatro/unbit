import '@/assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createNotivue } from 'notivue'

import App from '@/App.vue'
import router from '@/router'
import { settingMetaAPI } from '@/api/setting'
import { useSettingStore } from '@/stores/setting'

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

  const settingStore = useSettingStore()
  settingStore.setSetting(metaData)

  app.use(router)
  app.mount('#app')
}

bootstrap()
