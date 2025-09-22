import { ref } from 'vue'

// Toast 类型定义
export interface ToastMessage {
  id: string
  type: 'success' | 'error' | 'warning' | 'info'
  title: string
  message?: string
  duration?: number
  timestamp: number
}

// 全局 toast 状态
const toasts = ref<ToastMessage[]>([])
let toastIdCounter = 0

// 生成唯一 ID
const generateId = (): string => {
  return `toast-${++toastIdCounter}-${Date.now()}`
}

// 添加 toast
const addToast = (toast: Omit<ToastMessage, 'id' | 'timestamp'>): string => {
  const id = generateId()
  const newToast: ToastMessage = {
    ...toast,
    id,
    timestamp: Date.now(),
    duration: toast.duration ?? 4000
  }
  
  toasts.value.push(newToast)
  
  // 自动移除
  if (newToast.duration > 0) {
    setTimeout(() => {
      removeToast(id)
    }, newToast.duration)
  }
  
  return id
}

// 移除 toast
const removeToast = (id: string): void => {
  const index = toasts.value.findIndex(toast => toast.id === id)
  if (index > -1) {
    toasts.value.splice(index, 1)
  }
}

// 清空所有 toast
const clearToasts = (): void => {
  toasts.value = []
}

// Toast API
export const useToast = () => {
  const toast = {
    success: (title: string, message?: string, duration?: number) => {
      return addToast({ type: 'success', title, message, duration })
    },
    
    error: (title: string, message?: string, duration?: number) => {
      return addToast({ type: 'error', title, message, duration })
    },
    
    warning: (title: string, message?: string, duration?: number) => {
      return addToast({ type: 'warning', title, message, duration })
    },
    
    info: (title: string, message?: string, duration?: number) => {
      return addToast({ type: 'info', title, message, duration })
    },
    
    remove: removeToast,
    clear: clearToasts
  }
  
  const showToast = (message: string, type: 'success' | 'error' | 'warning' | 'info' = 'info') => {
    return toast[type](message)
  }

  return {
    toast,
    showToast,
    toasts: toasts.value,
    removeToast,
    clearToasts
  }
}

// 导出全局状态供组件使用
export { toasts }
