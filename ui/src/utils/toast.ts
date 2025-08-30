// Vue兼容的Toast通知系统

export interface ToastOptions {
  duration?: number
  type?: 'success' | 'error' | 'warning' | 'info'
  position?: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left' | 'top-center' | 'bottom-center'
}

export interface ToastAction {
  label: string
  onClick: () => void
}

class ToastManager {
  private container: HTMLElement | null = null
  private toastId = 0

  constructor() {
    this.createContainer()
  }

  private createContainer() {
    if (typeof window === 'undefined') return
    
    this.container = document.createElement('div')
    this.container.id = 'toast-container'
    this.container.style.cssText = `
      position: fixed;
      top: 20px;
      right: 20px;
      z-index: 9999;
      pointer-events: none;
    `
    document.body.appendChild(this.container)
  }

  private createToast(
    message: string, 
    options: ToastOptions = {}, 
    action?: ToastAction
  ): HTMLElement {
    const { 
      duration = 4000, 
      type = 'info', 
      position = 'top-right' 
    } = options

    const toast = document.createElement('div')
    const toastId = ++this.toastId
    toast.id = `toast-${toastId}`
    
    // 设置样式
    const typeColors = {
      success: '#10b981',
      error: '#ef4444',
      warning: '#f59e0b',
      info: '#3b82f6'
    }
    
    toast.style.cssText = `
      background: white;
      border-left: 4px solid ${typeColors[type]};
      border-radius: 8px;
      box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
      padding: 16px;
      margin-bottom: 12px;
      max-width: 400px;
      pointer-events: auto;
      transform: translateX(100%);
      transition: all 0.3s ease;
      display: flex;
      align-items: center;
      justify-content: space-between;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      font-size: 14px;
      line-height: 1.4;
    `

    // 创建消息内容
    const messageEl = document.createElement('div')
    messageEl.textContent = message
    messageEl.style.cssText = `
      flex: 1;
      color: #374151;
    `
    toast.appendChild(messageEl)

    // 添加操作按钮
    if (action) {
      const actionBtn = document.createElement('button')
      actionBtn.textContent = action.label
      actionBtn.style.cssText = `
        background: ${typeColors[type]};
        color: white;
        border: none;
        border-radius: 4px;
        padding: 6px 12px;
        margin-left: 12px;
        cursor: pointer;
        font-size: 12px;
        font-weight: 500;
      `
      actionBtn.onclick = () => {
        action.onClick()
        this.removeToast(toast)
      }
      toast.appendChild(actionBtn)
    }

    // 添加关闭按钮
    const closeBtn = document.createElement('button')
    closeBtn.innerHTML = '×'
    closeBtn.style.cssText = `
      background: none;
      border: none;
      color: #9ca3af;
      cursor: pointer;
      font-size: 18px;
      font-weight: bold;
      margin-left: 8px;
      padding: 0;
      width: 20px;
      height: 20px;
      display: flex;
      align-items: center;
      justify-content: center;
    `
    closeBtn.onclick = () => this.removeToast(toast)
    toast.appendChild(closeBtn)

    return toast
  }

  private removeToast(toast: HTMLElement) {
    toast.style.transform = 'translateX(100%)'
    toast.style.opacity = '0'
    setTimeout(() => {
      if (toast.parentNode) {
        toast.parentNode.removeChild(toast)
      }
    }, 300)
  }

  private showToast(
    message: string, 
    options: ToastOptions = {}, 
    action?: ToastAction
  ) {
    if (!this.container) return

    const toast = this.createToast(message, options, action)
    this.container.appendChild(toast)

    // 触发动画
    setTimeout(() => {
      toast.style.transform = 'translateX(0)'
    }, 10)

    // 自动移除
    if (options.duration !== 0) {
      setTimeout(() => {
        this.removeToast(toast)
      }, options.duration || 4000)
    }
  }

  success(message: string, options?: ToastOptions) {
    this.showToast(message, { ...options, type: 'success' })
  }

  error(message: string, options?: ToastOptions & { action?: ToastAction }) {
    const { action, ...toastOptions } = options || {}
    this.showToast(message, { ...toastOptions, type: 'error' }, action)
  }

  warning(message: string, options?: ToastOptions) {
    this.showToast(message, { ...options, type: 'warning' })
  }

  info(message: string, options?: ToastOptions) {
    this.showToast(message, { ...options, type: 'info' })
  }
}

// 创建全局实例
const toastManager = new ToastManager()

// 导出便捷方法
export const toast = {
  success: (message: string, options?: ToastOptions) => toastManager.success(message, options),
  error: (message: string, options?: ToastOptions & { action?: ToastAction }) => toastManager.error(message, options),
  warning: (message: string, options?: ToastOptions) => toastManager.warning(message, options),
  info: (message: string, options?: ToastOptions) => toastManager.info(message, options)
}

export default toast