// Toast桥接器，为errorHandler提供统一的toast接口
// 由于errorHandler是单例类，不能直接使用Vue composables
// 所以需要这个桥接器来连接两个系统

interface ToastBridge {
  success: (message: string) => void
  error: (message: string) => void
  warning: (message: string) => void
  info: (message: string) => void
}

// 全局toast方法，会被App.vue设置
let globalToastMethods: ToastBridge | null = null

export const setGlobalToastMethods = (methods: ToastBridge) => {
  globalToastMethods = methods
}

export const toastBridge: ToastBridge = {
  success: (message: string) => {
    if (globalToastMethods) {
      globalToastMethods.success(message)
    } else {
      console.warn('Toast bridge not initialized, message:', message)
    }
  },
  error: (message: string) => {
    if (globalToastMethods) {
      globalToastMethods.error(message)
    } else {
      console.error('Toast bridge not initialized, error:', message)
    }
  },
  warning: (message: string) => {
    if (globalToastMethods) {
      globalToastMethods.warning(message)
    } else {
      console.warn('Toast bridge not initialized, warning:', message)
    }
  },
  info: (message: string) => {
    if (globalToastMethods) {
      globalToastMethods.info(message)
    } else {
      console.info('Toast bridge not initialized, info:', message)
    }
  }
}