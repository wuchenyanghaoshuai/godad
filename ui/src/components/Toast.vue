<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-50 space-y-2">
      <TransitionGroup
        name="toast"
        tag="div"
        class="space-y-2"
      >
        <div
          v-for="toast in toastsRef"
          :key="toast.id"
          :class="[
            'max-w-sm w-full bg-white shadow-lg rounded-lg pointer-events-auto ring-1 ring-black ring-opacity-5 overflow-hidden',
            'transform transition-all duration-300 ease-in-out'
          ]"
        >
          <div class="p-4">
            <div class="flex items-start">
              <div class="flex-shrink-0">
                <!-- Success Icon -->
                <CheckCircle
                  v-if="toast.type === 'success'"
                  class="h-6 w-6 text-green-400"
                />
                <!-- Error Icon -->
                <XCircle
                  v-else-if="toast.type === 'error'"
                  class="h-6 w-6 text-red-400"
                />
                <!-- Warning Icon -->
                <AlertTriangle
                  v-else-if="toast.type === 'warning'"
                  class="h-6 w-6 text-yellow-400"
                />
                <!-- Info Icon -->
                <Info
                  v-else
                  class="h-6 w-6 text-blue-400"
                />
              </div>
              <div class="ml-3 w-0 flex-1 pt-0.5">
                <p class="text-sm font-medium text-gray-900">
                  {{ toast.title }}
                </p>
                <p v-if="toast.message" class="mt-1 text-sm text-gray-500">
                  {{ toast.message }}
                </p>
              </div>
              <div class="ml-4 flex-shrink-0 flex">
                <button
                  @click="removeToast(toast.id)"
                  class="bg-white rounded-md inline-flex text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                  <span class="sr-only">关闭</span>
                  <X class="h-5 w-5" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import {
  CheckCircle,
  XCircle,
  AlertTriangle,
  Info,
  X
} from 'lucide-vue-next'
import { toasts, useToast } from '@/composables/useToast'

const { removeToast } = useToast()

// 直接使用响应式的 toasts
const toastsRef = toasts
</script>

<style scoped>
/* Toast 动画 */
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.toast-move {
  transition: transform 0.3s ease;
}
</style>
