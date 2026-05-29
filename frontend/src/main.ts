import './styles/global.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { vTabIndent } from './directives/tabIndent'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.directive('tab-indent', vTabIndent)

app.mount('#app')
