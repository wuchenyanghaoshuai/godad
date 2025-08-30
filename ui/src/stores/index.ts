// 状态管理统一导出
import { createPinia } from 'pinia'

// 创建Pinia实例
export const pinia = createPinia()

// 导出所有store
export { useAuthStore } from './auth'

// 默认导出pinia实例
export default pinia