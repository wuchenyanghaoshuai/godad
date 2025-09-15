import { ref } from 'vue'

// 全局通知状态管理
const notificationCount = ref(0)

// 事件类型
type NotificationEvent = 'refresh' | 'mark-read' | 'clear-all'

// 事件监听器
const listeners = new Map<NotificationEvent, Set<() => void>>()

export function useNotificationSync() {
  // 触发通知刷新事件
  const triggerRefresh = () => {
    const refreshListeners = listeners.get('refresh')
    if (refreshListeners) {
      refreshListeners.forEach(listener => listener())
    }
  }

  // 监听通知事件
  const onNotificationEvent = (event: NotificationEvent, callback: () => void) => {
    if (!listeners.has(event)) {
      listeners.set(event, new Set())
    }
    listeners.get(event)!.add(callback)

    // 返回取消监听的函数
    return () => {
      listeners.get(event)?.delete(callback)
    }
  }

  // 更新通知数量
  const updateNotificationCount = (count: number) => {
    notificationCount.value = count
  }

  return {
    notificationCount,
    triggerRefresh,
    onNotificationEvent,
    updateNotificationCount
  }
}