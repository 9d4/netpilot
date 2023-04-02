import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'

import './assets/base.css'
import { useBoardStore } from './stores/boards'

const app = createApp(App)

let wsBaseUrl = import.meta.env.VITE_WS_BASE_URL
if (wsBaseUrl === undefined) {
  const protocol = location.protocol == 'https:' ? 'wss:' : 'ws:'
  wsBaseUrl = `${protocol}//ws`
}

app.provide('wsBaseUrl', wsBaseUrl)

app.use(createPinia())
app.use(router)

useBoardStore()

app.mount('#app')
