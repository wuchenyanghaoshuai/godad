import { ref } from 'vue'
import AdminApi from '@/api/admin'

export function usePendingReports() {
  const pendingCount = ref(0)
  let timer: any = null

  const refresh = async () => {
    try {
      const res: any = await AdminApi.getReports({ status: 'pending', page: 1, size: 1 })
      pendingCount.value = Number(res?.data?.total || 0)
    } catch {
      pendingCount.value = 0
    }
  }

  const startPolling = (intervalMs = 60000) => {
    stopPolling()
    timer = setInterval(refresh, intervalMs)
  }

  const stopPolling = () => {
    if (timer) { clearInterval(timer); timer = null }
  }

  return { pendingCount, refresh, startPolling, stopPolling }
}

