import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import { pinia } from './stores'
import { initApiInterceptor } from './utils/apiInterceptor'
import { errorHandler } from './utils/errorHandler'

const app = createApp(App)

// 初始化API拦截器
initApiInterceptor()

// 使用Pinia状态管理
app.use(pinia)

// 使用Vue Router
app.use(router)

// 全局错误处理
app.config.errorHandler = (error, instance, info) => {
  errorHandler.handleGlobalError(new ErrorEvent('error', {
    error: error as Error,
    message: (error as Error)?.message || 'Unknown error',
    filename: info,
    lineno: 0,
    colno: 0
  }))
}

// 挂载应用
app.mount('#app')
